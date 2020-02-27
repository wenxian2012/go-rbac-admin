package dto

import (
	"github.com/wenxian2012/go-rbac-admin/pkg/gormkit"
	"gopkg.in/go-playground/validator.v8"
)

// RoleMenus 菜单对象
type Menu struct {
	ID         int               `json:"id"`                                                    // 主键
	Name       string            `json:"name" binding:"required,max=10,min=3" comment:"菜单名称"`   // 菜单名称
	Router     string            `json:"router" binding:"required,max=20,min=3" comment:"访问路由"` // 访问路由
	Hidden     int               `json:"hidden" binding:"gte=0" comment:"隐藏设置"`                 // 隐藏菜单(0:不隐藏 1:隐藏)
	Sequence   int               `json:"sequence"`                                              // 排序值
	Icon       string            `json:"icon"`                                                  // 菜单图标
	ParentID   int               `json:"parent_id"`                                             // 父级ID
	ParentPath string            `json:"parent_path"`                                           // 父级路径
	CreatedAt  gormkit.LocalTime `json:"created_at"`                                            // 创建时间
}

type FillMenu struct {
	Menu
	Resources []*MenuResource `json:"resources" binding:"required" comment:"资源"` // 资源列表
	Actions   []*MenuAction   `json:"actions" binding:"required" comment:"动作"`   // 动作列表
}

// FillMenus 菜单列表
type FillMenus []*FillMenu

// FillMenuTree 菜单树
type FillMenuTree struct {
	FillMenu
	Children *[]*FillMenuTree `json:"children,omitempty"` // 子级树
}

// MenuAction 菜单动作对象
type MenuAction struct {
	ID   int    `json:"id"`   // 主键
	Code string `json:"code"` // 动作编号
	Name string `json:"name"` // 动作名称
}

// MenuResource 菜单资源对象
type MenuResource struct {
	ID     int    `json:"id"`     // 主键
	Code   string `json:"code"`   // 资源编号
	Name   string `json:"name"`   // 资源名称
	Method string `json:"method"` // 请求方式
	Path   string `json:"path"`   // 请求路径
}

// FillMenusTree 菜单树列表
type FillMenusTree []*FillMenuTree

// MenuQueryParams 查询条件
type MenuQueryParams struct {
	Name             string  // 菜单名称
	ParentID         *string // 父级内码
	Hidden           *int    // 隐藏菜单
	IncludeActions   bool    // 包含动作列表
	IncludeResources bool    // 包含资源列表
	PageNum          int     // 分页计算
	PageSize         int     // 分页条数
	PageIndex        int     // 当前页
}

// ToNested 转换为树形结构
func (a FillMenusTree) ToNested() FillMenusTree {
	mi := make(map[int]*FillMenuTree)
	for _, item := range a {
		if item.Actions == nil {
			item.Actions = []*MenuAction{}
		}
		if item.Resources == nil {
			item.Resources = []*MenuResource{}
		}
		if item.Children == nil {
			item.Children = &[]*FillMenuTree{}
		}
		mi[item.ID] = item
	}

	var list FillMenusTree
	for _, item := range a {
		if item.ParentID == 0 {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

// validator.v8 错误处理示例
func (r *Menu) GetError(err validator.ValidationErrors) string {
	msg := "参数错误"
	if val, exist := err["RoleMenus.Name"]; exist {
		if val.Field == "Name" {
			switch val.Tag {
			case "required":
				msg = "菜单名称不能为空"
			case "min":
				msg = "菜单名称不能少于" + val.Param + "个字符"
			case "max":
				msg = "菜单名称不能大于" + val.Param + "个字符"
			}
		}
	}

	if val, exist := err["RoleMenus.Router"]; exist {
		if val.Field == "Router" {
			switch val.Tag {
			case "required":
				msg = "访问路由不能为空"
			case "min":
				msg = "访问路由不能少于" + val.Param + "个字符"
			case "max":
				msg = "访问路由不能大于" + val.Param + "个字符"
			}
		}
	}

	if val, exist := err["RoleMenus.Hidden"]; exist {
		if val.Field == "Hidden" {
			switch val.Tag {
			case "required":
				msg = "隐藏路由不能为空"
			}
		}
	}

	return msg
}
