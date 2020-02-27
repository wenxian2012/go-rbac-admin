package dto

import (
	"github.com/wenxian2012/go-rbac-admin/pkg/gormkit"
)

type User struct {
	ID        int               `json:"id"`         // 主键
	Username  string            `json:"username"`   // 用户名
	Password  string            `json:"password"`   // 密码
	Nickname  string            `json:"nickname"`   // 姓名
	Phone     string            `json:"phone"`      // 手机号
	Email     string            `json:"email"`      // 邮箱
	Disabled  int               `json:"disabled"`   // 禁用
	Roles     []*Role           `json:"roles"`      // 用户角色
	UpdatedAt gormkit.LocalTime `json:"updated_at"` // 更新时间
}

type SafeUser struct {
	ID        int               `json:"id"`         // 主键
	Username  string            `json:"username"`   // 用户名
	Nickname  string            `json:"nickname"`   // 姓名
	Phone     string            `json:"phone"`      // 手机号
	Email     string            `json:"email"`      // 邮箱
	Disabled  int               `json:"disabled"`   // 禁用
	Roles     []*Role           `json:"roles"`      // 用户角色
	UpdatedAt gormkit.LocalTime `json:"updated_at"` // 更新时间
}

func (a User) ToSafeUser() *SafeUser {
	return &SafeUser{
		ID:        a.ID,
		Username:  a.Username,
		Nickname:  a.Nickname,
		Phone:     a.Phone,
		Email:     a.Email,
		Disabled:  a.Disabled,
		Roles:     a.Roles,
		UpdatedAt: a.UpdatedAt,
	}
}

type UserQueryParams struct {
	Username  string `form:"username"` // 用户名
	Nickname  string `form:"nickname"` // 姓名
	Phone     string `form:"phone"`    // 手机号
	Email     string `form:"email"`    // 邮箱
	RoleId    int    `form:"roleId"`   // 角色
	Disabled  *int   `form:"disabled"` // 禁用
	PageNum   int    // 分页计算
	PageSize  int    // 分页条数
	PageIndex int    // 当前页
}
