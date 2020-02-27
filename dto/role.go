package dto

import "github.com/wenxian2012/go-rbac-admin/pkg/gormkit"

type Role struct {
	ID        int               `json:"id"`         // 主键
	Name      string            `json:"name" `      // 角色名
	Memo      string            `json:"memo"`       // 备注
	RoleMenus []*RoleMenu       `json:"role_menus"` // 拥有菜单
	UpdatedAt gormkit.LocalTime `json:"updated_at"` // 创建时间
}

// RoleMenu 角色菜单对象
type RoleMenu struct {
	ID        int      `json:"id"`        // 主键
	MenuID    int      `json:"menu_id"`   // 菜单ID
	Actions   []string `json:"actions"`   // 动作权限列表
	Resources []string `json:"resources"` // 资源权限列表
}

type RoleQueryParams struct {
	Name      string // 角色名
	PageNum   int    // 分页计算
	PageSize  int    // 分页条数
	PageIndex int    // 当前页
}
