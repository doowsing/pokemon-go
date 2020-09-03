package controllers

import (
	"github.com/gin-gonic/gin"
	"pokemon/game/ginapp"
	"strconv"
)

// 已接受任务列表
func AcceptTaskList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	taskData := gapp.OptSvc.TaskSrv.GetAcceptTaskData(gapp.Id())

	gapp.JSONDATAOK("", taskData)
}

// 可接受任务列表
func EnableAcceptTaskList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	typeStr := c.Query("type")
	typeId, err := strconv.Atoi(typeStr)
	if err != nil || typeId < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	taskData := gapp.OptSvc.TaskSrv.GetEnableAcceptTaskData(gapp.Id(), typeId)

	gapp.JSONDATAOK("", taskData)
}

// 任务展示-未接受的任务
func TaskInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	gapp.JSONDATAOK("", gapp.OptSvc.TaskSrv.GetTaskInfo(gapp.Id(), id))

}

// 任务展示-已接受的任务
func AcceptTaskInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	gapp.JSONDATAOK("", gapp.OptSvc.TaskSrv.GetUserTaskInfo(gapp.Id(), id))

}

// 接受任务
func AcceptTask(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	ok, msg := gapp.OptSvc.TaskSrv.AcceptTask(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})

}

// 完成任务
func FinishTask(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	ok, msg := gapp.OptSvc.TaskSrv.FinishTask(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})

}

// 放弃任务
func ThrowTask(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数错误！", nil)
		return
	}
	ok, msg := gapp.OptSvc.TaskSrv.ThrowTask(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})

}
