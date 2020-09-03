package controllers

import (
	"github.com/gin-gonic/gin"
	"pokemon/pkg/ginapp"
	"pokemon/pkg/models"
)

// 接受任务列表
func GainTaskList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	//tasks := gapp.OptSrv.TaskSrv.GetGainTaskList(gapp.Id())
	//taskData := []gin.H{}
	//for _, t := range tasks {
	//
	//}
}

func CheckTaskEnableComplite(task *models.UTask) bool {
	// 检查背包

	// 检查要求

	// 检查限制
	return false
}
