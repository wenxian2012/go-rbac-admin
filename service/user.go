package service

import (
	"fmt"

	"github.com/wenxian2012/go-rbac-admin/dto"

	"github.com/casbin/casbin"
	"github.com/wenxian2012/go-rbac-admin/models"
)

type User struct {
	Enforcer *casbin.Enforcer `inject:""`
}

func (a *User) GetAll(params *dto.UserQueryParams) (*dto.ResponseList, error) {
	result := dto.ResponseList{
		Pagination: &dto.PaginationParam{
			PageIndex: params.PageIndex,
			PageSize:  params.PageSize,
		},
	}
	user := models.User{}
	users, count, err := user.Query(params)
	if err != nil {
		return nil, err
	}

	list := []*dto.SafeUser{}
	for _, v := range users {
		list = append(list, v.ToDtoUser().ToSafeUser())
	}

	result.Pagination.Total = count
	result.List = list

	return &result, nil
}

// LoadAllPolicy 加载所有的用户策略
func (a *User) LoadAllPolicy() error {
	users, err := models.GetUsersAll()
	if err != nil {
		return err
	}
	for _, user := range users {
		if len(user.Roles) != 0 {
			err = a.LoadPolicy(user.ID)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("角色权限关系", a.Enforcer.GetGroupingPolicy())
	return nil
}

// LoadPolicy 加载用户权限策略
func (a *User) LoadPolicy(id int) error {
	user := models.User{}
	_, err := user.Find(id)
	if err != nil {
		return err
	}

	// a.Enforcer.DeleteRolesForUser(user.Username)
	//
	// for _, ro := range user.Roles {
	// 	a.Enforcer.AddRoleForUser(user.Username, ro.Name)
	// }
	fmt.Println("更新角色权限关系", a.Enforcer.GetGroupingPolicy())
	return nil
}
