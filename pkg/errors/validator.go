package errors

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

var (
	trans        ut.Translator
	uni          *ut.UniversalTranslator
	errorHandler *AppErrorHandler
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		uni = ut.New(zh.New())
		trans, _ = uni.GetTranslator("zh")
		// 收集结构体中的comment标签，用于替换英文字段名称
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("comment")
		})
		// this is usually know or extracted from http 'Accept-Language' header
		// also see uni.FindTranslator(...)
		zh_translations.RegisterDefaultTranslations(v, trans)

		errorHandler = NewErrorHandler(uni, trans, v)
	}
}

// 错误处理
func ValidateError() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorHandler.HandleErrors(c)
	}
}

type AppErrorHandler struct {
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	validate *validator.Validate
}

func NewErrorHandler(uni *ut.UniversalTranslator, trans ut.Translator, validate *validator.Validate) *AppErrorHandler {
	return &AppErrorHandler{
		uni:      uni,
		trans:    trans,
		validate: validate,
	}
}

func (h *AppErrorHandler) HandleErrors(c *gin.Context) {
	c.Next()
	errorToPrint := c.Errors.ByType(gin.ErrorTypePublic).Last()
	if errorToPrint != nil {
		if errs, ok := errorToPrint.Err.(validator.ValidationErrors); ok {
			trans, _ := h.uni.GetTranslator("zh")
			errors := make(map[string]interface{})
			for _, v := range errs {
				errors[v.StructNamespace()] = v.Translate(trans)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"msg":    errs[0].Translate(trans),
				"errors": errors,
			})
		}
		// deal with other errors ...
	}
}
