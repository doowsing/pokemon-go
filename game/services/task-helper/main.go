package task_helper

import (
	"fmt"
	"github.com/psampaz/slice"
	"github.com/unknwon/com"
	"pokemon/common/persistence"
	"pokemon/game/models"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 一个玩家有多个任务，一个任务可能有多个要求
var TaskHelperInstance = NewTaskHelper()

type TaskHelper struct {
	UserTasks map[int]*UserTask
	mu        sync.RWMutex
}

func NewTaskHelper() *TaskHelper {
	return &TaskHelper{
		UserTasks: make(map[int]*UserTask),
	}
}

func (th *TaskHelper) UpdateTaskState(userId, gpcId int) {
	userTask, ok := th.GetUserTask(userId)
	if !ok {
		userTask = th.UpdateUserTask(userId)
	}
	userTask.PushGpc(gpcId)

}

func (th *TaskHelper) UpdateUserTask(userId int) *UserTask {
	userTask, ok := th.GetUserTask(userId)
	if !ok {
		userTask = NewUserTask()
		tasks := []*models.UTask{}
		persistence.GetOrm().Where("uid=?", userId).Find(&tasks)
		for _, t := range tasks {
			taskInfo := NewTaskInfo(t)
			if taskInfo != nil {
				userTask.Append(taskInfo)
			}
		}
	}
	th.SetUserTask(userId, userTask)
	return userTask
}

func (th *TaskHelper) GetUserTask(userId int) (*UserTask, bool) {
	th.mu.RLock()
	handler, ok := th.UserTasks[userId]
	th.mu.RUnlock()
	return handler, ok
}

func (th *TaskHelper) SetUserTask(userId int, userTask *UserTask) {
	th.mu.Lock()
	th.UserTasks[userId] = userTask
	th.mu.Unlock()
}

func (th *TaskHelper) ClearUnActive() {
	now := int(time.Now().Unix())
	deleteList := []int{}
	th.mu.RLock()
	for userId, userTask := range th.UserTasks {
		if now-userTask.lastPushTime > 60*10 {
			deleteList = append(deleteList, userId)
		}
	}
	th.mu.RUnlock()

	th.mu.Lock()
	for _, id := range deleteList {
		delete(th.UserTasks, id)
	}
	th.mu.Unlock()
}

type UserTask struct {
	lastPushTime int
	taskList     []*TaskInfo
}

func NewUserTask() *UserTask {
	return &UserTask{}
}

func (ut *UserTask) Append(taskInfo *TaskInfo) {
	ut.taskList = append(ut.taskList, taskInfo)
}

func (ut *UserTask) Delete(index int) {
	ut.taskList = append(ut.taskList[:index], ut.taskList[index+1:]...)
}

func (ut *UserTask) PushGpc(gpcId int) {
	if gpcId == 0 {
		return
	}
	ut.lastPushTime = int(time.Now().Unix())
	gpcIdStr := strconv.Itoa(gpcId)
	for i, taskInfo := range ut.taskList {
		for _, needInfo := range taskInfo.taskNeedInfos {
			if needInfo.aimNum <= needInfo.nowNum {
				continue
			}
			if slice.ContainsString(needInfo.gpcList, gpcIdStr) {
				// find
				needInfo.nowNum++

				// save task
				taskInfo.SaveTaskDb()

				if taskInfo.CheckFinish() {
					// task finished, delete this
					ut.Delete(i)
				}
				return
			}
		}
	}
}

// 单个任务，多个要求
type TaskInfo struct {
	taskId        int
	otherState    []string
	taskNeedInfos []*TaskNeedInfo
}

func NewTaskInfo(task *models.UTask) *TaskInfo {
	taskInfo := &TaskInfo{}
	otherState := []string{}
	taskNeedInfos := []*TaskNeedInfo{}
	task.GetM()
	if strings.Contains(task.MModel.OkNeed, "killmon") {
		for _, need := range strings.Split(task.MModel.OkNeed, ",") {
			if strings.Contains(need, "killmon") {
				items := strings.Split(need, ":")
				taskNeedInfo := &TaskNeedInfo{}
				taskNeedInfo.gpcIdStr = items[1]
				taskNeedInfo.gpcList = strings.Split(items[1], "|")
				taskNeedInfo.aimNum = com.StrTo(items[2]).MustInt()
				taskNeedInfos = append(taskNeedInfos, taskNeedInfo)
			}
		}
	}
	if len(taskNeedInfos) > 0 {
		for _, state := range strings.Split(task.State, ",") {
			if strings.Contains(state, "killmon") {
				items := strings.Split(state, ":")
				for _, needInfo := range taskNeedInfos {
					if needInfo.gpcIdStr == items[1] {
						needInfo.nowNum = com.StrTo(items[2]).MustInt()
					}
				}
			} else {
				otherState = append(otherState, state)
			}
		}
	} else {
		// 不需要杀怪的任务不需要监控
		return nil
	}

	taskInfo.taskId = task.Id
	taskInfo.otherState = otherState
	taskInfo.taskNeedInfos = taskNeedInfos
	if taskInfo.CheckFinish() {
		// 已完成的任务不需要监控
		return nil
	}

	return taskInfo
}

func (ti *TaskInfo) CheckFinish() bool {
	for _, needInfo := range ti.taskNeedInfos {
		if needInfo.nowNum < needInfo.aimNum {
			return false
		}
	}
	return true
}

func (ti *TaskInfo) SaveTaskDb() {
	var states = []string{}
	states = append(states, ti.otherState...)
	for _, needInfo := range ti.taskNeedInfos {
		states = append(states, needInfo.StateString())
	}
	persistence.GetOrm().Model(&models.UTask{Id: ti.taskId}).Update(map[string]interface{}{"state": strings.Join(states, ",")})
}

// 任务要求
type TaskNeedInfo struct {
	gpcIdStr string
	gpcList  []string
	nowNum   int
	aimNum   int
}

func (tni *TaskNeedInfo) StateString() string {
	return fmt.Sprintf("killmon:%s:%d", tni.gpcIdStr, tni.nowNum)
}
