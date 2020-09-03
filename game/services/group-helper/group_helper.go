package group_helper

import (
	"pokemon/common/rcache"
	"pokemon/game/models"
	"pokemon/game/services/common"
	"pokemon/game/utils"
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

func (g *UserGroup) GetGpc() *models.Gpc {
	if g.Level < 2 {
		if g.GpcIndex >= len(g.Gpc) {
			return nil
		}
		return common.GetGpc(g.Gpc[g.GpcIndex])
	} else {
		if len(g.Gpc) == 0 || g.GpcIndex != len(g.Gpc) {
			return nil
		}
		return common.GetGpc(g.Gpc[g.GpcIndex-1])
	}

}

func (g *UserGroup) FbNeverStart() bool {
	if g.Multiple == 0 && len(g.Gpc) == 0 && g.Level == 0 && g.Process == 0 {
		return true
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

func (g *UserGroup) Save() {
	// 保存信息，方便后面扩展
}
