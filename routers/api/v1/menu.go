package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/wenxian2012/go-rbac-admin/dto"
	"github.com/wenxian2012/go-rbac-admin/models"
	"github.com/wenxian2012/go-rbac-admin/pkg/errors"
	"github.com/wenxian2012/go-rbac-admin/pkg/ginplus"
	"github.com/wenxian2012/go-rbac-admin/pkg/util"
	"github.com/wenxian2012/go-rbac-admin/service"
)

// @Summary   获取单个菜单
// @Tags menu
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/menus/:id  [GET]
func GetMenu(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	menuSrv := service.Menu{}
	menu, err := menuSrv.Get(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, menu)
}

// @Summary   获取所有菜单
// @Tags menu
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/menus  [GET]
func GetMenus(c *gin.Context) {

	menuService := service.Menu{}
	params := dto.MenuQueryParams{
		Name:             c.Query("name"),
		IncludeActions:   c.Query("includeActions") == "1",
		IncludeResources: c.Query("includeResources") == "1",
		PageNum:          ginplus.GetPageNum(c),
		PageSize:         ginplus.GetPageSize(c),
		PageIndex:        ginplus.GetPageIndex(c),
	}

	if v := c.Query("parentId"); v != "" {
		params.ParentID = &v
	}

	if v := c.Query("hidden"); v != "" {
		if hidden := util.S(v).DefaultInt(0); hidden >= 0 {
			params.Hidden = &hidden
		}
	}

	result, err := menuService.GetAll(&params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result)
}

// @Summary   获取所有菜单 - 树结构
// @Tags menu
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/menus-tree  [GET]
func GetTree(c *gin.Context) {
	menuService := service.Menu{}
	result, err := menuService.GetTree()
	if err != nil {
		ginplus.ResError(c, err)
	}
	ginplus.ResSuccess(c, result)
}

// @Summary   增加菜单
// @Tags menu
// @Accept json
// @Produce  json
// @Param   body  body   models.RoleMenus   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/menus  [POST]
func AddMenu(c *gin.Context) {
	var (
		reqInfo dto.FillMenu
	)
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		ginplus.ValidatorErrorHandle(c, err)
		return
	}

	if reqInfo.Sequence == 0 {
		reqInfo.Sequence = 1000
	}

	menu := models.BelongToFillMenu(&reqInfo)
	if err := menu.Create(); err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResOK(c)
}

// @Summary   更新菜单
// @Tags menu
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Param   body  body   models.RoleMenus   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/menus/:id  [PUT]
func EditMenu(c *gin.Context) {

	id := com.StrTo(c.Param("id")).MustInt()
	reqInfo := dto.FillMenu{
		Menu: dto.Menu{
			ID: id,
		},
	}
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		ginplus.ValidatorErrorHandle(c, err)
		return
	}
	menu := &models.Menu{}
	findOne, err := menu.Find(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	} else if findOne.ID == reqInfo.ParentID {
		ginplus.ResError(c, errors.ErrInvalidParent)
	}
	if reqInfo.Sequence == 0 {
		reqInfo.Sequence = 1000
	}
	menu = models.BelongToFillMenu(&reqInfo)
	menuSrv := service.Menu{}
	err = menuSrv.Edit(menu)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResOK(c)
}

// @Summary   删除菜单
// @Tags menu
// @Accept json
// @Produce  json
// @Param  id  path  string true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/menus/:id  [DELETE]
func DeleteMenu(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	menu := models.Menu{}
	findOne, err := menu.Find(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	} else if findOne.ID == 0 {
		ginplus.ResError(c, errors.ErrNotFound)
		return
	}
	menuSrv := service.Menu{}
	err = menuSrv.Delete(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResOK(c)
}
