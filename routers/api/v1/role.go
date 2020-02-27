package v1

import (
	"github.com/wenxian2012/go-rbac-admin/dto"

	"github.com/wenxian2012/go-rbac-admin/models"
	"github.com/wenxian2012/go-rbac-admin/pkg/errors"

	"github.com/unknwon/com"
	"github.com/wenxian2012/go-rbac-admin/middleware/inject"
	"github.com/wenxian2012/go-rbac-admin/pkg/ginplus"
	"github.com/wenxian2012/go-rbac-admin/service"

	"github.com/gin-gonic/gin"
)

// @Summary   获取单个角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles/:id  [GET]
func GetRole(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	role := models.Role{}
	findOne, err := role.Find(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if findOne.ID == 0 {
		ginplus.ResError(c, errors.ErrNotFound)
		return
	}

	ginplus.ResSuccess(c, findOne)
}

// @Summary   获取所有角色
// @Tags role
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles  [GET]
func GetRoles(c *gin.Context) {
	roleSrv := service.Role{}
	params := dto.RoleQueryParams{
		Name:      c.Query("name"),
		PageNum:   ginplus.GetPageNum(c),
		PageSize:  ginplus.GetPageSize(c),
		PageIndex: ginplus.GetPageIndex(c),
	}

	result, err := roleSrv.GetAll(&params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result)
}

// @Summary   增加角色
// @Tags role
// @Accept json
// @Produce  json
// @Param   body  body   models.Roles   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles  [POST]
func AddRole(c *gin.Context) {
	var (
		reqInfo struct {
			Name      string          `json:"name" binding:"required,max=100,min=2" comment:"角色名"` // 角色名
			Memo      string          `json:"memo"`                                                // 备注
			RoleMenus []*dto.RoleMenu `json:"role_menus" binding:"required,min=1" comment:"菜单"`    // 菜单列表
		}
	)
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		ginplus.ValidatorErrorHandle(c, err)
		return
	}

	if exist, err := models.CheckRoleName(reqInfo.Name); err != nil {
		ginplus.ResError(c, err)
		return
	} else if exist {
		ginplus.ResError(c, errors.New400Response("该角色名已存在"))
		return
	}

	role := models.Role{
		Name: &reqInfo.Name,
		Memo: &reqInfo.Memo,
	}

	if err := role.Create(reqInfo.RoleMenus); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)

}

// @Summary   更新角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Param   body  body   models.Roles   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles/:id  [PUT]
func EditRole(c *gin.Context) {
	var (
		reqInfo struct {
			Name      string          `json:"name" binding:"required,max=100,min=2" comment:"角色名"` // 角色名
			Memo      string          `json:"memo"`                                                // 备注
			RoleMenus []*dto.RoleMenu `json:"role_menus" binding:"required,min=1" comment:"菜单"`    // 菜单列表
		}
	)
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		ginplus.ValidatorErrorHandle(c, err)
		return
	}
	id := com.StrTo(c.Param("id")).MustInt()
	role := models.Role{}

	_, err := role.Find(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	role = models.Role{
		Model: models.Model{
			ID: id,
		},
		Name: &reqInfo.Name,
		Memo: &reqInfo.Memo,
	}

	if err := role.Update(reqInfo.RoleMenus); err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResOK(c)
}

// @Summary   删除角色
// @Tags role
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/roles/:id  [DELETE]
func DeleteRole(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	role := models.Role{}
	findOne, err := role.Find(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	} else if findOne.ID == 0 {
		ginplus.ResError(c, errors.ErrNotFound)
		return
	}
	if err := findOne.Delete(id); err != nil {
		ginplus.ResError(c, err)
		return
	}

	inject.Obj.Enforcer.DeleteRole(*findOne.Name)

	ginplus.ResOK(c)
}
