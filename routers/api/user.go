package api

import (
	"github.com/wenxian2012/go-rbac-admin/dto"
	"github.com/wenxian2012/go-rbac-admin/middleware/inject"
	"github.com/wenxian2012/go-rbac-admin/pkg/errors"

	"github.com/wenxian2012/go-rbac-admin/models"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/wenxian2012/go-rbac-admin/pkg/ginplus"
	"github.com/wenxian2012/go-rbac-admin/pkg/util"
	"github.com/wenxian2012/go-rbac-admin/service"
)

const zero = 0

// @Summary   获取登录token 信息
// @Tags auth
// @Accept json
// @Produce  json
// @Param   body  body   models.AuthSwag   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /auth  [POST]
func Auth(c *gin.Context) {
	var (
		reqInfo struct {
			Username string `json:"username" binding:"required,max=100,min=3" comment:"用户名"` // 用户名
			Password string `json:"password" binding:"required,max=100,min=3" comment:"密码"`  // 密码
		}
	)
	// dataByte, _ := ioutil.ReadAll(c.Request.Body)
	// fsion := gofasion.NewFasion(string(dataByte))
	// fmt.Println(fsion.Get("username").ValueStr())

	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		ginplus.ValidatorErrorHandle(c, err)
		return
	}

	if exist, err := models.CheckUsername(reqInfo.Username); err != nil {
		ginplus.ResError(c, err)
		return
	} else if !exist {
		ginplus.ResError(c, errors.New400Response("该用户名不存在"))
		return
	}
	user, err := models.CheckUser(reqInfo.Username, util.EncodeMD5(reqInfo.Password))
	if err != nil {
		ginplus.ResError(c, err)
		return
	} else if user == nil {
		ginplus.ResError(c, errors.New400Response("用户名密码不正确"))
		return
	}

	token, err := util.GenerateToken(reqInfo.Username, reqInfo.Password)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	menuService := service.Menu{}
	menus, err := menuService.GetTree()
	if err != nil {
		ginplus.ResError(c, err)
	}

	ginplus.ResSuccess(c, gin.H{
		"token": token,
		"user":  user.ToDtoUser().ToSafeUser(),
		"menus": menus,
	})
}

func Token(c *gin.Context) {
	token := ginplus.GetToken(c)
	claims, err := util.ParseToken(token)
	if err != nil {
		ginplus.ResError(c, errors.ErrTimeOutToken)
		return
	}
	user, err := models.CheckUser(claims.Username, claims.Password)
	if err != nil {
		ginplus.ResError(c, err)
		return
	} else if user == nil {
		ginplus.ResError(c, errors.ErrInvalidUser)
		return
	}
	ginplus.ResSuccess(c, user.ToDtoUser().ToSafeUser())
}

// @Summary   获取单个用户信息
// @Tags  users
// @Accept json
// @Produce  json
// @Param  id  path   int true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users/:id  [GET]
func GetUser(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	user := models.User{}
	findOne, err := user.Find(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	} else if findOne.ID == 0 {
		ginplus.ResError(c, errors.ErrNotFound)
		return
	}
	safeUser := findOne.ToDtoUser().ToSafeUser()
	ginplus.ResSuccess(c, safeUser)
}

// @Summary   获取所有用户
// @Tags  users
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [GET]
func GetUsers(c *gin.Context) {
	userSrv := service.User{}
	params := dto.UserQueryParams{
		PageNum:   ginplus.GetPageNum(c),
		PageIndex: ginplus.GetPageIndex(c),
		PageSize:  ginplus.GetPageSize(c),
	}
	_ = c.ShouldBindQuery(&params)
	result, err := userSrv.GetAll(&params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result)
}

// @Summary   增加用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users  [POST]
func AddUser(c *gin.Context) {
	var (
		reqInfo struct {
			Username string `json:"username" binding:"required,max=100,min=3" comment:"用户名"`
			Nickname string `json:"nickname" binding:"required,max=100,min=3" comment:"姓名"`
			Password string `json:"password" binding:"required,max=100,min=3" comment:"密码"`
			Phone    string `json:"phone" binding:"omitempty,len=11" comment:"手机号"`
			Email    string `json:"email" binding:"omitempty,email" comment:"邮箱"`
			Disabled int    `json:"disabled" binding:"gte=0,lte=1" comment:"禁用"`
			Roles    []int  `json:"roles" binding:"required,max=10,min=1" comment:"角色"`
		}
	)
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		ginplus.ValidatorErrorHandle(c, err)
		return
	}

	if exist, err := models.CheckUsername(reqInfo.Username); err != nil {
		ginplus.ResError(c, err)
		return
	} else if exist {
		ginplus.ResError(c, errors.New400Response("该用户名已存在"))
		return
	}

	reqInfo.Password = util.EncodeMD5(reqInfo.Password)

	user := models.User{
		Username: &reqInfo.Username,
		Password: &reqInfo.Password,
		Nickname: &reqInfo.Nickname,
		Phone:    &reqInfo.Phone,
		Email:    &reqInfo.Email,
		Disabled: &reqInfo.Disabled,
	}

	if err := user.Create(reqInfo.Roles); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// @Summary   更新用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users/:id  [PUT]
func EditUser(c *gin.Context) {
	var (
		reqInfo struct {
			Username string `json:"username" binding:"required,max=100,min=3" comment:"用户名"`
			Nickname string `json:"nickname" binding:"required,max=100,min=3" comment:"姓名"`
			Password string `json:"password" binding:"omitempty,max=100,min=3" comment:"密码"`
			Phone    string `json:"phone" binding:"omitempty,len=11" comment:"手机号"`
			Email    string `json:"email" binding:"omitempty,email" comment:"邮箱"`
			Disabled int    `json:"disabled" binding:"gte=0,lte=1" comment:"禁用"`
			Roles    []int  `json:"roles" binding:"required,max=10,min=1" comment:"角色"`
		}
	)
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		ginplus.ValidatorErrorHandle(c, err)
		return
	}
	id := com.StrTo(c.Param("id")).MustInt()
	user := models.User{}

	findOne, err := user.Find(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	} else if findOne.ID == 0 {
		ginplus.ResError(c, errors.ErrNotFound)
		return
	}

	if reqInfo.Password == "" {
		reqInfo.Password = *findOne.Password
	} else {
		reqInfo.Password = util.EncodeMD5(reqInfo.Password)
	}

	if reqInfo.Username != *findOne.Username {
		if exist, err := models.CheckUsername(reqInfo.Username); err != nil {
			ginplus.ResError(c, err)
			return
		} else if exist {
			ginplus.ResError(c, errors.New400Response("该用户名已存在"))
			return
		}
	}

	user = models.User{
		Model: models.Model{
			ID: id,
		},
		Username: &reqInfo.Username,
		Password: &reqInfo.Password,
		Nickname: &reqInfo.Nickname,
		Phone:    &reqInfo.Phone,
		Email:    &reqInfo.Email,
		Disabled: &reqInfo.Disabled,
	}

	if err := user.Update(reqInfo.Roles); err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResOK(c)
}

// @Summary   删除用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param  id  path  int true "id"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Router /api/v1/users/:id  [DELETE]
func DeleteUser(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	user := models.User{}
	findOne, err := user.Find(id)
	if err != nil {
		ginplus.ResError(c, err)
		return
	} else if findOne.ID == 0 {
		ginplus.ResError(c, errors.ErrNotFound)
		return
	} else if findOne.ID == 1 {
		ginplus.ResError(c, errors.ErrNotAllowDelete)
		return
	}

	if err := findOne.Delete(id); err != nil {
		ginplus.ResError(c, err)
		return
	}

	inject.Obj.Enforcer.DeleteUser(*findOne.Username)

	ginplus.ResOK(c)
}
