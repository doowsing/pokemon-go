package ginapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pokemon/pkg/services"
	"pokemon/pkg/utils"
	"sync"
)

type Gapp struct {
	C       *gin.Context
	OptSrv  *services.OptService
	session utils.Session
}

var gappPool = &sync.Pool{
	New: func() interface{} {
		gapp := &Gapp{OptSrv: services.NewOptService()}
		return gapp
	},
}

func NewGapp(c *gin.Context) *Gapp {
	gapp := gappPool.Get().(*Gapp)
	//gapp := &Gapp{OptSrv: services.NewOptService()}
	gapp.C = c
	gapp.session = nil
	gapp.OptSrv.ReSet()
	//gapp.OptSrv = services.NewOptService()

	return gapp
}

func DropGapp(gapp *Gapp) {
	//services.DropOptService(gapp.OptSrv)
	//gapp.OptSrv = nil
	gappPool.Put(gapp)
}

func NewGappHandler(gHandler func(gapp *Gapp)) func(*gin.Context) {

	return func(c *gin.Context) {
		gapp := NewGapp(c)
		defer DropGapp(gapp)
		gHandler(gapp)
	}
}

func (g *Gapp) Session() utils.Session {
	if g.session == nil {
		g.session = utils.GetSession(g.C)
	}
	return g.session
}

func (g *Gapp) Redirect(location string) {
	g.C.Redirect(http.StatusFound, location)
}

func (g *Gapp) String(format string, values ...interface{}) {
	g.C.Header("Content-Type", fmt.Sprintf("text/html; charset=utf-8"))
	g.C.String(http.StatusOK, format, values...)
}

func (g *Gapp) HTML(name string, obj interface{}) {

	g.C.Render(http.StatusOK, NewHTML(name, obj))
}

func (g *Gapp) JSON(obj interface{}) {
	g.C.JSON(http.StatusOK, obj)
}

func (g *Gapp) JSONDATA(code int, msg string, obj interface{}) {
	g.C.JSON(http.StatusOK, gin.H{"code": code, "msg": msg, "data": obj})
}
func (g *Gapp) JSONDATAOK(msg string, obj interface{}) {
	g.JSONDATA(200, msg, obj)
}

func (g *Gapp) Id() int {
	return g.C.MustGet("id").(int)
}

func (g *Gapp) OkJson(data interface{}) {
	g.C.JSON(http.StatusOK, NewJSONData(http.StatusOK, "", data))
}

func NewJSONData(code int, msg string, data interface{}) gin.H {
	return gin.H{"Code": code, "Message": msg, "Data": data}
}
