package models

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/wenxian2012/go-rbac-admin/dto"
)

type Role struct {
	Model
	Name      *string     `gorm:"column:name;"`
	Memo      *string     `gorm:"column:memo"`
	RoleMenus []*RoleMenu `gorm:"association_save_reference:false;"`
}

func (a Role) TableName() string {
	return a.Model.TableName("role")
}

// RoleMenu 角色菜单关联实体
type RoleMenu struct {
	Model
	RoleID   int     `gorm:"column:role_id;"`  // 角色内码
	MenuID   int     `gorm:"column:menu_id;"`  // 菜单内码
	Action   *string `gorm:"column:action;"`   // 动作权限(多个以英文逗号分隔)
	Resource *string `gorm:"column:resource;"` // 资源权限(多个以英文逗号分隔)
}

func (a *Role) ToDtoRole() *dto.Role {
	item := &dto.Role{
		ID:        a.ID,
		Name:      *a.Name,
		Memo:      *a.Memo,
		UpdatedAt: a.UpdatedAt,
	}
	roleMenus := []*dto.RoleMenu{}
	for _, v := range a.RoleMenus {
		roleMenu := &dto.RoleMenu{
			ID:        v.ID,
			MenuID:    v.MenuID,
			Actions:   []string{},
			Resources: []string{},
		}
		if b := v.Action; b != nil && *b != "" {
			roleMenu.Actions = strings.Split(*b, ",")
		}
		if b := v.Resource; b != nil && *b != "" {
			roleMenu.Resources = strings.Split(*b, ",")
		}
		roleMenus = append(roleMenus, roleMenu)
	}
	item.RoleMenus = roleMenus
	return item
}

func ToDtoRoles(a []*Role) []*dto.Role {
	roles := []*dto.Role{}
	for _, v := range a {
		roles = append(roles, v.ToDtoRole())
	}
	return roles
}

func CheckRoleName(name string) (bool, error) {
	var role Role
	err := db.Where("name = ?", name).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if role.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetRolesAll() ([]*Role, error) {
	var role []*Role
	err := db.Preload("RoleMenus").Find(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return role, nil
}

func (a *Role) Query(params *dto.RoleQueryParams) ([]*Role, int, error) {
	var list []*Role
	var count int
	db := db.Model(&Role{})
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	db.Count(&count)
	err := db.Preload("RoleMenus").Offset(params.PageNum).Limit(params.PageSize).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, count, err
	}
	return list, count, nil
}

func (a *Role) Find(id int) (*Role, error) {
	var role Role
	err := db.Preload("RoleMenus").Where("id = ?", id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &role, nil
}

func (a *Role) Update(dtoRoleMenus []*dto.RoleMenu) error {
	if err := db.Save(&a).Error; err != nil {
		return err
	}
	var roleMenus []*RoleMenu
	for _, item := range dtoRoleMenus {
		resources := strings.Join(item.Resources, ",")
		actions := strings.Join(item.Actions, ",")
		roleMenu := &RoleMenu{
			Model: Model{
				ID: item.ID,
			},
			RoleID:   a.ID,
			MenuID:   item.MenuID,
			Resource: &resources,
			Action:   &actions,
		}
		roleMenus = append(roleMenus, roleMenu)
	}

	if err := db.Model(a).Association("RoleMenus").Replace(roleMenus).Error; err != nil {
		return err
	}

	return nil
}

func (a *Role) Create(dtoRoleMenus []*dto.RoleMenu) error {
	if err := db.Create(&a).Error; err != nil {
		return err
	}
	var roleMenus []*RoleMenu
	for _, item := range dtoRoleMenus {
		resources := strings.Join(item.Resources, ",")
		actions := strings.Join(item.Actions, ",")
		roleMenu := &RoleMenu{
			RoleID:   a.ID,
			MenuID:   item.MenuID,
			Resource: &resources,
			Action:   &actions,
		}
		roleMenus = append(roleMenus, roleMenu)
	}

	if err := db.Model(a).Association("RoleMenus").Append(roleMenus).Error; err != nil {
		return err
	}

	return nil
}

func (a *Role) Delete(id int) error {
	if err := db.Delete(Role{Model: Model{ID: id}}).Association("RoleMenus").Clear().Error; err != nil {
		return err
	}
	return nil
}
