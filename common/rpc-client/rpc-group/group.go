package rpc_group

import (
	"errors"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"pokemon/common/rpc-client"
	"pokemon/game/models"
	"pokemon/game/services/group-helper"
	"time"
)

const (
	Group_ADDRRESS = "127.0.0.1:1234"
)

var client *rpc.Client
var _needConnect = true

func InitGroupRpcClient() {
	go func() {
		for {
			if _needConnect {
				connectGroupRpc()
			}
			time.Sleep(time.Second)
			CheckConnect()
		}
	}()
}

func connectGroupRpc() bool {
	var err error
	client, err = jsonrpc.Dial("tcp", Group_ADDRRESS)
	if err != nil {
		log.Print("rpc dialing error:", err)
		return false
	}
	_needConnect = false
	return true
}

func dealRpcError(err error) {
	log.Println("rpc error:", err)
	_needConnect = true
}

func GetGroup(uuid string) *group_helper.UserGroup {
	group := &group_helper.UserGroup{}
	reply := rpc_client.DataResponse{
		Data: group,
	}
	err := client.Call("GroupHandle.GetGroup", &rpc_client.IdRequest{UUid: uuid}, &reply)
	if err != nil {
		dealRpcError(err)
	}
	if !reply.Exist {
		return nil
	}
	return group
}

func GetMapAllGroup(mapId int) []map[string]interface{} {
	data := []map[string]interface{}{}
	err := client.Call("GroupHandle.GetMapAllGroup", mapId, &data)
	if err != nil {
		dealRpcError(err)
	}
	return data
}

func NewUserGroup(leader int, leaderName string, mapId int) *group_helper.UserGroup {
	group := &group_helper.UserGroup{}
	reply := rpc_client.DataResponse{
		Data: group,
	}
	err := client.Call("GroupHandle.NewUserGroup", &rpc_client.CreateGroupRequest{UserId: leader, Nickname: leaderName, MapId: mapId}, &reply)
	if err != nil {
		dealRpcError(err)
		return nil
	}
	return group
}

func DropGroup(uuid string, userId int) (bool, string) {
	req := &rpc_client.IdRequest{UUid: uuid, UserId: userId}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.DropGroup", req, &res)
	if err != nil {
		dealRpcError(err)
		return false, "服务器错误！"
	}
	if !res.Exist {
		return false, "队伍不存在！"
	}

	return res.Ok, res.Msg
}

func AddGroupRequest(uuid string, userId int, nickname string) (bool, string) {
	req := &rpc_client.IdRequest{UUid: uuid, UserId: userId, Nickname: nickname}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.AddGroupRequest", req, &res)
	if err != nil {
		dealRpcError(err)
		return false, "服务器错误！"
	}
	if !res.Exist {
		return false, "队伍不存在！"
	}

	return res.Ok, res.Msg
}

func ReceiveUser(uuid string, userId, inmap int) (bool, string) {
	req := &rpc_client.ParamIntRequest{UUid: uuid, UserId: userId, ParamInt: inmap}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.ReceiveUser", req, &res)
	if err != nil {
		dealRpcError(err)
		return false, "服务器错误！"
	}
	if !res.Exist {
		return false, "队伍不存在！"
	}

	return res.Ok, res.Msg
}

func RefuseUser(uuid string, userId int) (bool, string) {
	req := &rpc_client.IdRequest{UUid: uuid, UserId: userId}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.RefuseUser", req, &res)
	if err != nil {
		dealRpcError(err)
		return false, "服务器错误！"
	}
	if !res.Exist {
		return false, "队伍不存在！"
	}

	return true, ""
}

func KickOut(uuid string, memberId int) (bool, string) {
	req := &rpc_client.IdRequest{UUid: uuid, UserId: memberId}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.KickOut", req, &res)
	if err != nil {
		dealRpcError(err)
		return false, "服务器错误！"
	}
	if !res.Exist {
		return false, "队伍不存在！"
	}

	return res.Ok, res.Msg
}

func ExitGroup(uuid string, memberId int) (bool, string) {
	req := &rpc_client.IdRequest{UUid: uuid, UserId: memberId}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.ExitGroup", req, &res)
	if err != nil {
		dealRpcError(err)
		return false, "服务器错误！"
	}
	if !res.Exist {
		return false, "队伍不存在！"
	}

	return res.Ok, res.Msg
}

func AddAwardProp(uuid string, rateset []*models.RatePid) {
	req := &rpc_client.RateSetRequest{UUid: uuid, RateSets: rateset}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.AddAwardProp", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func ResetFightInfo(group *group_helper.UserGroup) {
	group.ResetFightInfo()
	req := &rpc_client.IdRequest{UUid: group.UUId}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.ResetFightInfo", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func ResetCardAwards(group *group_helper.UserGroup) {
	group.ResetCardAwards()
	req := &rpc_client.IdRequest{UUid: group.UUId}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.ResetCardAwards", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetFightUser(group *group_helper.UserGroup, userId int) {
	group.FightUserId = userId
	req := &rpc_client.IdRequest{UUid: group.UUId, UserId: userId}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetFightUser", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetMultiple(group *group_helper.UserGroup, multiple int) {
	group.Multiple = multiple
	req := &rpc_client.ParamIntRequest{UUid: group.UUId, ParamInt: multiple}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetMultiple", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func AddMoney(group *group_helper.UserGroup, money int) {
	group.GetMoney += money
	req := &rpc_client.ParamIntRequest{UUid: group.UUId, ParamInt: money}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.AddMoney", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func AddExp(group *group_helper.UserGroup, exp int) {
	group.GetExp += exp
	req := &rpc_client.ParamIntRequest{UUid: group.UUId, ParamInt: exp}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.AddExp", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetGpcs(group *group_helper.UserGroup, gpcs []int) {
	group.Gpc = gpcs
	req := &rpc_client.ParamSliceIntRequest{UUid: group.UUId, ParamInts: gpcs}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetGpcs", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetGpcIndex(group *group_helper.UserGroup, index int) {
	group.GpcIndex = index
	req := &rpc_client.ParamIntRequest{UUid: group.UUId, ParamInt: index}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetGpcIndex", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetGpcDeHp(group *group_helper.UserGroup, deHp int) {
	group.GpcDeHp = deHp
	req := &rpc_client.ParamIntRequest{UUid: group.UUId, ParamInt: deHp}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetGpcDeHp", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetFbEnd(group *group_helper.UserGroup, end bool) {
	group.End = end
	req := &rpc_client.ParamBoolRequest{UUid: group.UUId, ParamBool: end}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetFbEnd", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetFbLevel(group *group_helper.UserGroup, level int) {
	group.Level = level
	req := &rpc_client.ParamIntRequest{UUid: group.UUId, ParamInt: level}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetFbLevel", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetFbProcess(group *group_helper.UserGroup, process int) {
	group.Process = process
	req := &rpc_client.ParamIntRequest{UUid: group.UUId, ParamInt: process}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetFbProcess", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetFbStartCardTime(group *group_helper.UserGroup, t int) {
	group.StartCardTime = t
	req := &rpc_client.ParamIntRequest{UUid: group.UUId, ParamInt: t}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetFbStartCardTime", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetFbCardAwards(group *group_helper.UserGroup, CardAwards []*group_helper.CardInfo) {
	group.CardAwards = CardAwards
	req := &rpc_client.CardAwardsRequest{UUid: group.UUId, Param: CardAwards}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetFbCardAwards", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetFbCardAward(group *group_helper.UserGroup, position int, CardAward *group_helper.CardInfo) {
	group.CardAwards[position] = CardAward
	req := &rpc_client.CardAwardRequest{UUid: group.UUId, Position: position, Param: CardAward}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetFbCardAward", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func AddFbBossCardEffects(group *group_helper.UserGroup, BossCardEffect *group_helper.CardEffectInfo) {
	group.BossCardEffect = append(group.BossCardEffect, BossCardEffect)
	req := &rpc_client.BossCardEffectRequest{UUid: group.UUId, Param: BossCardEffect}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.AddFbBossCardEffects", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetMemberStatus(group *group_helper.UserGroup, members []*group_helper.MemberStatus) {
	group.Member = members
	req := &rpc_client.MemberStatusRequest{UUid: group.UUId, Members: members}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetMemberStatus", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetUserStatus(group *group_helper.UserGroup, userId int, ready bool) {
	req := &rpc_client.ParamBoolRequest{UUid: group.UUId, UserId: userId, ParamBool: ready}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetUserStatus", req, &res)
	if err != nil {
		dealRpcError(err)
		return
	}
	if !res.Exist {
		return
	}
}

func SetNextFightUserId(group *group_helper.UserGroup) bool {
	req := &rpc_client.IdRequest{UUid: group.UUId}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.SetNextFightUserId", req, &res)
	if err != nil {
		dealRpcError(err)
		return false
	}
	if !res.Exist {
		return false
	}
	return res.Ok
}

func GetGroupID(userId int) string {
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.GetGroupID", userId, &res)
	if err != nil {
		dealRpcError(err)
		return ""
	}
	return res.Msg
}

func CheckConnect() {
	if client == nil {
		dealRpcError(errors.New("nil client"))
		return
	}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.CheckConnect", 1, &res)
	if err != nil {
		dealRpcError(err)
	}
}

func CheckConnectError() error {
	if client == nil {
		return errors.New("nil client")
	}
	res := rpc_client.ResultResponse{}
	err := client.Call("GroupHandle.CheckConnect", 1, &res)
	return err
}

func GetClient() *rpc.Client {
	return client
}
