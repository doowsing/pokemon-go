package loginmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
	"pokemon/common/config"
	"pokemon/common/persistence"
	"pokemon/game/services"
	"pokemon/game/utils"
	"pokemon/game/utils/redistore"
)

var store = getRedisStore("secret")

var UnLoginUrl = []string{}

func getRedisStore(name string) sessions.Store {
	store, _ := redistore.NewRediStoreWithPool(redistore.NewRedisClusterPool(persistence.GetRedisCluster()), []byte(name))
	store.SetMaxAge(86400)
	return store
}

func SessionMdl() gin.HandlerFunc {
	return utils.Sessions("ptoken", store, config.Config().SessionExpire)
}

func CharSet(CharSet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", fmt.Sprintf("text/html; charset=%s", CharSet))
	}
}

func LoginPage() gin.HandlerFunc {
	return loginRequeire("page")
}

func LoginRequest() gin.HandlerFunc {
	return loginRequeire("request")
}

func loginRequeire(RType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		needLogin := true
		for _, v := range UnLoginUrl {
			if v == c.Request.URL.Path {
				needLogin = false
				break
			}
		}
		if needLogin {
			abort := true
			session := utils.GetSession(c)
			if id := GetUserId(session); id > 0 {

				rr := services.UserService{}
				realToken, err := rr.GetIdToken(id)
				if err != nil || realToken == session.SessionId() {

					//rr.UpdateIdToken(claims.ID)
					abort = false
					c.Set("id", id)
				}

			}
			if abort {
				//if RType == "page" {
				//	gapp.Redirect("/passport/login.php")
				//} else {
				//	gapp.String(`<a href="/passport/login.php?rd=%s" target=_top>网络中断，请重新登录!!</a>`, strconv.Itoa(rand.Int()))
				//}
				c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "请登录后再操作", "data": nil})

				c.Abort()
			}
		}
		c.Next()
	}
}

func GetUserId(s utils.Session) int {
	idStr, ok := s.Get("id")
	if !ok {
		return -1
	}
	return idStr.(int)
}
