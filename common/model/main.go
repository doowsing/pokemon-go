package model

type CreateGroupRequest struct {
	UUid     string
	UserId   int
	Nickname string
	MapId    int
}

type MemberStatusRequest struct {
	UUid    string
	Members []*MemberStatus
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

type CardAwardsRequest struct {
	UUid  string
	Param []*CardInfo
}

type CardAwardRequest struct {
	UUid     string
	Position int
	Param    *CardInfo
}

type BossCardEffectRequest struct {
	UUid  string
	Param *CardEffectInfo
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
	RateSets []*RateSet
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

type RateSet struct {
	Pid  int
	Rate int
}
