package service

import (
	"github.com/wenxian2012/go-rbac-admin/dto"

	"github.com/casbin/casbin"
	"github.com/wenxian2012/go-rbac-admin/models"
)

type Role struct {
	Enforcer *casbin.Enforcer `inject:""`
}

func (a *Role) GetAll(params *dto.RoleQueryParams) (*dto.ResponseList, error) {
	result := dto.ResponseList{
		Pagination: &dto.PaginationParam{
			PageIndex: params.PageIndex,
			PageSize:  params.PageSize,
		},
	}
	role := models.Role{}
	roles, count, err := role.Query(params)
	if err != nil {
		return nil, err
	}

	result.Pagination.Total = count
	result.List = models.ToDtoRoles(roles)

	return &result, nil
}

// LoadAllPolicy 加载所有的角色策略
func (a *Role) LoadAllPolicy() error {
	roles, err := models.GetRolesAll()
	if err != nil {
		return err
	}

	for _, role := range roles {
		err = a.LoadPolicy(role.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadPolicy 加载角色权限策略
func (a *Role) LoadPolicy(id int) error {
	role := &models.Role{}
	role, err := role.Find(id)
	if err != nil {
		return err
	}
	a.Enforcer.DeleteRole(*role.Name)

	//for _, menu := range role.RoleMenus {
	//	if menu.Path == "" || menu.Method == "" {
	//		continue
	//	}
	//	a.Enforcer.AddPermissionForUser(role.Name, menu.Path, menu.Method)
	//}
	return nil
}
