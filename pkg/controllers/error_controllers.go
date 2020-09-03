package controllers

import (
	"github.com/gin-gonic/gin"
	"pokemon/pkg/ginapp"
)

type ErrorController struct {
}

func NewErrorController() *ErrorController {
	return &ErrorController{}
}

func (*ErrorController) ServerError(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.C.String(500, "服务器错误！")
}

func (*ErrorController) PageNotFound(c *gin.Context) {
	c.JSON(404, gin.H{"code": "NOT_FOUND", "message": "Page not found"})
}

func (*ErrorController) MethodNotFound(c *gin.Context) {
	c.JSON(405, gin.H{"code": "NOT_FOUND", "message": "Method not found"})
}
