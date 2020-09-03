package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"pokemon/common/rcache"
	"pokemon/game/utils"
	"time"
)

func LimitIp() gin.HandlerFunc {
	return func(context *gin.Context) {
		id := com.StrTo(context.MustGet("id").(int)).MustInt()
		if id <= 0 {
			return
		}
		ip := utils.GetClientIp(context)
		users := rcache.GetIPUsers(ip)
		if users == nil {
			users = make(map[int]int)
		}
		find := false
		for i, _ := range users {
			if i == id {
				find = true
				break
			}
		}
		if !find && len(users) >= rcache.IpLimitCount {
			context.JSON(http.StatusOK, gin.H{"code": 401, "msg": "当前IP活跃账号过多，强制退出登录状态！", "data": nil})
			rcache.DelIdToken(id)
			context.Abort()
			return
		}
		now := int(time.Now().Unix())
		users[id] = now
		rcache.SetIPUsers(ip, users)

	}
}
