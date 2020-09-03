package model

import (
	"github.com/devfeel/dotweb/framework/crypto/uuid"
	"log"
	"math/rand"
	common2 "pokemon/common"
	"pokemon/common/rcache"
	"pokemon/common/utils"
	"sync"
)

const MaxMember = 5

type UserGroup struct {
	UUId          string
	LeaderName    string
	Leader        int
	MapId         int
	Member        []*MemberStatus
	RequestMember []*MemberStatus

	// 本次战斗信息
	// 需要存放玩家dehp,demp,怪物hp

	FightUserId int
	Multiple    int
	GetMoney    int
	GetExp      int
	GetProps    map[int][]int //玩家获得的道具
	Gpc         []int
	GpcIndex    int
	GpcDeHp     int
	FightStatus *struct {
		ForceFinish bool
		FinishMsg   string

		MeDie    bool
		EmeryDie bool
	}

	// 遗忘副本属性
	// 玩家状态：打过，刚进还没开始，开始了而且打完了
	End            bool // 是否已经结束了
	Level          int  // 层级，总共3层，第一、二层为刷怪，第三层为boss翻牌
	Process        int  // 怪物进度
	StartCardTime  int  // 进入翻牌时间，总共可以翻30秒，期间可以出去再进
	CardAwards     []*CardInfo
	BossCardEffect []*CardEffectInfo
}

type MemberStatus struct {
	Id       int
	Nickname string
	Ready    bool
}

// 遗忘一二层通关后添加，固定数量10个
type CardInfo struct {
	TarotCardId int
	Content     string
	UserId      int // 0则并没有被人翻过
	Nickname    string
}

// 遗忘第三层通关后添加，固定数量10个
type CardEffectInfo struct {
	Position int
	Img      string
}

func NewUserGroup(leader int, leaderName string, mapId int) *UserGroup {
	g := &UserGroup{
		UUId:       uuid.NewV4().String(),
		Leader:     leader,
		LeaderName: leaderName,
		MapId:      mapId,
		Member:     []*MemberStatus{{leader, leaderName, true}},
		FightStatus: &struct {
			ForceFinish bool
			FinishMsg   string

			MeDie    bool
			EmeryDie bool
		}{},
	}
	rcache.SetGroupID(leader, g.UUId)
	log.Printf("set user %d uuid %s\n", leader, g.UUId)
	common2.GroupUpdateId(g.UUId, leader)
	AddGroup(g)
	return g
}

func (g *UserGroup) ResetFightInfo() {
	g.FightStatus.ForceFinish = false
	g.FightStatus.FinishMsg = ""
	g.FightStatus.MeDie = false
	g.FightStatus.EmeryDie = false

	g.GetMoney = 0
	g.GetExp = 0
	g.GetProps = make(map[int][]int)
	g.GpcIndex = 0
	g.GpcDeHp = 0
}

func (g *UserGroup) ResetCardAwards() {
	g.CardAwards = []*CardInfo{}
}

func (g *UserGroup) IsLeader(userId int) bool {
	if userId != g.Leader {
		return false
	}
	return true
}

func (g *UserGroup) GetGpc() int {
	if g.Level < 2 {
		if g.GpcIndex >= len(g.Gpc) {
			return 0
		}
		return g.Gpc[g.GpcIndex]
	} else {
		if len(g.Gpc) == 0 || g.GpcIndex != len(g.Gpc) {
			return 0
		}
		return g.Gpc[g.GpcIndex-1]
	}

}

func (g *UserGroup) AddAwardProp(rateSet []*RateSet) {
	for _, m := range g.Member {
		rateNum := 2
		if m.Id == g.FightUserId {
			rateNum = 1
		}
		for _, set := range rateSet {
			if rand.Intn(set.Rate*rateNum) == 0 {
				g.GetProps[m.Id] = append(g.GetProps[m.Id], set.Pid)
			}
		}
	}
}

func (g *UserGroup) FbNeverStart() bool {
	if g.Multiple == 0 && len(g.Gpc) == 0 && g.Level == 0 && g.Process == 0 {
		return true
	}
	return false
}

func (g *UserGroup) SetNextFightUserId() bool {
	find := false
	for _, m := range g.Member {
		if find == true {
			g.FightUserId = m.Id
			return true
		}
		if m.Id == g.FightUserId {
			find = true
		}
	}

	return false
}

func (g *UserGroup) FbNeedSj() bool {
	today := utils.ToDayStartUnix()
	for _, m := range g.Member {
		if t := rcache.GetYiWangTime(m.Id); today < t {
			return true
		}
	}
	return false
}

func (g *UserGroup) SetFbTime() bool {
	now := utils.NowUnix()
	for _, m := range g.Member {
		rcache.SetYiWangTime(m.Id, now)
	}
	return false
}

func (g *UserGroup) GetMember(userId int) *MemberStatus {
	for _, m := range g.Member {
		if m.Id == userId {
			return m
		}
	}
	return nil
}

func (g *UserGroup) InGroup(userId int) bool {
	if g.GetMember(userId) != nil {
		return true
	}
	return false
}

func (g *UserGroup) AddRequest(userId int, nickname string) (bool, string) {
	if len(g.Member) >= MaxMember {
		return false, "成员已满！"
	}
	for _, m := range g.RequestMember {
		if m.Id == userId {
			return false, "已经申请过该队伍！"
		}
	}
	g.RequestMember = append(g.RequestMember, &MemberStatus{
		Id:       userId,
		Ready:    false,
		Nickname: nickname,
	})
	common2.GroupNotice(&common2.GroupCommandMsg{
		Command: "update-uuid",
		UUid:    g.UUId,
		Receive: g.Leader,
		Data:    nil,
	})

	common2.GroupUpdate(g.UUId, g.Leader)
	return true, "申请加入成功！"
}

func (g *UserGroup) ReceiveUser(userId, inMap int) (bool, string) {
	if len(g.Member) >= MaxMember {
		return false, "成员已满！"
	}
	requestMembers := []*MemberStatus{}
	var member *MemberStatus
	for _, m := range g.RequestMember {
		if m.Id == userId {
			member = m
		} else {
			requestMembers = append(requestMembers, m)
		}
	}
	if member == nil {
		return false, "该成员尚未申请入队！"
	}
	if g.InGroup(userId) {
		return false, "该成员已在队伍内！"
	}
	g.RequestMember = requestMembers
	if inMap != g.MapId {
		log.Printf("userId :%d, user map:%d, group map:%d\n", userId, inMap, g.MapId)
		return true, "该玩家已不在本地图！"
	}
	if rcache.GetGroupID(userId) != "" {
		return true, "该玩家已在其他队伍内！"
	}
	g.Member = append(g.Member, member)
	rcache.SetGroupID(userId, g.UUId)

	common2.GroupUpdateId(g.UUId, member.Id)

	common2.GroupUpdate(g.UUId, 0)
	return true, "进入队伍成功！"
}

func (g *UserGroup) RefuseUser(userId int) {
	find := false
	requestMembers := []*MemberStatus{}
	for _, m := range g.RequestMember {
		if m.Id == userId {
			find = true
		} else {
			requestMembers = append(requestMembers, m)
		}
	}
	if find {
		g.RequestMember = requestMembers
	}
}

func (g *UserGroup) DelUser(userId int, Exit bool) (bool, string) {
	if g.IsLeader(userId) {
		return false, "队长无法退出队伍！"
	}
	newMembers := []*MemberStatus{}
	find := false
	for _, m := range g.Member {
		if m.Id == userId {
			find = true
		} else {
			newMembers = append(newMembers, m)
		}
	}
	if find {
		g.Member = newMembers
		rcache.DelGroupID(userId)
		common2.GroupUpdate(g.UUId, 0)
		common2.GroupUpdateId("", userId)
		if Exit {
			return true, "退出成功！"
		} else {
			return true, "踢出成功！"
		}
	}
	if Exit {
		return false, "您已不在队伍中！"
	} else {
		return false, "成员不在队伍之中！"
	}
}

func (g *UserGroup) IsReady(userId int) (bool, string) {
	if m := g.GetMember(userId); m != nil {
		return m.Ready, ""
	}
	return false, "该用户不在队伍中！"
}

func (g *UserGroup) AllReady() bool {
	for _, m := range g.Member {
		if !m.Ready {
			return false
		}
	}
	return true
}

func (g *UserGroup) SetReady(userId int) (bool, string) {
	if m := g.GetMember(userId); m != nil {
		m.Ready = true
		common2.GroupUpdate(g.UUId, 0)
		return true, ""
	}
	return false, "该用户不在队伍中！"
}

func (g *UserGroup) SetUnReady(userId int) (bool, string) {
	if m := g.GetMember(userId); m != nil {
		m.Ready = false
		common2.GroupUpdate(g.UUId, 0)
		return true, ""
	}
	return false, "该用户不在队伍中！"
}

func (g *UserGroup) Drop(userId int) (bool, string) {
	if !g.IsLeader(userId) {
		return false, "不是队长无法解散队伍！"
	}
	DelGroup(g.UUId)
	common2.GroupDissolve(g.UUId)
	for _, m := range g.Member {
		rcache.DelGroupID(m.Id)
	}
	g.Member = nil
	return true, "队伍已被解散！"
}

func (g *UserGroup) Save() {
	// 保存信息，方便后面扩展
}

var GroupHub sync.Map

func DelGroup(uuid string) {
	GroupHub.Delete(uuid)
}

func AddGroup(g *UserGroup) {
	GroupHub.Store(g.UUId, g)
}

func GetGroup(uuid string) *UserGroup {
	if data, ok := GroupHub.Load(uuid); ok {
		return data.(*UserGroup)
	}
	return nil
}

func GetMapAllGroup(mapId int) []*UserGroup {
	var groups []*UserGroup
	GroupHub.Range(func(key, value interface{}) bool {
		if g := value.(*UserGroup); g.MapId == mapId {
			groups = append(groups, g)
		}
		return true
	})
	return groups
}
