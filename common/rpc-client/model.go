package rpc_client

import (
	"pokemon/game/models"
	group_helper "pokemon/game/services/group-helper"
)

type CreateGroupRequest struct {
	UUid     string
	UserId   int
	Nickname string
	MapId    int
}

type ParamIntRequest struct {
	UUid     string
	UserId   int
	ParamInt int
}
type ParamBoolRequest struct {
	UUid      string
	UserId    int
	ParamBool bool
}

type MemberStatusRequest struct {
	UUid    string
	Members []*group_helper.MemberStatus
}

type CardAwardsRequest struct {
	UUid  string
	Param []*group_helper.CardInfo
}

type CardAwardRequest struct {
	UUid     string
	Position int
	Param    *group_helper.CardInfo
}

type BossCardEffectRequest struct {
	UUid  string
	Param *group_helper.CardEffectInfo
}

type ParamSliceIntRequest struct {
	UUid      string
	UserId    int
	ParamInts []int
}

type IdRequest struct {
	UUid     string
	UserId   int
	Nickname string
}

type RateSetRequest struct {
	UUid     string
	RateSets []*models.RatePid
}

type ResultResponse struct {
	Exist bool
	Ok    bool
	Msg   string
}

type DataResponse struct {
	Exist bool
	Data  interface{}
}
