package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pokemon/pkg/services"
	"pokemon/pkg/utils"
	"strings"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = http.StatusOK
		token := c.GetHeader("Authorization")
		if token == "" {
			code = http.StatusUnauthorized
		} else {
			token = strings.Replace(token, "Bearer ", "", -1)
			claims, err := utils.ParseToken(token)
			if err != nil {
				code = http.StatusUnauthorized
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = http.StatusNonAuthoritativeInfo
			} else {
				rr := services.UserService{}
				realToken, err := rr.GetIdToken(claims.ID)
				if err != nil {
					// 没有登录或者已经过期了
					code = http.StatusUnauthorized
				} else if realToken != token {
					// 在其他地方登录了
					code = http.StatusUnauthorized
				} else {
					//rr.UpdateIdToken(claims.ID)
					c.Set("id", claims.ID)
					c.Set("account", claims.Account)
				}
			}
		}

		if code != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "请登录后再进行操作",
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
