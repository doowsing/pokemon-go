package rcache

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"pokemon/pkg/common"
	"pokemon/pkg/models"
	"strconv"
)

//Redis缓存管理, 文章之类的 不细分
type PetRedisRepository struct {
	BaseRedisRepository
}

func NewPetRedisRepository() *PetRedisRepository {
	return &PetRedisRepository{BaseRedisRepository: BaseRedisRepository{}}
}

func (rr *PetRedisRepository) GetMPet(MPetId int) (*models.MPet, error) {
	mpetStr, err := Hget(common.MPetKey, strconv.Itoa(MPetId))
	if err != nil {
		return nil, err
	}
	var mpet models.MPet
	err = json.Unmarshal(mpetStr, &mpet)
	if err != nil {
		fmt.Println("还原结构体失败！", err)
		return nil, err
	}
	return &mpet, err
}

func (rr *PetRedisRepository) SetMPets(MPets *[]models.MPet) (interface{}, error) {
	_, _ = Delete(common.MPetKey)
	mpetsMap := make(map[string]interface{})
	for _, v := range *MPets {
		mpetsMap[strconv.Itoa(v.ID)] = v
	}
	return Hmset(common.MPetKey, &mpetsMap)
}

func (rr *PetRedisRepository) GetMSkill(MSkillId int) (*models.MSkill, error) {
	mskillStr, err := Hget(common.MSkillKey, strconv.Itoa(MSkillId))
	if err != nil {
		return nil, err
	}
	var mskill models.MSkill

	err = json.Unmarshal(mskillStr, &mskill)
	if err != nil {
		fmt.Println("还原结构体失败！", err)
		return nil, err
	}
	return &mskill, err
}

func (rr *PetRedisRepository) SetMSkill(MSkills *[]models.MSkill) (bool, error) {
	_, _ = Delete(common.MSkillKey)
	mskillMap := make(map[string]interface{})
	for _, v := range *MSkills {
		mskillMap[strconv.Itoa(v.ID)] = v
	}
	return Hmset(common.MSkillKey, &mskillMap)
}

func (rr *PetRedisRepository) GetExpByLv(level int) (int, error) {
	return redis.Int(Hget(common.ExpListKey, strconv.Itoa(level)))
}
func (rr *PetRedisRepository) GetExps() (*[]models.ExpList, error) {
	ExpListStr, err := Hget(common.ExpListKey, common.ExpListKey)
	if err != nil {
		return nil, err
	}
	var ExpList []models.ExpList
	err = json.Unmarshal(ExpListStr, &ExpList)
	if err != nil {
		fmt.Println("还原技能结构体失败！", err)
		return nil, err
	}
	return &ExpList, err
}

func (rr *PetRedisRepository) SetExps(ExpList *[]models.ExpList) (bool, error) {
	_, _ = Delete(common.ExpListKey)
	expList := make(map[string]interface{})

	for _, v := range *ExpList {
		expList[strconv.Itoa(v.Level)] = v.NextLvExp
	}
	vJson, err := json.Marshal(*ExpList)
	if err != nil {
		fmt.Println("序列化经验表模型失败！")
	} else {
		expList[common.ExpListKey] = vJson
	}
	return Hmset(common.ExpListKey, &expList)
}
