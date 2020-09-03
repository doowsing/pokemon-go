package rcache

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"pokemon/pkg/common"
	"pokemon/pkg/models"
	"strconv"
	"time"
)

type FightRedisRepository struct {
	BaseRedisRepository
}

func NewFightRedisRepository() *FightRedisRepository {
	return &FightRedisRepository{BaseRedisRepository: BaseRedisRepository{}}
}
func (fr *FightRedisRepository) SetMMaps(maps *[]models.Map) (bool, error) {
	_, _ = Delete(common.MMapKey)

	mapMap := make(map[string]interface{})
	for _, v := range *maps {
		//fr.UnSerializerMap(&v)
		mapMap[strconv.Itoa(v.ID)] = v
	}
	return Hmset(common.MMapKey, &mapMap)
}

func (fr *FightRedisRepository) SetMGpcs(mgpcs *[]models.Gpc) (bool, error) {
	_, _ = Delete(common.MGpcKey)
	gpcmap := make(map[string]interface{})
	for _, v := range *mgpcs {
		//fr.UnSerializerGpc(&v)
		gpcmap[strconv.Itoa(v.ID)] = v
	}
	return Hmset(common.MGpcKey, &gpcmap)
}

func (fr *FightRedisRepository) SetMGpcGroup(mgpcGroup *[]models.GpcGroup) (bool, error) {
	_, _ = Delete(common.MGpcGroupKey)
	gpcGroupMap := make(map[string]interface{})
	for _, v := range *mgpcGroup {
		//fr.UnSerializerGpcGroup(&v)
		gpcGroupMap[strconv.Itoa(v.ID)] = v
	}
	return Hmset(common.MGpcGroupKey, &gpcGroupMap)
}

func (fr *FightRedisRepository) UnSerializerGpcGroup(gpcGroup *models.GpcGroup) {

	drops := make([][]int, 0)
	err := json.Unmarshal([]byte(gpcGroup.DropList), &drops)
	if err != nil {
		return
	}
	gpcGroup.Drops = &drops

	gpcs := make([]int, 0)
	err = json.Unmarshal([]byte(gpcGroup.GpcList), &gpcs)
	if err != nil {
		return
	}
	gpcGroup.Gpcs = &gpcs
}

func (fr *FightRedisRepository) GetMap(id int) (*models.Map, error) {
	mapStr, err := Hget(common.MMapKey, strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	var mmap models.Map
	err = json.Unmarshal(mapStr, &mmap)
	if err != nil {
		fmt.Println("还原结构体失败！", err)
		return nil, err
	}
	return &mmap, err
}

func (fr *FightRedisRepository) GetGpc(id int) (*models.Gpc, error) {
	gpcStr, err := Hget(common.MGpcKey, strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	var mgpc models.Gpc
	err = json.Unmarshal(gpcStr, &mgpc)
	if err != nil {
		fmt.Println("还原结构体失败！", err)
		return nil, err
	}
	return &mgpc, err
}

func (fr *FightRedisRepository) GetGpcGroup(id int) (*models.GpcGroup, error) {
	gpcGroupStr, err := Hget(common.MGpcGroupKey, strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	var mgpcGroup models.GpcGroup
	err = json.Unmarshal(gpcGroupStr, &mgpcGroup)
	if err != nil {
		fmt.Println("还原结构体失败！", err)
		return nil, err
	}
	return &mgpcGroup, err
}

func (fr *FightRedisRepository) GetZbAttribute(PetId int) (*models.ZbAttribute, error) {
	Str, err := Hget(common.MZbAttributeKey, strconv.Itoa(PetId))
	if err != nil {
		return nil, err
	}
	var attribute *models.ZbAttribute
	err = json.Unmarshal(Str, &attribute)
	if err != nil {
		return nil, err
	}
	return attribute, nil
}

func (fr *FightRedisRepository) SetZbAttribute(PetId int, zbAttribute *models.ZbAttribute) {
	_, _ = Hset(common.MZbAttributeKey, strconv.Itoa(PetId), zbAttribute)
}

func (fr *FightRedisRepository) ClearZbAttribute(PetId int) {
	_, _ = Hdel(common.MZbAttributeKey, strconv.Itoa(PetId))
}

func (fr *FightRedisRepository) SetFightTime(UserId int) {
	_, _ = Hset(common.MFightTimeKey, strconv.Itoa(UserId), time.Now().Unix())
}

func (fr *FightRedisRepository) GetFightTime(UserId int) int {
	_time, err := redis.Int(Hget(common.MFightTimeKey, strconv.Itoa(UserId)))
	if err != nil {
		return 0
	}
	return _time
}
