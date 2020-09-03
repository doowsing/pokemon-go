package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pokemon/common/rcache"
	"pokemon/game/services/common"
)

func TokenMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(common.TokenKey)
		//fmt.Printf("tokenkey:%s\n", token)
		login := false
		if token != "" {
			info, err := common.GetVerifyInfo(token)
			//fmt.Printf("id:%s\n", info.Id)
			//fmt.Printf("rand:%s\n", info.Token)
			if err == nil {
				realToken, err := rcache.GetIdToken(info.Id)
				if err == nil {
					if realToken == token {
						c.Set("id", info.Id)
						c.Set("account", info.Account)
						login = true
					}
				}
				common.DropVerifyInfo(info)
			}

		}
		if !login {
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "请登录后再进行操作",
				"data": nil,
			})

			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
