package middleware

import (
	"github.com/gin-gonic/gin"
	"pokemon/pkg/rcache"
)

//修改jwt源码 jwt.go defaultCheckJWT() 如果是sessions 登录登出接口 不做处理
// 不如此处理的话. dotweb的中间件调用貌似有点问题. /api/的中间件会影响所有的中间件
type CrosMiddleware struct {
}

func (cm *CrosMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := new(rcache.SysRedisRepository)
		repo.InsertIp(c.GetHeader("X-Real-IP"))
		//c.Header("Access-Control-Allow-Origin", "*")
		//c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,Sign")
		//c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	}
	//插入ip统计

	//ctx.Response().SetHeader("Content-Type", "application/json")             //返回数据格式是json
}
func NewCROSMiddleware() *CrosMiddleware {
	return &CrosMiddleware{}
}
