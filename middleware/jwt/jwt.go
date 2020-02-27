package jwt

import (
	"strings"

	"github.com/wenxian2012/go-rbac-admin/pkg/errors"
	"github.com/wenxian2012/go-rbac-admin/pkg/ginplus"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/wenxian2012/go-rbac-admin/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		Authorization := c.GetHeader("Authorization")
		token := strings.Split(Authorization, " ")

		if Authorization == "" {
			ginplus.ResError(c, errors.NewResponse(403, "Not Found Token", 403))
			return
		} else {
			_, err := util.ParseToken(token[1])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ginplus.ResError(c, errors.ErrInvalidToken)
					return
				default:
					ginplus.ResError(c, errors.NewResponse(403, "Invalid Token"), 403)
					return
				}
			}
		}

		c.Next()
	}
}
