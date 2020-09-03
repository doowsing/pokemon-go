package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pokemon/pkg/common"
	"pokemon/pkg/ginapp"
	"pokemon/pkg/models"
	"pokemon/pkg/services"
	"pokemon/pkg/utils"
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

func (uc *SysController) CreateUser(c *gin.Context) {
	us := services.NewUserService(nil)
	if us.InitCreate() {
		c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "创建用户成功！"), Data: nil})
	} else {
		c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "创建用户失败！"), Data: nil})
	}
}

func (uc *SysController) InitUserTable(c *gin.Context) {
	us := services.NewUserService(nil)
	if err := us.InitTable(); err != nil {
		fmt.Print(err)
		c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "异常！"), Data: nil})
	} else {
		c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "成功！"), Data: nil})
	}
}
func (uc *SysController) GetTime(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.String(fmt.Sprintf("%d", utils.NowUnix()))
}
