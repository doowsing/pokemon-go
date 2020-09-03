package rdc_server

import (
	"pokemon/common/model"
	"pokemon/common/rcache"
)

type GroupHandle struct {
}

func (gh *GroupHandle) NewUserGroup(req *model.CreateGroupRequest, res *model.DataResponse) error {
	userGroup := model.NewUserGroup(req.UserId, req.Nickname, req.MapId)
	res.Data = userGroup
	return nil
}

func (gh *GroupHandle) GetGroup(req *model.CreateGroupRequest, res *model.DataResponse) (e error) {

	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true
	res.Data = group
	return
}

func (gh *GroupHandle) ResetFightInfo(req *model.IdRequest, res *model.ResultResponse) error {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return nil
	}
	res.Exist = true

	group.ResetFightInfo()
	return nil
}

func (gh *GroupHandle) ResetCardAwards(req *model.IdRequest, res *model.ResultResponse) error {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return nil
	}
	res.Exist = true

	group.ResetCardAwards()
	return nil
}

func (gh *GroupHandle) AddAwardProp(req *model.RateSetRequest, res *model.ResultResponse) error {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return nil
	}
	res.Exist = true

	group.AddAwardProp(req.RateSets)
	return nil
}

func (gh *GroupHandle) SetNextFightUserId(req *model.IdRequest, res *model.ResultResponse) error {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return nil
	}
	res.Exist = true

	ok := group.SetNextFightUserId()
	res.Ok = ok
	return nil
}

func (gh *GroupHandle) SetFbTime(req *model.IdRequest, res *model.ResultResponse) error {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return nil
	}
	res.Exist = true

	ok := group.SetFbTime()
	res.Ok = ok
	return nil
}

func (gh *GroupHandle) GetMapAllGroup(mapId int, data *[]map[string]interface{}) error {
	_data := []map[string]interface{}{}
	for _, g := range model.GetMapAllGroup(mapId) {
		_data = append(_data, map[string]interface{}{"uuid": g.UUId, "leader_name": g.LeaderName, "member_num": len(g.Member)})
	}
	*data = _data
	return nil
}

func (gh *GroupHandle) DropGroup(req *model.IdRequest, res *model.ResultResponse) error {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return nil
	}
	res.Exist = true

	res.Ok, res.Msg = group.Drop(req.UserId)
	return nil
}

func (gh *GroupHandle) AddGroupRequest(req *model.IdRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	res.Ok, res.Msg = group.AddRequest(req.UserId, req.Nickname)
	return
}

func (gh *GroupHandle) ReceiveUser(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	res.Ok, res.Msg = group.ReceiveUser(req.UserId, req.ParamInt)
	return
}

func (gh *GroupHandle) RefuseUser(req *model.IdRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.RefuseUser(req.UserId)
	return
}

func (gh *GroupHandle) KickOut(req *model.IdRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	res.Ok, res.Msg = group.DelUser(req.UserId, false)
	return
}

func (gh *GroupHandle) ExitGroup(req *model.IdRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	res.Ok, res.Msg = group.DelUser(req.UserId, true)
	return
}

func (gh *GroupHandle) SetFightUser(req *model.IdRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true
	group.FightUserId = req.UserId
	return
}

func (gh *GroupHandle) SetMultiple(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.Multiple = req.ParamInt
	return
}

func (gh *GroupHandle) AddMoney(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.GetMoney += req.ParamInt
	return
}

func (gh *GroupHandle) AddExp(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.GetExp += req.ParamInt
	return
}

func (gh *GroupHandle) SetGpcs(req *model.ParamSliceIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.Gpc = req.ParamInts
	return
}

func (gh *GroupHandle) SetGpcIndex(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.GpcIndex = req.ParamInt
	return
}

func (gh *GroupHandle) SetGpcDeHp(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.GpcDeHp = req.ParamInt
	return
}

func (gh *GroupHandle) SetFbEnd(req *model.ParamBoolRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.End = req.ParamBool
	return
}

func (gh *GroupHandle) SetFbLevel(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.Level = req.ParamInt
	return
}

func (gh *GroupHandle) SetFbProcess(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.Process = req.ParamInt
	return
}

func (gh *GroupHandle) SetFbStartCardTime(req *model.ParamIntRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.StartCardTime = req.ParamInt
	return
}

func (gh *GroupHandle) SetFbCardAwards(req *model.CardAwardsRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.CardAwards = req.Param
	return
}

func (gh *GroupHandle) SetFbCardAward(req *model.CardAwardRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.CardAwards[req.Position] = req.Param
	return
}

func (gh *GroupHandle) AddFbBossCardEffects(req *model.BossCardEffectRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.BossCardEffect = append(group.BossCardEffect, req.Param)
	return
}

func (gh *GroupHandle) SetMemberStatus(req *model.MemberStatusRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	group.Member = req.Members
	return
}

func (gh *GroupHandle) SetUserStatus(req *model.ParamBoolRequest, res *model.ResultResponse) (e error) {
	group := model.GetGroup(req.UUid)
	if group == nil {
		return
	}
	res.Exist = true

	if req.ParamBool {
		group.SetReady(req.UserId)
	} else {
		group.SetUnReady(req.UserId)
	}
	return
}

func (gh *GroupHandle) GetGroupID(userId int, res *model.ResultResponse) (e error) {
	res.Msg = rcache.GetGroupID(userId)
	return
}

func (gh *GroupHandle) CheckConnect(a int, res *model.ResultResponse) (e error) {
	return
}
