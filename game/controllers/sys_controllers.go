package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pokemon/game/ginapp"
	"pokemon/game/services"
	"pokemon/game/utils"
)

var SysCtl = NewSysController()

type SysController struct {
	service *services.SysService
}

func NewSysController() *SysController {
	return &SysController{service: &services.SysService{}}
}

func (sc *SysController) InitRdModels(c *gin.Context) {
	success := sc.service.InitRdModels()
	c.JSON(200, gin.H{"code": 200, "msg": "设置redis中宠物、技能、经验表原型成功！", "data": success})
}

func (uc *SysController) GetTime(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.String(fmt.Sprintf("%d", utils.NowUnix()))
}
