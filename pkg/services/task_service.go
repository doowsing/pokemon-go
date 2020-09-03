package services

import "pokemon/pkg/models"

type TaskServices struct {
	baseService
}

func NewTaskServices(osrc *OptService) *TaskServices {
	us := &TaskServices{}
	us.SetOptSrc(osrc)
	return us
}

func (ts *TaskServices) GetGainTaskList(userId int) []*models.UTask {
	tasks := []*models.UTask{}
	ts.GetDb().Where("uid=?", userId).Find(tasks)
	return tasks
}
