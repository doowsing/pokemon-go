package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
	"math/rand"
	"pokemon/game/models"
	"pokemon/game/services/common"
	task_helper "pokemon/game/services/task-helper"
	"pokemon/game/utils"
	"strconv"
	"strings"
)

type TaskServices struct {
	BaseService
}

func NewTaskServices(osrc *OptService) *TaskServices {
	us := &TaskServices{}
	us.SetOptSrc(osrc)
	return us
}

func (ts *TaskServices) GetAcceptTaskData(userId int) []gin.H {
	tasks := ts.GetAcceptTask(userId)
	//bagCnt := ts.OptSvc.PropSrv.GetCarryPropsCnt(userId)
	user := ts.OptSvc.UserSrv.GetUserById(userId)
	userInfo := ts.OptSvc.UserSrv.GetUserInfoById(userId)
	carryProp := ts.OptSvc.PropSrv.GetCarryProps(userId, false)
	mainPet := ts.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()
	taskLogs := []*models.TaskLog{}
	ts.GetDb().Where("uid=?", userId).Find(&taskLogs)
	taskDatas := []gin.H{}
	for _, t := range tasks {
		t.GetM()
		data := gin.H{}
		data["id"] = t.Id
		data["title"] = t.MModel.Title
		enableFinish := true

		// 任务要求
		enableFinish, _, _ = ts.GetAimProcessData(t, user, mainPet, carryProp)
		// 是否可重复完成以及是否为系列任务
		if enableFinish {
			enableFinish = ts.PassCid(t.MModel, taskLogs)
		}

		// 限制要求
		if enableFinish {
			enableFinish, _ = ts.PassLimit(t.MModel, user, userInfo, mainPet)
		}

		data["enable_finish"] = enableFinish
		taskDatas = append(taskDatas, data)
	}
	return taskDatas
}

func (ts *TaskServices) GetEnableAcceptTaskData(userId, typeId int) []gin.H {
	acceptTasks := ts.GetAcceptTask(userId)
	taskLogs := []*models.TaskLog{}
	ts.GetDb().Where("uid=?", userId).Find(&taskLogs)
	tasks := []*models.Task{}
	ts.GetDb().Where("color=?", typeId).Find(&tasks)

	user := ts.OptSvc.UserSrv.GetUserById(userId)
	userInfo := ts.OptSvc.UserSrv.GetUserInfoById(userId)
	mainPet := ts.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()

	taskDatas := []gin.H{}
	for _, t := range tasks {
		enableAccept := true

		// 是否隐藏
		if t.Hide == 0 || t.Hide == 2 {
			continue
		}
		if !ts.PassCid(t, taskLogs) {
			continue
		}

		// 是否正在接受中
		find := false
		for _, atask := range acceptTasks {
			if atask.TaskId == t.Id {
				find = true
				break
			}
		}
		if find {
			// 已经领取了
			continue
		}
		// 检查是否符合接受条件
		enableAccept, _ = ts.PassLimit(t, user, userInfo, mainPet)

		taskDatas = append(taskDatas, gin.H{
			"id":            t.Id,
			"title":         t.Title,
			"enbale_accept": enableAccept,
		})
	}
	return taskDatas
}

func (ts *TaskServices) GetTask(userId, taskId int) *models.UTask {
	task := &models.UTask{}
	ts.GetDb().Where("uid=? and id=?", userId, taskId).First(task)
	if task.Id == 0 {
		return nil
	}

	return task
}

func (ts *TaskServices) GetTaskInfo(userId, taskId int) gin.H {
	task := common.GetTask(taskId)
	if task == nil {
		return nil
	}
	if task.Hide == 0 || task.Hide == 2 {
		return nil
	}

	acceptTasks := ts.GetAcceptTask(userId)
	taskLogs := []*models.TaskLog{}
	ts.GetDb().Where("uid=?", userId).Find(&taskLogs)

	user := ts.OptSvc.UserSrv.GetUserById(userId)
	userInfo := ts.OptSvc.UserSrv.GetUserInfoById(userId)
	mainPet := ts.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()
	enableAccept := true

	enableAccept = ts.PassCid(task, taskLogs)

	// 是否正在接受中
	find := false
	for _, atask := range acceptTasks {
		if atask.TaskId == task.Id {
			find = true
			break
		}
	}
	if find {
		// 已经领取了
		enableAccept = false
	}
	// 检查是否符合接受条件

	if enableAccept {
		enableAccept, _ = ts.PassLimit(task, user, userInfo, mainPet)
	}

	aimDatas := ts.GetAimData(task)

	awardDatas := ts.GetAwardData(task)

	return gin.H{
		"id":            task.Id,
		"title":         task.Title,
		"description":   task.FromMsg,
		"aim":           aimDatas,
		"award":         awardDatas,
		"enable_accept": enableAccept,
	}
}

func (ts *TaskServices) GetUserTaskInfo(userId, utaskId int) gin.H {
	task := ts.GetTask(userId, utaskId)
	if task == nil {
		return nil
	}
	task.GetM()
	if task.MModel.Hide == 0 || task.MModel.Hide == 2 {
		return nil
	}

	//acceptTasks := ts.GetAcceptTask(userId)
	taskLogs := []*models.TaskLog{}
	ts.GetDb().Where("uid=?", userId).Find(&taskLogs)

	user := ts.OptSvc.UserSrv.GetUserById(userId)
	userInfo := ts.OptSvc.UserSrv.GetUserInfoById(userId)
	mainPet := ts.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()
	enableFinish := false

	stateMap := make(map[string]int)
	for _, state := range strings.Split(task.State, ",") {
		sItems := strings.Split(state, ":")
		if len(sItems) > 2 && sItems[0] == "killmon" {
			stateMap[sItems[1]] = com.StrTo(sItems[2]).MustInt()
		}
	}
	// 检查是否符合接受条件
	enableFinish = ts.PassCid(task.MModel, taskLogs)
	if enableFinish {
		enableFinish, _ = ts.PassLimit(task.MModel, user, userInfo, mainPet)
	}

	awardDatas := ts.GetAwardData(task.MModel)
	carryProps := ts.OptSvc.PropSrv.GetCarryProps(userId, false)
	_enableFinish, aimDatas, processDatas := ts.GetAimProcessData(task, user, mainPet, carryProps)
	if enableFinish {
		enableFinish = _enableFinish
	}

	return gin.H{
		"id":            task.Id,
		"title":         task.MModel.Title,
		"description":   task.MModel.OkMsg,
		"aim":           aimDatas,
		"process":       processDatas,
		"award":         awardDatas,
		"enable_finish": enableFinish,
	}

}

func (ts *TaskServices) GetAcceptTask(userId int) []*models.UTask {
	tasks := []*models.UTask{}
	ts.GetDb().Where("uid=?", userId).Find(&tasks)
	return tasks
}

func (ts *TaskServices) PassLimit(task *models.Task, user *models.User, userInfo *models.UserInfo, mainPet *models.UPet) (bool, string) {
	enablePass := true
	msg := ""
	for _, limit := range strings.Split(task.LimitLv, ",") {
		lItems := strings.Split(limit, ":")
		if len(limit) > 1 {
			switch lItems[0] {
			case "level", "lv":
				lvItems := strings.Split(lItems[1], "|")
				if len(lvItems) == 1 {
					if mainPet.Level < com.StrTo(lvItems[0]).MustInt() {
						enablePass = false
						msg = "主宠等级不足！"
					}
				} else if len(lvItems) == 2 {
					if mainPet.Level < com.StrTo(lvItems[0]).MustInt() || (com.StrTo(lvItems[1]).MustInt() > 0 && mainPet.Level > com.StrTo(lvItems[1]).MustInt()) {
						enablePass = false
						msg = "主宠等级不符合要求！"
					}

				}
				break
			case "czl":
				czlItems := strings.Split(lItems[1], "|")
				if len(czlItems) == 1 {
					if mainPet.CC < com.StrTo(czlItems[0]).MustFloat64() {
						enablePass = false
						msg = "主宠成长不足！"
					}
				} else if len(czlItems) == 2 {
					if mainPet.CC < com.StrTo(czlItems[0]).MustFloat64() || mainPet.CC > com.StrTo(czlItems[1]).MustFloat64() {
						enablePass = false
						msg = "主宠成长不符合要求！"
					}
				}
				break
			case "comself":
				if !com.IsSliceContainsStr(strings.Split(lItems[1], "|"), strconv.Itoa(mainPet.MModel.ID)) {
					enablePass = false
					msg = "主宠不符合要求！"
				}
				break
			case "jifen", "jf":
				if user.Score < com.StrTo(lItems[1]).MustInt() {
					enablePass = false
					msg = "积分不足！"
				}
				break
			case "wx":
				if mainPet.MModel.Wx != com.StrTo(lItems[1]).MustInt() {
					enablePass = false
					msg = "主宠五行不符合要求！"
				}
				break
			case "cishu":
				// 限制次数
				cplCnt := 0
				ts.GetDb().Model(&models.TaskLog{}).Where("uid=? and time>? and taskid=?", user.ID, utils.ToDayStartUnix(), task.Id).Count(&cplCnt)
				if cplCnt >= com.StrTo(lItems[1]).MustInt() {
					enablePass = false
					msg = fmt.Sprintf("每天只能完成 %s 次任务", lItems[1])
				}
				break
			case "timelimit":
				//if utils.NowUnix()-task.Time-com.StrTo(items[1]).MustInt()*3600 > 0 {
				//				//	// 超时了
				//				//	enablePass = false
				//				//}
				break
			case "xfsj":
				enablePass = false
				break
			case "paihang":
				if user.Paihang != com.StrTo(lItems[1]).MustInt() {
					// 排行不对
					enablePass = false
					msg = "排行不符合要求！"
				}
				break
			case "all_rmb":
				if user.AllRmb < com.StrTo(lItems[1]).MustInt() {
					enablePass = false
					msg = "累计充值不符合要求！"
				}
				break
			case "vip":
				if user.Vip < com.StrTo(lItems[1]).MustInt() {
					enablePass = false
					msg = "VIP积分不足！"
				}
				break
			case "merge":
				if userInfo.Merge == 0 {
					enablePass = false
					msg = "需要结婚后才可完成任务！"
				}
				break
			}
			if !enablePass {
				return enablePass, msg
			}
		}
	}
	return enablePass, msg
}

func (ts *TaskServices) PassCid(task *models.Task, taskLogs []*models.TaskLog) bool {
	enablePass := true
	if task.Cid == "0" {
		// 检查是否完成过
		find := false
		for _, tlog := range taskLogs {
			if tlog.TaskId == task.Id {
				find = true
				break
			}
		}
		if find {
			// 只能完成一次且已经完成过了
			enablePass = false
		}

	} else if strings.Index(task.Cid, "rwl") > -1 {
		cidItems := strings.Split(task.Cid, ":")
		if cidItems[0] == "rwl" && len(cidItems) > 1 {
			items := strings.Split(cidItems[1], "|")
			if task.Xulie > 0 {
				find := false
				for _, tlog := range taskLogs {
					if tlog.Xulie == task.Xulie {
						if tlog.TaskId == com.StrTo(items[0]).MustInt() {
							find = true
							break
						}
					}
				}
				if !find && task.Id != com.StrTo(items[0]).MustInt() {
					// 没完成过前置任务且前置任务不为自身
					enablePass = false
				}
			} else {
				if task.Id == com.StrTo(items[0]).MustInt() {
					// 下架了
					enablePass = false
				}
			}

		}
	}
	return enablePass
}

func (ts *TaskServices) GetAimData(task *models.Task) []string {
	aimDatas := []string{}
	for _, a := range strings.Split(task.OkNeed, ",") {
		needItems := strings.Split(a, ":")
		switch needItems[0] {
		case "wx":
			aimDatas = append(aimDatas, fmt.Sprintf("需要五系：%s 系", utils.GetWxName(com.StrTo(needItems[1]).MustInt())))
			break
		case "see":
			break
		case "killmon":
			if len(needItems) > 2 {
				gids := strings.Split(needItems[1], "|")
				num := com.StrTo(needItems[2]).MustInt()
				gpcMap := make(map[string]bool)
				for _, gid := range gids {
					g := common.GetGpc(com.StrTo(gid).MustInt())
					if g != nil {
						gpcMap[g.Name] = true
					}
				}
				if len(gpcMap) > 0 {
					nameSet := []string{}
					for name, _ := range gpcMap {
						nameSet = append(nameSet, name)
					}
					aimDatas = append(aimDatas, fmt.Sprintf("击败怪物：%s %d 个", strings.Join(nameSet, "、"), num))
				}
			}
			break
		case "giveitem":
			if len(needItems) > 2 {
				pid := com.StrTo(needItems[1]).MustInt()
				num := com.StrTo(needItems[2]).MustInt()
				p := common.GetMProp(pid)
				if p != nil {
					aimDatas = append(aimDatas, fmt.Sprintf("收集：%s %d 个", p.Name, num))
				}
			}
			break
		case "money":
			aimDatas = append(aimDatas, fmt.Sprintf("需要金币：%d 个", com.StrTo(needItems[1]).MustInt()))
			break
		case "ww":
			aimDatas = append(aimDatas, fmt.Sprintf("需要威望：%d 点", com.StrTo(needItems[1]).MustInt()))
			break
		case "jifen":
			aimDatas = append(aimDatas, fmt.Sprintf("需要积分：%d 点", com.StrTo(needItems[1]).MustInt()))
			break
		case "dianjuan":
			aimDatas = append(aimDatas, fmt.Sprintf("需要交纳点券：%d 点", com.StrTo(needItems[1]).MustInt()))
			break
		case "all_rmb":
			aimDatas = append(aimDatas, fmt.Sprintf("需要累计充值：%d 元", com.StrTo(needItems[1]).MustInt()))
			break
		case "lv":
			lvItems := strings.Split(needItems[1], "|")
			var lvStr string
			if len(lvItems) == 1 || lvItems[1] == "0" {
				lvStr = lvItems[0]
			} else {
				lvStr = strings.Join(lvItems, "-")
			}
			aimDatas = append(aimDatas, fmt.Sprintf("需要主宠等级：%s 级", lvStr))
			break
		}
	}
	return aimDatas
}

func (ts *TaskServices) GetAimProcessData(utask *models.UTask, user *models.User, mainPet *models.UPet, carryProps []*models.UProp) (bool, []string, []string) {
	aimDatas := []string{}
	processDatas := []string{}
	task := utask.GetM()
	mainPet.GetM()
	stateMap := make(map[string]int)
	for _, state := range strings.Split(utask.State, ",") {
		sItems := strings.Split(state, ":")
		if len(sItems) > 2 && sItems[0] == "killmon" {
			stateMap[sItems[1]] = com.StrTo(sItems[2]).MustInt()
		}
	}
	enableFinish := true
	for _, a := range strings.Split(task.OkNeed, ",") {
		needItems := strings.Split(a, ":")
		process := "未完成"
		switch needItems[0] {
		case "wx":
			wx := com.StrTo(needItems[1]).MustInt()
			if wx == mainPet.MModel.Wx {
				process = "已完成"
			} else {
				enableFinish = false
			}
			aimDatas = append(aimDatas, fmt.Sprintf("需要五系：%s 系", utils.GetWxName(wx)))
			processDatas = append(processDatas, fmt.Sprintf("主宠五系：%s 系 %s", utils.GetWxName(mainPet.MModel.Wx), process))
			break
		case "see":
			break
		case "killmon":
			if len(needItems) > 2 {
				gids := strings.Split(needItems[1], "|")
				num := com.StrTo(needItems[2]).MustInt()
				gpcMap := make(map[string]bool)
				for _, gid := range gids {
					g := common.GetGpc(com.StrTo(gid).MustInt())
					if g != nil {
						gpcMap[g.Name] = true
					}
				}
				if len(gpcMap) > 0 {
					nameSet := []string{}
					for name, _ := range gpcMap {
						nameSet = append(nameSet, name)
					}
					cplNum, _ := stateMap[needItems[1]]

					aimDatas = append(aimDatas, fmt.Sprintf("击败怪物：%s %d 个", strings.Join(nameSet, "、"), num))
					if cplNum >= num {
						process = "已完成"
					} else {
						enableFinish = false
					}
					processDatas = append(processDatas, fmt.Sprintf("已击败 %s %d/%d %s", strings.Join(nameSet, "、"), cplNum, num, process))

				}
			}
			break
		case "giveitem":
			if len(needItems) > 2 {
				pid := com.StrTo(needItems[1]).MustInt()
				num := com.StrTo(needItems[2]).MustInt()
				p := common.GetMProp(pid)
				if p != nil {
					cplNum := 0
					for _, cp := range carryProps {
						cp.GetM()
						if cp.Pid == pid {
							cplNum = cp.Sums
						}
					}
					aimDatas = append(aimDatas, fmt.Sprintf("收集：%s %d 个", p.Name, num))
					if cplNum >= num {
						process = "已完成"
					} else {
						enableFinish = false
					}
					processDatas = append(processDatas, fmt.Sprintf("已收集 %s %d/%d %s", p.Name, cplNum, num, process))

				}
			}
			break
		case "money":
			num := com.StrTo(needItems[1]).MustInt()
			aimDatas = append(aimDatas, fmt.Sprintf("需要金币：%d 个", num))
			if user.Money >= num {
				process = "已完成"
			} else {
				enableFinish = false
			}
			processDatas = append(processDatas, fmt.Sprintf("金币 %d/%d %s", user.Money, num, process))

			break
		case "ww":
			num := com.StrTo(needItems[1]).MustInt()
			aimDatas = append(aimDatas, fmt.Sprintf("需要威望：%d 点", num))
			if user.Prestige >= num {
				process = "已完成"
			} else {
				enableFinish = false
			}
			processDatas = append(processDatas, fmt.Sprintf("威望 %d/%d %s", user.Prestige, num, process))
			break
		case "jifen":
			num := com.StrTo(needItems[1]).MustInt()
			aimDatas = append(aimDatas, fmt.Sprintf("需要积分：%d 点", num))
			if user.Score >= num {
				process = "已完成"
			} else {
				enableFinish = false
			}
			processDatas = append(processDatas, fmt.Sprintf("积分 %d/%d %s", user.Score, num, process))
			break
		case "dianjuan":
			num := com.StrTo(needItems[1]).MustInt()
			aimDatas = append(aimDatas, fmt.Sprintf("需要交纳点券：%d 点", com.StrTo(needItems[1]).MustInt()))
			if user.ActiveScore >= num {
				process = "已完成"
			} else {
				enableFinish = false
			}
			processDatas = append(processDatas, fmt.Sprintf("点券 %d/%d %s", user.ActiveScore, num, process))
			break
		case "all_rmb":
			num := com.StrTo(needItems[1]).MustInt()
			aimDatas = append(aimDatas, fmt.Sprintf("需要累计充值：%d 元", com.StrTo(needItems[1]).MustInt()))
			if user.AllRmb >= num {
				process = "已完成"
			} else {
				enableFinish = false
			}
			processDatas = append(processDatas, fmt.Sprintf("累计充值 %d/%d %s", user.AllRmb, num, process))
			break
		case "lv":
			lvItems := strings.Split(needItems[1], "|")
			if len(lvItems) == 1 {
				lv := com.StrTo(lvItems[0]).MustInt()
				if mainPet.Level >= lv {
					process = "已完成"
				} else {
					enableFinish = false
				}
				aimDatas = append(aimDatas, fmt.Sprintf("需要主宠等级：%d 级", lv))
				processDatas = append(processDatas, fmt.Sprintf("主宠等级 %d %s", mainPet.Level, process))
			} else if len(lvItems) == 2 {
				lv1 := com.StrTo(lvItems[0]).MustInt()
				var lv2 int
				if len(lvItems) > 1 {
					lv2 = com.StrTo(lvItems[1]).MustInt()
				} else {
					lv2 = 0
				}

				if mainPet.Level >= lv1 && (lv2 == 0 || mainPet.Level <= lv2) {
					process = "已完成"
				} else {
					enableFinish = false
				}
				if lv2 > 0 {
					aimDatas = append(aimDatas, fmt.Sprintf("需要主宠等级：%d-%d 级", lv1, lv2))
				} else {
					aimDatas = append(aimDatas, fmt.Sprintf("需要主宠等级：%d 级", lv1))
				}

				processDatas = append(processDatas, fmt.Sprintf("主宠等级 %d %s", mainPet.Level, process))

			}
			break
		}
	}
	return enableFinish, aimDatas, processDatas
}

func (ts *TaskServices) GetAwardData(task *models.Task) []string {
	awardDatas := []string{}

	for _, rs := range strings.Split(task.Result, ",") {
		rItems := strings.Split(rs, ":")
		switch rItems[0] {
		case "item", "itemrand":
			awardDatas = append(awardDatas, "随机获得道具")
			break
		case "lvprops":
			// 不同等级获得不同道具，暂时不做
			break
		case "props":
			if len(rItems) > 2 {
				pid := com.StrTo(rItems[1]).MustInt()
				num := com.StrTo(rItems[2]).MustInt()
				p := common.GetMProp(pid)
				if p != nil {
					awardDatas = append(awardDatas, fmt.Sprintf("获得物品：%s %d个", p.Name, num))
				}
			}
			break
		case "bprops":
			if len(rItems) > 2 {
				pid := com.StrTo(rItems[1]).MustInt()
				num := com.StrTo(rItems[2]).MustInt()
				p := common.GetMProp(pid)
				if p != nil {
					awardDatas = append(awardDatas, fmt.Sprintf("获得可交易物品：%s %d个", p.Name, num))
				}
			}
			break
		case "exp":
			if len(rItems) > 3 {
				wwItems := strings.Split(rItems[2], "|")
				if len(wwItems) == 1 {
					awardDatas = append(awardDatas, fmt.Sprintf("当交纳的威望大于%d时获得经验：%d", com.StrTo(wwItems[0]).MustInt(), com.StrTo(rItems[1]).MustInt()))
				} else if len(wwItems) == 2 {
					awardDatas = append(awardDatas, fmt.Sprintf("当交纳的威望在%d-%d之间时获得经验：%d", com.StrTo(wwItems[0]).MustInt(), com.StrTo(wwItems[1]).MustInt(), com.StrTo(rItems[1]).MustInt()))
				}
			} else {
				awardDatas = append(awardDatas, fmt.Sprintf("获得经验：%d", com.StrTo(rItems[1]).MustInt()))
			}
			break
		case "mon", "money":
			if len(rItems) > 3 {
				wwItems := strings.Split(rItems[2], "|")
				if len(wwItems) == 1 {
					awardDatas = append(awardDatas, fmt.Sprintf("当交纳的威望大于%d时获得金币：%d 个", com.StrTo(wwItems[0]).MustInt(), com.StrTo(rItems[1]).MustInt()))
				} else if len(wwItems) == 2 {
					awardDatas = append(awardDatas, fmt.Sprintf("当交纳的威望在%d-%d之间时获得金币：%d 个", com.StrTo(wwItems[0]).MustInt(), com.StrTo(wwItems[1]).MustInt(), com.StrTo(rItems[1]).MustInt()))
				}
			} else {
				awardDatas = append(awardDatas, fmt.Sprintf("获得金币：%d 个", com.StrTo(rItems[1]).MustInt()))
			}
			break
		}
	}
	return awardDatas
}

func (ts *TaskServices) PropExist(carryProps *[]*models.UProp, pid, num int) bool {
	for _, p := range *carryProps {
		if p.Pid == pid && p.Sums >= num && p.Zbing == 0 {
			return true
		}
	}
	return false
}

func (ts *TaskServices) AcceptTask(userId, taskId int) (bool, string) {
	task := common.GetTask(taskId)
	if task == nil {
		return false, "任务不存在！"
	}
	if task.Hide == 0 || task.Hide == 2 {
		return false, "任务不存在！"
	}

	acceptTasks := ts.GetAcceptTask(userId)
	taskLogs := []*models.TaskLog{}
	ts.GetDb().Where("uid=?", userId).Find(&taskLogs)

	user := ts.OptSvc.UserSrv.GetUserById(userId)
	userInfo := ts.OptSvc.UserSrv.GetUserInfoById(userId)
	mainPet := ts.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()
	enableAccept := false

	enableAccept = ts.PassCid(task, taskLogs)

	// 是否正在接受中
	find := false
	for _, atask := range acceptTasks {
		if atask.TaskId == task.Id {
			find = true
			break
		}
	}
	if find {
		// 已经领取了
		return false, "该任务已领取！"
	}
	// 检查是否符合接受条件

	enableAccept, msg := ts.PassLimit(task, user, userInfo, mainPet)
	if enableAccept {
		if len(acceptTasks) >= 15 {
			return false, "您最多领取15个任务！"
		}
	} else {
		return false, msg
	}
	if strings.Contains(task.OkNeed, "killmon") {
		go task_helper.TaskHelperInstance.UpdateUserTask(userId)
	}
	newTask := &models.UTask{
		Uid:    userId,
		TaskId: taskId,
		Time:   utils.NowUnix(),
	}
	ts.GetDb().Create(newTask)
	return true, "成功领取任务！"
}

func (ts *TaskServices) FinishTask(userId, taskId int) (bool, string) {
	uTask := ts.GetTask(userId, taskId)
	if uTask == nil {
		return false, "任务不存在！"
	}
	task := uTask.GetM()

	// 是否可重复完成以及是否为系列任务
	cidItems := strings.Split(uTask.MModel.Cid, ":")
	if len(cidItems) > 1 {
		// 系列任务
		if cidItems[0] == "rwl" && len(cidItems) > 1 {
			cidItem := strings.Split(cidItems[1], "|")
			tid1 := com.StrTo(cidItem[0]).MustInt()
			//tid2 := com.StrTo(cidItem[1]).MustInt()
			if task.Xulie > 0 {
				if tid1 != uTask.TaskId {
					fTaskLog := &models.TaskLog{}
					ts.GetDb().Where("uid = ? and taskid = ?", userId, tid1).First(fTaskLog)
					if fTaskLog.Id == 0 {
						return false, "该任务的前置任务没有完成！"
					}
				}
			} else {
				if tid1 == uTask.TaskId {
					// 下架了
					return false, "该任务已过时！"
				}
			}
		}
	} else {
		// 只可以完成一次的任务
		if uTask.MModel.Cid == "0" {
			fTaskLog := &models.TaskLog{}
			ts.GetDb().Where("uid = ? and taskid = ?", userId, task.Id).First(fTaskLog)
			if fTaskLog.Id > 0 {
				return false, "该任务只能完成一次！"
			}
		}
	}

	user := ts.OptSvc.UserSrv.GetUserById(userId)
	userInfo := ts.OptSvc.UserSrv.GetUserInfoById(userId)

	bagCnt := ts.OptSvc.PropSrv.GetCarryPropsCnt(userId)
	if needPlace := strings.Count(task.Result, "itemrand") + strings.Count(task.Result, "props"); needPlace+bagCnt > user.BagPlace {
		return false, fmt.Sprintf("请至少准备好 %d 的背包剩余空间", needPlace)
	}
	carryProp := ts.OptSvc.PropSrv.GetCarryProps(userId, false)
	mainPet := ts.OptSvc.PetSrv.GetPet(userId, user.Mbid)
	mainPet.GetM()
	// 限制要求
	enablePass, msg := ts.PassLimit(task, user, userInfo, mainPet)
	if !enablePass {
		return enablePass, msg
	}

	stateMap := make(map[string]int)
	for _, state := range strings.Split(uTask.State, ",") {
		sItems := strings.Split(state, ":")
		if len(sItems) > 2 && sItems[0] == "killmon" {
			stateMap[sItems[1]] = com.StrTo(sItems[2]).MustInt()
		}
	}
	// 任务要求
	ts.OptSvc.Begin()
	defer ts.OptSvc.Rollback()
	for _, str := range strings.Split(uTask.MModel.OkNeed, ",") {
		items := strings.Split(str, ":")
		switch items[0] {
		case "giveitem":
			pid, num := com.StrTo(items[1]).MustInt(), com.StrTo(items[2]).MustInt()
			if !ts.PropExist(&carryProp, pid, num) || !ts.OptSvc.PropSrv.DecrPropByPid(userId, pid, num) {
				return false, "任务所需道具不足！"
			}
			break
		case "zx":
			break
		case "givejifen":
			num := com.StrTo(items[1]).MustInt()
			if user.Score < num || ts.GetDb().Model(user).Where("score>=?", num).Update(gin.H{"score": gorm.Expr("score-?", num)}).RowsAffected == 0 {
				return false, "任务所需积分不足！"
			}
			break
		case "all_rmb":
			num := com.StrTo(items[1]).MustInt()
			if user.AllRmb < num {
				return false, "任务所需累计充值金额不足！"
			}
			break
		case "giveww":
			num := com.StrTo(items[1]).MustInt()
			if user.Prestige < num || ts.GetDb().Model(user).Where("prestige>=?", num).Update(gin.H{"prestige": gorm.Expr("prestige-?", num)}).RowsAffected == 0 {
				return false, "任务所需威望不足！"
			}
			break
		case "givevip":
			num := com.StrTo(items[1]).MustInt()
			if user.Vip < num || ts.GetDb().Model(user).Where("vip>=?", num).Update(gin.H{"vip": gorm.Expr("vip-?", num)}).RowsAffected == 0 {
				return false, "任务所需VIP积分不足！"
			}
			break
		case "giveml":
			num := com.StrTo(items[1]).MustInt()
			if userInfo.Ml < num || ts.GetDb().Model(user).Where("ml>=?", num).Update(gin.H{"ml": gorm.Expr("ml-?", num)}).RowsAffected == 0 {
				return false, "任务所需魅力不足！"
			}
			break
		case "givemoney":
			num := com.StrTo(items[1]).MustInt()
			if user.Money < num || ts.GetDb().Model(user).Where("money>=?", num).Update(gin.H{"money": gorm.Expr("money-?", num)}).RowsAffected == 0 {
				return false, "任务所需金币不足！"
			}
			break
		case "givedianjuan":
			num := com.StrTo(items[1]).MustInt()
			if user.ActiveScore < num || ts.GetDb().Model(user).Where("active_score>=?", num).Update(gin.H{"active_score": gorm.Expr("active_score-?", num)}).RowsAffected == 0 {
				return false, "任务所需点券不足！"
			}
			break
		case "monself":
			if mainPet.MModel.ID != com.StrTo(items[1]).MustInt() {
				return false, "任务所需主宠要求不符合！"
			}
			break
		case "lv":
			lvItems := strings.Split(items[1], "|")
			if len(lvItems) == 1 {
				if mainPet.Level < com.StrTo(lvItems[0]).MustInt() {
					return false, "任务所需主宠等级不足！"
				}
			} else if len(lvItems) == 2 {
				if mainPet.Level < com.StrTo(lvItems[0]).MustInt() || mainPet.Level > com.StrTo(lvItems[1]).MustInt() {
					return false, "任务所需主宠等级不符合！"
				}
			}
			break
		case "wx":
			if !com.IsSliceContainsStr(strings.Split(items[1], "|"), strconv.Itoa(mainPet.MModel.Wx)) {
				return false, "任务所需主宠五行不符合！"
			}
			break
		case "killmon":
			if num := stateMap[items[1]]; num < com.StrTo(items[2]).MustInt() {
				return false, "任务所需击败怪物目标不足！！"
			}
		}
	}

	addMsg := ""
	logNotes := []string{}

	for _, rs := range strings.Split(uTask.MModel.Result, ",") {
		rItems := strings.Split(rs, ":")
		switch rItems[0] {
		case "item", "itemrand":
			str := strings.ReplaceAll(rs, "itemrand:", "")
			for _, pstr := range strings.Split(str, "|") {
				randItems := strings.Split(pstr, ":")
				if len(randItems) > 2 {
					pid := com.StrTo(randItems[0]).MustInt()
					randNum := com.StrTo(randItems[1]).MustInt()
					num := com.StrTo(randItems[2]).MustInt()
					if rand.Intn(randNum) == 0 {
						p := common.GetMProp(pid)
						if p != nil && ts.OptSvc.PropSrv.AddProp(userId, pid, num, false) {
							dlog := fmt.Sprintf("获得道具 %s * %d", p.Name, num)
							addMsg += dlog
							logNotes = append(logNotes, dlog)
							break
						}
					}
				}
			}
			break
		case "lvprops":
			// 不同等级获得不同道具，暂时不做
			break
		case "props":
			if len(rItems) > 2 {
				pid := com.StrTo(rItems[1]).MustInt()
				num := com.StrTo(rItems[2]).MustInt()
				p := common.GetMProp(pid)
				if p != nil {
					ts.OptSvc.PropSrv.AddProp(userId, pid, num, false)
					logNotes = append(logNotes, fmt.Sprintf("获得道具 %s * %d", p.Name, num))
				}
			}
			break
		case "bprops":
			if len(rItems) > 2 {
				pid := com.StrTo(rItems[1]).MustInt()
				num := com.StrTo(rItems[2]).MustInt()
				p := common.GetMProp(pid)
				if p != nil {
					newprop, ok := ts.OptSvc.PropSrv.AddOrCreateProp(userId, pid, num, false)
					if ok && newprop != nil {
						ts.GetDb().Model(newprop).Update(gin.H{"cantrade": 1})
						logNotes = append(logNotes, fmt.Sprintf("获得道具可交易 %s * %d", p.Name, num))
					} else {
						return false, "任务出错！请联系管理员！"
					}
				}
			}
			break
		case "exp":
			num := com.StrTo(rItems[1]).MustInt()
			if len(rItems) > 3 {
				wwItems := strings.Split(rItems[2], "|")
				if len(wwItems) == 1 {
					if user.Jprestige < com.StrTo(wwItems[0]).MustInt() {
						num = 0
					}
				} else if len(wwItems) == 2 {
					if user.Jprestige < com.StrTo(wwItems[0]).MustInt() || user.Jprestige > com.StrTo(wwItems[1]).MustInt() {
						num = 0
					}
				}
			}
			if num > 0 {
				ts.OptSvc.PetSrv.IncreaseExp2Pet(mainPet, num)
				logNotes = append(logNotes, fmt.Sprintf("获得经验  %d", num))
			}
			break
		case "mon", "money":
			var num int
			num = com.StrTo(rItems[1]).MustInt()
			if len(rItems) > 3 {
				wwItems := strings.Split(rItems[2], "|")
				if len(wwItems) == 1 {
					if user.Jprestige < com.StrTo(wwItems[0]).MustInt() {
						num = 0
					}
				} else if len(wwItems) == 2 {
					if user.Jprestige < com.StrTo(wwItems[0]).MustInt() || user.Jprestige > com.StrTo(wwItems[1]).MustInt() {
						num = 0
					}
				}
			}
			if num > 0 {
				money := num + user.Money
				if money >= utils.MaxJinBi {
					money = utils.MaxJinBi
				}
				ts.GetDb().Model(user).Update(gin.H{"money": money})
				logNotes = append(logNotes, fmt.Sprintf("获得金币 %d", num))
			}
			break
		case "gonggao":
			AnnounceAll(user.Nickname, rItems[1])
			break
		}
	}
	ts.GetDb().Delete(uTask)
	if task.Cid == "0" || task.Cid == "self" && strings.Index(task.LimitLv, "cishu") > -1 {
		ts.GetDb().Create(&models.TaskLog{
			Uid:     userId,
			TaskId:  task.Id,
			Xulie:   task.Xulie,
			Time:    utils.NowUnix(),
			FromNpc: task.Color,
		})
	} else if task.Xulie > 0 && strings.Index(task.Cid, "rwl") > -1 {
		taskLog := &models.TaskLog{}
		ts.GetDb().Model(taskLog).Where("xulie=?", task.Xulie).First(taskLog)
		if taskLog.Id > 0 {
			ts.GetDb().Model(taskLog).Update(gin.H{"taskid": task.Id, "time": utils.NowUnix()})
		} else {
			ts.GetDb().Create(&models.TaskLog{
				Uid:     userId,
				TaskId:  task.Id,
				Xulie:   task.Xulie,
				Time:    utils.NowUnix(),
				FromNpc: task.Color,
			})
		}
	}
	ts.OptSvc.Commit()
	SelfGameLog(userId, fmt.Sprintf("任务id：%d, 任务标题：%s\n 任务结果：%s", task.Id, task.Title, strings.Join(logNotes, "\n")), 161)
	return true, task.Title + " 任务完成！您获得了相应任务奖励！" + addMsg
}

func (ts *TaskServices) ThrowTask(userId, taskId int) (bool, string) {
	task := ts.GetTask(userId, taskId)
	if task == nil {
		return false, "任务不存在！"
	}
	if ts.GetDb().Delete(task).RowsAffected > 0 {

		if strings.Contains(task.GetM().OkNeed, "killmon") {
			go task_helper.TaskHelperInstance.UpdateUserTask(userId)
		}
		return true, "成功放弃任务！"
	}
	return false, "任务不存在！"
}

func FinishFightTask(userId, gpcId int) {
	task_helper.TaskHelperInstance.UpdateTaskState(userId, gpcId)
	//tasks := []*models.UTask{}
	//persistence.GetOrm().Where("uid=?", userId).Find(&tasks)
	//for _, task := range tasks {
	//	task.GetM()
	//	if !strings.Contains(task.MModel.OkNeed, "killmon") {
	//		continue
	//	}
	//	var ids string
	//	var num int
	//	find := false
	//	for _, need := range strings.Split(task.MModel.OkNeed, ",") {
	//		if strings.Contains(need, "killmon") {
	//			items := strings.Split(need, ":")
	//			ids = items[1]
	//			num = com.StrTo(items[2]).MustInt()
	//			if slice.ContainsString(strings.Split(ids, "|"), strconv.Itoa(gpcId)) {
	//				find = true
	//			}
	//		}
	//	}
	//	if find {
	//		stateFind := false
	//		stateItems := strings.Split(task.State, ",")
	//		index := 0
	//		stateNum := 0
	//		for i, state := range stateItems {
	//			if strings.Contains(state, "killmon") {
	//				items := strings.Split(state, ":")
	//				if slice.ContainsString(strings.Split(items[1], "|"), strconv.Itoa(gpcId)) {
	//					stateFind = true
	//					index = i
	//					stateNum = com.StrTo(items[2]).MustInt()
	//					break
	//				}
	//			}
	//		}
	//		if stateNum < num {
	//			stateNum++
	//			if stateFind {
	//				stateItems[index] = fmt.Sprintf("killmon:%s:%d", stateNum)
	//			} else {
	//				stateItems = append(stateItems, fmt.Sprintf("killmon:%s:%d", stateNum))
	//			}
	//			persistence.GetOrm().Model(task).Update(gin.H{"state": strings.Join(stateItems, ",")})
	//		}
	//	}
	//}
}
