package models

import (
	"strconv"

	"github.com/wenxian2012/go-rbac-admin/pkg/errors"

	"github.com/jinzhu/gorm"
	"github.com/wenxian2012/go-rbac-admin/dto"
)

type Menu struct {
	Model
	Name       *string         `gorm:"column:name;"`                                                                  // 菜单名称
	Sequence   *int            `gorm:"column:sequence;"`                                                              // 排序值
	Icon       *string         `gorm:"column:icon;"`                                                                  // 菜单图标
	Router     *string         `gorm:"column:router;"`                                                                // 访问路由
	Hidden     *int            `gorm:"column:hidden;"`                                                                // 隐藏菜单(0:不隐藏 1:隐藏)
	ParentID   *int            `gorm:"column:parent_id;"`                                                             // 父级内码
	ParentPath *string         `gorm:"column:parent_path;"`                                                           // 父级路径
	Resources  []*MenuResource `gorm:"association_save_reference:false;foreignkey:MenuID;association_foreignkey:ID;"` // 资源列表
	Actions    []*MenuAction   `gorm:"association_save_reference:false;foreignkey:MenuID;association_foreignkey:ID;"` // 动作列表
}

func (a Menu) TableName() string {
	return a.Model.TableName("menu")
}

func (a *Menu) ToDtoMenu() *dto.Menu {
	item := &dto.Menu{
		ID:         a.ID,
		Name:       *a.Name,
		Router:     *a.Router,
		Hidden:     *a.Hidden,
		Sequence:   *a.Sequence,
		Icon:       *a.Icon,
		ParentID:   *a.ParentID,
		ParentPath: *a.ParentPath,
		CreatedAt:  a.CreatedAt,
	}
	if a.Router != nil {
		item.Router = *a.Router
	}
	if a.ParentID != nil {
		item.ParentID = *a.ParentID
	}
	if a.ParentPath != nil {
		item.ParentPath = *a.ParentPath
	}
	return item
}

// 转实体Menu
func BelongToMenu(a *dto.Menu) *Menu {
	v := &Menu{
		Model: Model{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
		},
		Name:       &a.Name,
		Sequence:   &a.Sequence,
		Icon:       &a.Icon,
		Router:     &a.Router,
		Hidden:     &a.Hidden,
		ParentID:   &a.ParentID,
		ParentPath: &a.ParentPath,
	}
	return v
}

// 转实体FillMenu
func BelongToFillMenu(a *dto.FillMenu) *Menu {
	v := &Menu{
		Model: Model{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
		},
		Name:       &a.Name,
		Sequence:   &a.Sequence,
		Icon:       &a.Icon,
		Router:     &a.Router,
		Hidden:     &a.Hidden,
		ParentID:   &a.ParentID,
		ParentPath: &a.ParentPath,
	}
	if a.Resources != nil {
		var resources []*MenuResource
		for _, item := range a.Resources {
			resources = append(resources, BelongToResource(item, a.ID))
		}
		v.Resources = resources
	}
	if a.Actions != nil {
		var actions []*MenuAction
		for _, item := range a.Actions {
			actions = append(actions, BelongToAction(item, a.ID))
		}
		v.Actions = actions
	}
	return v
}

// 转实体Resource
func BelongToResource(a *dto.MenuResource, menuID int) *MenuResource {
	v := &MenuResource{
		Model: Model{
			ID: a.ID,
		},
		MenuID: menuID,
		Name:   a.Name,
		Code:   a.Code,
		Method: a.Method,
		Path:   a.Path,
	}
	return v
}

// 转实体Action
func BelongToAction(a *dto.MenuAction, menuID int) *MenuAction {
	v := &MenuAction{
		Model: Model{
			ID: a.ID,
		},
		MenuID: menuID,
		Code:   a.Code,
		Name:   a.Name,
	}
	return v
}

// Menus 菜单实体列表
type Menus []*Menu

// MenuAction 菜单动作关联实体
type MenuAction struct {
	Model
	MenuID int    `gorm:"column:menu_id;int;size:36;index;"` // 菜单ID
	Code   string `gorm:"column:code;size:50;index;"`        // 动作编号
	Name   string `gorm:"column:name;size:50;"`              // 动作名称
}

// TableName 表名
func (a MenuAction) TableName() string {
	return a.Model.TableName("menu_action")
}

func (a MenuAction) ToPogoMenuActive() *dto.MenuAction {
	return &dto.MenuAction{
		ID:   a.ID,
		Code: a.Code,
		Name: a.Name,
	}
}

// MenuActions 菜单动作关联实体列表
type MenuActions []*MenuAction

// MenuResource 菜单资源关联实体
type MenuResource struct {
	Model
	MenuID int    `gorm:"column:menu_id;size:36;index;"` // 菜单ID
	Code   string `gorm:"column:code;size:50;index;"`    // 资源编号
	Name   string `gorm:"column:name;size:50;"`          // 资源名称
	Method string `gorm:"column:method;size:50;"`        // 请求方式
	Path   string `gorm:"column:path;size:255;"`         // 请求路径
}

// TableName 表名
func (a MenuResource) TableName() string {
	return a.Model.TableName("menu_resource")
}

func (a *MenuResource) ToPogoMenuResource() *dto.MenuResource {
	return &dto.MenuResource{
		ID:     a.ID,
		Code:   a.Code,
		Name:   a.Name,
		Method: a.Method,
		Path:   a.Path,
	}
}

// MenuResources 菜单资源关联实体列表
type MenuResources []*MenuResource

// 填充菜单对象数据
func (a Menus) DoFillMenus(includeActions bool, includeResources bool) (dto.FillMenusTree, error) {

	menuIDs := make([]int, len(a))
	for i, item := range a {
		menuIDs[i] = item.ID
	}

	var actionList MenuActions
	var resourceList MenuResources
	var list1 MenuActions
	err := db.Where("menu_id IN(?)", menuIDs).Find(&list1).Error
	if err != nil {
		return nil, err
	}
	actionList = list1

	var list2 MenuResources
	err = db.Where("menu_id IN(?)", menuIDs).Find(&list2).Error
	if err != nil {
		return nil, err
	}
	resourceList = list2

	var result dto.FillMenusTree

	for _, item := range a {
		_item := &dto.FillMenuTree{
			FillMenu: dto.FillMenu{
				Menu: *item.ToDtoMenu(),
			},
		}

		if includeActions {
			_item.Actions = []*dto.MenuAction{}
		}

		if includeResources {
			_item.Resources = []*dto.MenuResource{}
		}

		if includeActions && len(actionList) > 0 {
			_item.Actions = actionList.GetByMenuID(item.ID)
		}

		if includeResources && len(resourceList) > 0 {
			_item.Resources = resourceList.GetByMenuID(item.ID)
		}
		result = append(result, _item)
	}

	return result, nil
}

// GetByMenuID 根据菜单ID获取菜单动作列表
func (a MenuActions) GetByMenuID(menuID int) []*dto.MenuAction {
	list := make([]*dto.MenuAction, 0)
	for _, item := range a {
		if item.MenuID == menuID {
			list = append(list, &dto.MenuAction{
				ID:   item.ID,
				Code: item.Code,
				Name: item.Name,
			})
		}
	}
	return list
}

// GetByMenuID 根据菜单ID获取菜单资源列表
func (a MenuResources) GetByMenuID(menuID int) []*dto.MenuResource {
	list := make([]*dto.MenuResource, 0)
	for _, item := range a {
		if item.MenuID == menuID {
			list = append(list, &dto.MenuResource{
				ID:     item.ID,
				Code:   item.Code,
				Name:   item.Name,
				Method: item.Method,
				Path:   item.Path,
			})
		}
	}
	return list
}

func GetTree() (dto.FillMenusTree, error) {
	var list Menus
	err := db.Order("sequence DESC,id DESC").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	result, err := list.DoFillMenus(true, true)
	if err != nil {
		return nil, err
	}

	result = result.ToNested()

	return result, nil
}

// 获取父级路径
func (a *Menu) getParentPath(parentID int) (string, error) {
	if parentID == 0 {
		return "", nil
	}
	pitem, err := a.Find(parentID)
	if err != nil {
		return "", errors.ErrInvalidParent
	}
	return a.joinParentPath(pitem.ToDtoMenu().ParentPath, strconv.Itoa(pitem.ID)), nil
}

func (a *Menu) joinParentPath(ppath string, code string) string {
	if ppath != "" {
		ppath += "/"
	}
	return ppath + code
}

func (a Menu) Query(params *dto.MenuQueryParams) (Menus, int, error) {
	var list = Menus{}
	var count int
	db := db.Model(&Menu{})
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}

	if v := params.ParentID; v != nil {
		db = db.Where("parent_id = ?", v)
	}

	if v := params.Hidden; v != nil {
		db = db.Where("hidden = ?", v)
	}
	db.Count(&count)
	err := db.Order("sequence DESC,id DESC").Offset(params.PageNum).Limit(params.PageSize).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, count, err
	}

	return list, count, nil
}

func (a *Menu) Find(id int) (*Menu, error) {
	var menu Menu
	err := db.Where("id = ?", id).First(&menu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &menu, nil
}

func (a *Menu) Create() error {
	parentId := *a.ParentID
	parentPath, err := a.getParentPath(parentId)
	if err != nil {
		return err
	}
	a.ParentPath = &parentPath

	if err := db.Create(&a).Error; err != nil {
		return err
	}

	for _, item := range a.Resources {
		item.MenuID = a.ID
	}

	for _, item := range a.Actions {
		item.MenuID = a.ID
	}

	if err := db.Model(a).Association("Resources").Append(a.Resources).Error; err != nil {
		return err
	}

	if err := db.Model(a).Association("Actions").Append(a.Actions).Error; err != nil {
		return err
	}

	return nil
}

func (a *Menu) Update() error {
	parentId := *a.ParentID
	parentPath, err := a.getParentPath(parentId)
	if err != nil {
		return err
	}
	a.ParentPath = &parentPath

	if err := db.Save(&a).Error; err != nil {
		return err
	}

	if err := db.Model(a).Association("Resources").Replace(a.Resources).Error; err != nil {
		return err
	}

	if err := db.Model(a).Association("Actions").Replace(a.Actions).Error; err != nil {
		return err
	}

	return nil
}

func (a *Menu) Delete(id int) error {
	if err := db.Where("id = ?", id).Delete(Menu{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *Menu) UpdateMenuGetRoles(id int) []int {
	var menu Menu
	var role []Role

	db.Model(&menu).Where("id = ?", id)
	db.Joins(" left join go_role_menu b on go_role.id=b.role_id left join go_menu c on c.id=b.menu_id").Where("c.id = ?", id).Find(&role)

	roleList := []int{}
	for _, v := range role {
		roleList = append(roleList, v.ID)
	}
	return roleList
}
