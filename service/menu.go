package service

import (
	"github.com/casbin/casbin"
	"github.com/wenxian2012/go-rbac-admin/dto"
	"github.com/wenxian2012/go-rbac-admin/models"
)

type Menu struct {
	Menu     *models.Menu     `inject:""`
	Enforcer *casbin.Enforcer `inject:""`
}

func (a *Menu) Get(id int) (*dto.FillMenuTree, error) {
	menu := models.Menu{}
	findOne, err := menu.Find(id)
	if err != nil {
		return nil, err
	}

	list, err := models.Menus{findOne}.DoFillMenus(true, true)
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

func (a *Menu) GetAll(params *dto.MenuQueryParams) (*dto.ResponseList, error) {
	result := dto.ResponseList{
		Pagination: &dto.PaginationParam{
			PageIndex: params.PageIndex,
			PageSize:  params.PageSize,
		},
	}
	menu := models.Menu{}
	menus, count, err := menu.Query(params)
	if err != nil {
		return nil, err
	}

	list := []*dto.Menu{}
	for _, v := range menus {
		list = append(list, v.ToDtoMenu())
	}

	if params.IncludeResources || params.IncludeActions {
		list2 := models.Menus{}
		for _, v := range list {
			vv := models.BelongToMenu(v)
			list2 = append(list2, vv)
		}
		list3, err := list2.DoFillMenus(params.IncludeActions, params.IncludeResources)
		if err != nil {
			return nil, err
		}
		result.List = list3
	} else {
		result.List = list
	}

	result.Pagination.Total = count

	return &result, nil
}

func (a *Menu) GetTree() (dto.FillMenusTree, error) {
	Menus, err := models.GetTree()
	if err != nil {
		return nil, err
	}

	return Menus, nil
}

func (a *Menu) Edit(menu *models.Menu) error {
	err := menu.Update()
	if err != nil {
		return err
	}
	// roleList := menu.UpdateMenuGetRoles(menu.ID)
	// roleService := Role{}
	// for _, v := range roleList {
	// 	err := roleService.LoadPolicy(v)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (a *Menu) Delete(id int) error {
	menu := models.Menu{}
	err := menu.Delete(id)
	if err != nil {
		return err
	}
	// roleList := menu.UpdateMenuGetRoles(id)
	// roleService := Role{}
	// for _, v := range roleList {
	// 	err := roleService.LoadPolicy(v)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
