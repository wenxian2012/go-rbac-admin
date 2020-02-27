package ginplus

import (
	"net/http"
	"strings"

	"github.com/wenxian2012/go-rbac-admin/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/wenxian2012/go-rbac-admin/dto"
	"github.com/wenxian2012/go-rbac-admin/pkg/errors"
	"github.com/wenxian2012/go-rbac-admin/pkg/util"
	"gopkg.in/go-playground/validator.v9"
)

func ValidatorErrorHandle(c *gin.Context, err error) {
	if err, exist := err.(validator.ValidationErrors); exist {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	ResError(c, errors.Wrap400Response(err, "请求参数错误"))
}

func GetPageNum(c *gin.Context) int {
	page := GetPageIndex(c)
	return (page - 1) * GetPageSize(c)
}

// GetPageIndex 获取分页的页索引
func GetPageIndex(c *gin.Context) int {
	defaultVal := 1
	if v := c.Query("pageIndex"); v != "" {
		if iv := util.S(v).DefaultInt(defaultVal); iv > 0 {
			return iv
		}
	}
	return defaultVal
}

// GetPageSize 获取分页的页大小
func GetPageSize(c *gin.Context) int {
	defaultVal := setting.AppSetting.PageSize
	if v := c.Query("pageSize"); v != "" {
		if iv := util.S(v).DefaultInt(defaultVal); iv > 0 {
			return iv
		} else if iv == -1 {
			return 10000000
		}
	}
	return defaultVal
}

// GetToken 获取用户令牌
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Token "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// ResPage 响应分页数据
func ResPage(c *gin.Context, v *dto.ResponseList) {
	ResSuccess(c, v)
}

// ResOK 响应OK
func ResOK(c *gin.Context) {
	ResSuccess(c, "ok")
}

// ResSuccess 响应成功
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

// ResJSON 响应JSON数据
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := util.JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// ResError 响应错误
func ResError(c *gin.Context, err error, status ...int) {
	var resError *errors.ResponseError
	if err != nil {
		if re, ok := err.(*errors.ResponseError); ok {
			resError = re
		} else {
			resError = errors.UnWrapResponse(errors.Wrap500Response(err))
		}
	} else {
		resError = errors.UnWrapResponse(errors.ErrInternalServer)
	}

	dtoError := dto.ResponseError{
		Code:    resError.Code,
		Message: resError.Message,
	}

	if err := resError.ERR; err != nil {
		dtoError.Error = err.Error()
		if status := resError.Status; status >= 400 && status < 500 {
			// logger.StartSpan(NewContext(c)).Warnf(err.Error())
		} else if status >= 500 {
			// span := logger.StartSpan(NewContext(c))
			// span = span.WithField("stack", fmt.Sprintf("%+v", err))
			// span.Errorf(err.Error())
		}
	}

	ResJSON(c, resError.Status, dtoError)
}
