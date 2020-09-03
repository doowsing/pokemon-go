package controllers

import (
	"github.com/gin-gonic/gin"
	"pokemon/game/ginapp"
)

type ErrorController struct {
}

func NewErrorController() *ErrorController {
	return &ErrorController{}
}

func (*ErrorController) ServerError(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	c.JSON(500, gin.H{"code": 500, "msg": "服务器错误！", "data": nil})
}

func (*ErrorController) PageNotFound(c *gin.Context) {
	c.JSON(404, gin.H{"code": 404, "msg": "Page not found", "data": nil})
}

func (*ErrorController) MethodNotFound(c *gin.Context) {
	c.JSON(405, gin.H{"code": "NOT_FOUND", "message": "Method not found"})
}
