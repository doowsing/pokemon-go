package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	common2 "pokemon/common"
	"pokemon/common/rcache"
	"pokemon/common/rpc-client/rpc-group"
	"pokemon/game/ginapp"
	"pokemon/game/services/common"
	"strconv"
	"strings"
)

var FightCtl = &FightController{}

type FightController struct {
}

func CheckOpenMap(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	//id := gapp.Id()
	mapId := com.StrTo(c.DefaultQuery("n", "0")).MustInt()
	if mapId < 1 {
		return
	}
	optType := c.Query("type")
	dMap := common.GetMMap(mapId)
	if dMap == nil {
		gapp.String("1")
		return
	}
	switch optType {
	case "1":
		//user := gapp.OptSvc.UserSrv.GetUserById(id)
		//openMap := strings.Split(user.OpenMap, ",")
		//for _, v := range openMap {
		//	if v == strconv.Itoa(mapId) {
		//		gapp.String("10")
		//		return
		//	}
		//}
		gapp.String("12")
		return
	}
}

func OpenMap(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	ok, msg := gapp.OptSvc.FightSrv.OpenMap(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func OpenMaps(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	gapp.JSONDATAOK("", gin.H{"maps": strings.Split(user.OpenMap, ",")})
}

func GoInToMap(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	info, msg := gapp.OptSvc.FightSrv.IntoMap(gapp.Id(), id)
	gapp.JSONDATAOK(msg, info)

}

func GoInToFbMap(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	info, msg := gapp.OptSvc.FightSrv.IntoFbMap(gapp.Id(), id)
	gapp.JSONDATAOK(msg, info)

}

func StartFight(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	multiple, err := strconv.Atoi(c.Query("multiple"))
	if err != nil || multiple < 1 {
		multiple = 0
	}
	info, msg := gapp.OptSvc.FightSrv.StartFight(gapp.Id(), multiple)
	gapp.JSONDATAOK(msg, info)
}

func Attack(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		//gapp.JSONDATAOK("参数出错！", nil)
		//return
		id = 0
	}
	info, msg := gapp.OptSvc.FightSrv.Attack(gapp.Id(), id)
	gapp.JSONDATAOK(msg, info)
}

func AutoStartFight(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	autoType := c.Query("type")
	if autoType == "" {
		autoType = "jb"
	}
	if autoType != "jb" && autoType != "yb" {
		gapp.JSONDATAOK("参数出错！", gin.H{"result": false})
		return
	}
	ok, msg := gapp.OptSvc.FightSrv.SetAutoStart(gapp.Id(), autoType)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func CancelAutoStartFight(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	rcache.DelAutoFightFlag(gapp.Id())
	gapp.JSONDATAOK("", gin.H{"result": true})
}

func AutoFightSkill(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	ok, msg := gapp.OptSvc.FightSrv.SetAutoSkill(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func TTUseSj(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.FightSrv.TTUseSj(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func TTUserRank(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	data := gapp.OptSvc.FightSrv.GetTTUserRank()
	gapp.JSONDATAOK("", gin.H{"users": data})
}

func MapUsers(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	data := gapp.OptSvc.FightSrv.GetMapUsers(gapp.Id(), id)
	gapp.JSONDATAOK("", gin.H{"users": data})
}

func CatchPet(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("精灵球数量不足！", gin.H{"result": false})
		return
	}
	ok, msg := gapp.OptSvc.FightSrv.Catch(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func StartFb(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("地图不存在！", gin.H{"result": false})
		return
	}
	t := c.Query("type")
	ok, msg := gapp.OptSvc.FightSrv.StartFb(gapp.Id(), id, t == "1")
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func GroupInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"map_id": 0,
		"list":   nil,
		"group":  nil,
	}
	mapId := rcache.GetInMap(gapp.Id())
	data["map_id"] = mapId

	find := false
	groupId := rpc_group.GetGroupID(gapp.Id())
	if groupId != "" {
		group := rpc_group.GetGroup(groupId)
		if group != nil {
			data["map_id"] = group.MapId
			find = true
			members := []gin.H{}
			requestM := []gin.H{}
			for _, u := range group.Member {
				members = append(members, gin.H{"id": u.Id, "nickname": u.Nickname, "ready": u.Ready})
			}
			if group.IsLeader(gapp.Id()) {
				for _, u := range group.RequestMember {
					requestM = append(requestM, gin.H{"id": u.Id, "nickname": u.Nickname})
				}
			}
			data["group"] = gin.H{"uuid": groupId, "leader_name": group.LeaderName, "leader_id": group.Leader, "member": members, "request": requestM}
		}
	}
	if !find {
		data["list"] = rpc_group.GetMapAllGroup(mapId)
	}
	gapp.JSONDATAOK(msg, data)
}

func CreateGroup(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"result": false,
	}

	groupId := rpc_group.GetGroupID(gapp.Id())
	if groupId != "" {
		msg = "您已在一个队伍中！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	mapId := rcache.GetInMap(gapp.Id())
	if mapId == 0 {
		msg = "该地图不可组队！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	rpc_group.NewUserGroup(user.ID, user.Nickname, mapId)
	msg = "建立成功！"
	data["result"] = true
	gapp.JSONDATAOK(msg, data)
}

func DissolveGroup(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"result": false,
	}
	groupId := rpc_group.GetGroupID(gapp.Id())
	if groupId == "" {
		msg = "您不在队伍之中！1"
		gapp.JSONDATAOK(msg, data)
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		msg = "您不在队伍之中！2"
		gapp.JSONDATAOK(msg, data)
		return
	}
	if !group.IsLeader(gapp.Id()) {
		msg = "队员不可解散队伍！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	ok, msg := rpc_group.DropGroup(groupId, gapp.Id())
	data["result"] = ok
	gapp.JSONDATAOK(msg, data)
}

func RequestGroup(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"result": false,
	}

	groupId := c.Query("group_id")
	myGroupId := rpc_group.GetGroupID(gapp.Id())
	if myGroupId != "" {
		msg = "您已在队伍之中！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	group := rpc_group.GetGroup(groupId)
	if group == nil {
		msg = "队伍不存在！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	ok, msg := rpc_group.AddGroupRequest(groupId, user.ID, user.Nickname)
	data["result"] = ok
	gapp.JSONDATAOK(msg, data)
}

func ReceiveGroup(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"result": false,
	}

	id := com.StrTo(c.Query("id")).MustInt()
	if id == 0 {
		msg = "参数出错！"
		gapp.JSONDATAOK(msg, data)
		return
	}

	myGroupId := rpc_group.GetGroupID(gapp.Id())
	if myGroupId == "" {
		msg = "您不在队伍之中！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	group := rpc_group.GetGroup(myGroupId)
	if group == nil {
		msg = "您的队伍不存在！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	if !group.IsLeader(gapp.Id()) {
		msg = "不是队长不可进行操作！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	userInMap := rcache.GetInMap(id)
	ok, msg := rpc_group.ReceiveUser(myGroupId, id, userInMap)
	data["result"] = ok
	gapp.JSONDATAOK(msg, data)
}

func RefuseGroup(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"result": false,
	}

	id := com.StrTo(c.Query("id")).MustInt()
	if id == 0 {
		msg = "参数出错！"
		gapp.JSONDATAOK(msg, data)
		return
	}

	myGroupId := rpc_group.GetGroupID(gapp.Id())
	if myGroupId == "" {
		msg = "您不在队伍之中！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	group := rpc_group.GetGroup(myGroupId)
	if group == nil {
		msg = "您的队伍不存在！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	if !group.IsLeader(gapp.Id()) {
		msg = "不是队长不可进行操作！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	ok, msg := rpc_group.RefuseUser(myGroupId, id)
	data["result"] = true
	if ok {
		msg = "拒绝成功！"
	} else {
		msg = "对方并没有申请您的队伍！"
	}
	gapp.JSONDATAOK(msg, data)
}

func KickOutGroup(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"result": false,
	}

	id := com.StrTo(c.Query("id")).MustInt()
	if id == 0 {
		msg = "参数出错！"
		gapp.JSONDATAOK(msg, data)
		return
	}

	myGroupId := rpc_group.GetGroupID(gapp.Id())
	if myGroupId == "" {
		msg = "您不在队伍之中！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	group := rpc_group.GetGroup(myGroupId)
	if group == nil {
		msg = "您的队伍不存在！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	if !group.IsLeader(gapp.Id()) {
		msg = "不是队长不可进行操作！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	ok, msg := rpc_group.KickOut(myGroupId, id)
	data["result"] = ok
	gapp.JSONDATAOK(msg, data)
}

func ExitGroup(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"result": false,
	}

	myGroupId := rpc_group.GetGroupID(gapp.Id())
	if myGroupId == "" {
		msg = "您不在队伍之中！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	group := rpc_group.GetGroup(myGroupId)
	if group == nil {
		msg = "您的队伍不存在！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	if group.IsLeader(gapp.Id()) {
		msg = "队长不可退出！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	ok, msg := rpc_group.ExitGroup(myGroupId, gapp.Id())
	data["result"] = ok
	gapp.JSONDATAOK(msg, data)
}

func InviteGroup(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msg := ""
	data := gin.H{
		"result": false,
	}

	id := com.StrTo(c.Query("id")).MustInt()
	if id == 0 {
		msg = "参数出错！"
		gapp.JSONDATAOK(msg, data)
		return
	}

	otherGroupId := rpc_group.GetGroupID(id)
	if otherGroupId != "" {
		msg = "玩家已在其他队伍之中！"
		gapp.JSONDATAOK(msg, data)
		return
	}

	myGroupId := rpc_group.GetGroupID(gapp.Id())
	if myGroupId == "" {
		msg = "您不在队伍之中！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	group := rpc_group.GetGroup(myGroupId)
	if group == nil {
		msg = "您的队伍不存在！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	if !group.IsLeader(gapp.Id()) {
		msg = "不是队长不可进行操作！"
		gapp.JSONDATAOK(msg, data)
		return
	}
	common2.GroupInvite(group.UUId, id, group.Leader)
	data["result"] = true
	gapp.JSONDATAOK(msg, data)
}

func GroupStartFight(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	useSj := false
	useSjStr := c.DefaultQuery("usesj", "0")
	if useSjStr == "1" {
		useSj = true
	}
	multiple := com.StrTo(c.Query("multiple")).MustInt()

	data, msg := gapp.OptSvc.FightSrv.StartGroupFight(gapp.Id(), multiple, useSj)
	gapp.JSONDATAOK(msg, data)
}

func GroupAttack(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	id := com.StrTo(c.Query("id")).MustInt()
	data, msg := gapp.OptSvc.FightSrv.GroupAttack(gapp.Id(), id)
	gapp.JSONDATAOK(msg, data)
}

func GroupEnterCard(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data, msg := gapp.OptSvc.FightSrv.GetGroupCardData(gapp.Id())
	gapp.JSONDATAOK(msg, data)
}

func GroupDoCard(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	position := com.StrTo(c.Query("position")).MustInt()
	data, msg := gapp.OptSvc.FightSrv.DoGroupCard(gapp.Id(), position)
	gapp.JSONDATAOK(msg, data)
}

func GroupAllCard(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data, msg := gapp.OptSvc.FightSrv.GetAllGroupCardData(gapp.Id())
	gapp.JSONDATAOK(msg, data)
}

func GroupEnterBossCard(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	data, msg := gapp.OptSvc.FightSrv.GetGroupBossCard(gapp.Id())
	gapp.JSONDATAOK(msg, data)
}

func GroupSetStatus(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	status, err := strconv.Atoi(c.DefaultQuery("status", "-1"))
	if err != nil {
		gapp.JSONDATAOK("参数出错！", gin.H{"result": false})
		return
	}
	data, msg := gapp.OptSvc.FightSrv.SetGroupStatus(gapp.Id(), status)
	gapp.JSONDATAOK(msg, data)
}

func GroupInmap(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	group := rpc_group.GetGroup(c.Query("uuid"))
	if group == nil {
		gapp.JSONDATAOK("队伍不存在！", gin.H{"result": false, "map_id": 0})
		return
	}
	gapp.JSONDATAOK("", gin.H{"result": true, "map_id": group.MapId})
}

func GroupDoBossCard(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	position := com.StrTo(c.Query("position")).MustInt()
	data, msg := gapp.OptSvc.FightSrv.DoGroupBossCard(gapp.Id(), position)
	gapp.JSONDATAOK(msg, data)
}
