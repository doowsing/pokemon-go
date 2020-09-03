package ginapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pokemon/game/services"
	"sync"
)

type Gapp struct {
	C      *gin.Context
	OptSvc *services.OptService
}

var gappPool = &sync.Pool{
	New: func() interface{} {
		gapp := &Gapp{OptSvc: services.NewOptService()}
		return gapp
	},
}

func NewGapp(c *gin.Context) *Gapp {
	gapp := gappPool.Get().(*Gapp)
	//gapp := &Gapp{OptSvc: services.NewOptService()}
	gapp.C = c
	gapp.OptSvc.ReSet()
	//gapp.OptSvc = services.NewOptService()

	return gapp
}

func DropGapp(gapp *Gapp) {
	//services.DropOptService(gapp.OptSvc)
	//gapp.OptSvc = nil
	gappPool.Put(gapp)
}

func NewGappHandler(gHandler func(gapp *Gapp)) func(*gin.Context) {

	return func(c *gin.Context) {
		gapp := NewGapp(c)
		defer DropGapp(gapp)
		gHandler(gapp)
	}
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

func (g *Gapp) Account() string {
	return g.C.MustGet("account").(string)
}

func (g *Gapp) OkJson(data interface{}) {
	g.C.JSON(http.StatusOK, NewJSONData(http.StatusOK, "", data))
}

func NewJSONData(code int, msg string, data interface{}) gin.H {
	return gin.H{"Code": code, "Message": msg, "Data": data}
}
