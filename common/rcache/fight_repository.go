package rcache

import (
	"encoding/json"
	"fmt"
	"pokemon/game/common"
	"pokemon/game/models"
	"pokemon/game/utils"
	"strconv"
	"sync"
)

const (
	ZBATTRIBUTE        = "ZBATTRIBUTE"
	FIGHTTIME          = "FIGHTTIME"
	FIGHTTIMECOOLTIME  = 10
	ATTACKTIME         = "ATTACKTIME"
	ATTACKTIMECOOLTIME = 2
	PETZBATTR          = "PETZBATTR_"
	INMAP              = "INMAP"
	INMAPTIME          = "INMAPTIME"
	PETSTATUS          = "PetStatus_"
	USERFIGHTSTATUS    = "UserFightStatus"
	USERAUTOFIGHTFLAG  = "UserAutoFightFlag"
	USERAUTOFIGHTSKILL = "UserAutoFightSkill"
	USERTTFLAG         = "UserTTFlag"
	USERTTRECORD       = "UserTTRecord"
	YIWANGFLAG         = "yiwangFlag"
	SSBattleEndTime    = "SSBattleEndTime"
	SSBattleStartTime  = "SSBattleStartTime"

	GroupKey = "GroupKey"

	AutoFightNone = 0
	AutoFightJb   = 1
	AutoFightYb   = 2

	UserTTNone = 0
	UserTTFh   = 1
	UserTT31   = 2
)

var mapUsers = make([]sync.Map, 500)

func DelMapUser(userId, mapId int) {
	mapUsers[mapId].Delete(userId)
}

func setMapUser(userId, mapId int) {
	mapUsers[mapId].Store(userId, 1)
}

func GetMapUserList(mapId int) []int {
	ids := []int{}
	mapUsers[mapId].Range(func(key, value interface{}) bool {
		ids = append(ids, key.(int))
		return true
	})
	return ids
}

func SetMMaps(maps *[]models.Map) (bool, error) {
	RdbOperator.Delete(common.MMapKey)

	mapMap := make(map[string]interface{})
	for _, v := range *maps {
		//fr.UnSerializerMap(&v)
		mapMap[strconv.Itoa(v.ID)] = v
	}
	return RdbOperator.Hmset(common.MMapKey, mapMap).Bool()
}

func SetMGpcs(mgpcs *[]models.Gpc) (bool, error) {
	RdbOperator.Delete(common.MGpcKey)
	gpcmap := make(map[string]interface{})
	for _, v := range *mgpcs {
		gpcmap[strconv.Itoa(v.ID)] = v
	}
	return RdbOperator.Hmset(common.MGpcKey, gpcmap).Bool()
}

func SetMGpcGroup(mgpcGroup *[]models.GpcGroup) (bool, error) {
	RdbOperator.Delete(common.MGpcGroupKey)
	gpcGroupMap := make(map[string]interface{})
	for _, v := range *mgpcGroup {
		//fr.UnSerializerGpcGroup(&v)
		gpcGroupMap[strconv.Itoa(v.ID)] = v
	}
	return RdbOperator.Hmset(common.MGpcGroupKey, gpcGroupMap).Bool()
}

func GetMap(id int) (*models.Map, error) {
	mmap := models.Map{}
	err := RdbOperator.Hget(common.MMapKey, strconv.Itoa(id)).Struct(&mmap)
	if err == nil {
		return &mmap, err
	}
	return nil, err
}

func GetGpc(id int) (*models.Gpc, error) {
	mgpc := models.Gpc{}
	err := RdbOperator.Hget(common.MGpcKey, strconv.Itoa(id)).Struct(&mgpc)
	if err == nil {
		return &mgpc, err
	}
	return nil, err
}

func GetGpcGroup(id int) (*models.GpcGroup, error) {
	mgpcGroup := models.GpcGroup{}
	err := RdbOperator.Hget(common.MGpcGroupKey, id).Struct(&mgpcGroup)
	if err == nil {
		return &mgpcGroup, err
	}
	return nil, err
}

func GetPetAttribute(PetId int) *models.PetZbAttr {
	data, ok := GetGCache().Get(PETZBATTR + strconv.Itoa(PetId))
	if !ok {
		return nil
	}
	bs, ok := data.([]byte)
	if !ok {
		return nil
	}
	attr := &models.PetZbAttr{}
	err := json.Unmarshal(bs, attr)
	if err == nil {
		return attr
	}
	return nil
}

func SetPetAttribute(PetId int, zbAttribute *models.PetZbAttr) {
	bs, _ := json.Marshal(zbAttribute)
	GetGCache().Set(PETZBATTR+strconv.Itoa(PetId), bs)
}

func ClearPetAttribute(PetId int) {
	GetGCache().Del(PETZBATTR + strconv.Itoa(PetId))
}

func GetFightTime(UserId int) int {
	data, ok := GetGCache().Get(FIGHTTIME + strconv.Itoa(UserId))
	if ok {
		return data.(int)
	} else {
		return 0
	}
}

func SetFightTime(UserId, unixTime int) {
	GetGCache().Set(FIGHTTIME+strconv.Itoa(UserId), unixTime)
}

// 常用
func GetFightCoolTime(userId, unixTime int) int {
	t := GetFightTime(userId)
	if coolTime := t + FIGHTTIMECOOLTIME - unixTime; coolTime > 0 {
		return coolTime
	}
	return 0
}

func SetAttackTime(UserId, unixTime int) {
	GetGCache().Set(ATTACKTIME+strconv.Itoa(UserId), unixTime)
}

func GetAttackCoolTime(UserId, unixTime int) int {
	var t int
	data, ok := GetGCache().Get(ATTACKTIME + strconv.Itoa(UserId))
	if ok {
		t = data.(int)
	} else {
		return 0
	}
	if coolTime := t + ATTACKTIMECOOLTIME - unixTime; coolTime > 0 {
		return coolTime
	}
	return 0
}

func SetInMap(userId, mapId int) {
	if oldId := GetInMap(userId); oldId != mapId {
		if oldId != 0 {
			DelMapUser(userId, oldId)
		}
		gCache.Set(INMAP+strconv.Itoa(userId), mapId)
		if mapId != 0 {
			setMapUser(userId, mapId)
		}
	}
	gCache.Set(INMAPTIME+strconv.Itoa(userId), utils.NowUnix())
}

func GetInMap(userId int) int {
	if data, ok := gCache.Get(INMAP + strconv.Itoa(userId)); ok {
		return data.(int)
	}
	return 0
}

func GetInMapTime(userId int) int {
	if data, ok := gCache.Get(INMAPTIME + strconv.Itoa(userId)); ok {
		return data.(int)
	}
	return 0
}

type PetStatus struct {
	DeHp int
	DeMp int
}

// 玩家战斗对象怪物血量
type FightStatus struct {
	GpcId    int
	Multiple int
	DeHp     int
}

func GetPetStatus(petId int) *PetStatus {
	if data, ok := gCache.Get(PETSTATUS + strconv.Itoa(petId)); ok {
		bs, ok := data.([]byte)
		if ok {
			status := &PetStatus{}
			err := json.Unmarshal(bs, status)
			if err == nil {
				return status
			}
		}
	}
	return nil
}

func SetPetStatus(petId, deHp, deMp int) {
	status := GetPetStatus(petId)
	if status == nil {
		status = &PetStatus{}
	}
	status.DeHp = deHp
	status.DeMp = deMp
	bs, _ := json.Marshal(status)
	gCache.Set(PETSTATUS+strconv.Itoa(petId), bs)
}

func DelPetStatus(petId int) {
	gCache.Del(PETSTATUS + strconv.Itoa(petId))
}

var fightStatusHub = make([]*FightStatus, 100000)

func GetFightStatus(userId int) *FightStatus {
	if userId < len(fightStatusHub) {
		data := fightStatusHub[userId]
		return data
	}
	return nil
}

func SetFightStatus(userId, gpcId, multiple, deHp int) {
	if userId < len(fightStatusHub) {
		status := GetFightStatus(userId)
		if status == nil {
			status = &FightStatus{}
		}
		status.GpcId = gpcId
		status.Multiple = multiple
		status.DeHp = deHp
		fightStatusHub[userId] = status
	}
}

func DelFightStatus(userId int) {
	if userId < len(fightStatusHub) {
		fightStatusHub[userId] = nil
	}
}

func SetAutoFightFlag(userId, autoFlag int) {
	// 0为非自动，1为金币自动，2为元宝自动
	GetGCache().Set(USERAUTOFIGHTFLAG+strconv.Itoa(userId), autoFlag)
}

func GetAutoFightFlag(userId int) int {
	if data, ok := GetGCache().Get(USERAUTOFIGHTFLAG + strconv.Itoa(userId)); ok {
		return data.(int)
	}
	return AutoFightNone
}

func DelAutoFightFlag(userId int) {
	GetGCache().Del(USERAUTOFIGHTFLAG + strconv.Itoa(userId))
}

func SetAutoGroupFightFlag(userId, autoFlag int) {
	// 0为非自动，1为金币自动，2为元宝自动
	GetGCache().Set(USERAUTOFIGHTFLAG+"group"+strconv.Itoa(userId), autoFlag)
}

func GetAutoGroupFightFlag(userId int) int {
	if data, ok := GetGCache().Get(USERAUTOFIGHTFLAG + "group" + strconv.Itoa(userId)); ok {
		return data.(int)
	}
	return AutoFightNone
}

func DelAutoGroupFightFlag(userId int) {
	GetGCache().Del(USERAUTOFIGHTFLAG + "group" + strconv.Itoa(userId))
}

func SetAutoFightSkill(userId, skill int) {
	GetGCache().Set(USERAUTOFIGHTSKILL+strconv.Itoa(userId), skill)
}

func GetAutoFightSkill(userId int) int {
	if data, ok := GetGCache().Get(USERAUTOFIGHTSKILL + strconv.Itoa(userId)); ok {
		return data.(int)
	}
	return 0
}

func DelAutoFightSkill(userId int) {
	GetGCache().Del(USERAUTOFIGHTSKILL + strconv.Itoa(userId))
}

func SetTTFlag(userId, autoFlag int) {
	RdbOperator.Hset(USERTTFLAG, userId, autoFlag)
}

func GetTTFlag(userId int) int {
	if data, err := RdbOperator.Hget(USERTTFLAG, userId).Int(); err == nil {
		return data
	}
	return UserTTNone
}

func DelTTFlag(userId int) {
	RdbOperator.Hdel(USERTTFLAG, userId)
}

func GetTTRecord(userId int) *models.TTRecord {
	record := &models.TTRecord{}
	err := RdbOperator.Hget(USERTTRECORD, userId).Struct(record)
	if err != nil {
		fmt.Printf("还原结构体失败！err:%s\n", err)
		return nil
	}
	return record
}

func SetTTRecord(userId int, record *models.TTRecord) {
	if record == nil {
		return
	}
	data, _ := marshal(record)
	RdbOperator.Hset(USERTTRECORD, userId, data)
}

func DelTTRecord(userId int) {
	RdbOperator.Hdel(USERTTRECORD, userId)
}

func SetYiWangTime(userId, t int) {
	RdbOperator.Hset(YIWANGFLAG, userId, t)
}

func GetYiWangTime(userId int) int {
	t, _ := RdbOperator.Hget(YIWANGFLAG, userId).Int()
	return t
}

func GetSSBattleEndTime() int {
	t, _ := RdbOperator.Get(SSBattleEndTime).Int()
	return t
}

func SetSSBattleEndTime(t int) {
	RdbOperator.Set(SSBattleEndTime, t)
}

func GetSSBattleStartTime() int {
	t, _ := RdbOperator.Get(SSBattleStartTime).Int()
	return t
}

func SetSSBattleStartTime(t int) {
	RdbOperator.Set(SSBattleStartTime, t)
}
