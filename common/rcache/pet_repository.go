package rcache

import (
	"pokemon/game/common"
	"pokemon/game/models"
	"strconv"
)

const (
	EVOLUTION           = "EVOLUTION"
	EVOLUTIONCOOLTIME   = 2
	MERGE               = "MERGE"
	MERGECOOLTIME       = 10
	ZHUANGSEHNG         = "ZHUANGSEHNG"
	ZHUANGSEHNGCOOLTIME = 16
)

var (
	EvolutionTimer   = NewTimeManager(EVOLUTION, EVOLUTIONCOOLTIME)
	MergeTimer       = NewTimeManager(MERGE, MERGECOOLTIME)
	ZhuangshengTimer = NewTimeManager(ZHUANGSEHNG, ZHUANGSEHNGCOOLTIME)
)

func GetMPet(MPetId int) (*models.MPet, error) {
	mpet := models.MPet{}
	err := RdbOperator.Hget(common.MPetKey, MPetId).Struct(&mpet)
	if err != nil {
		return &mpet, err
	}
	return nil, err
}

func SetMPets(MPets *[]models.MPet) error {
	RdbOperator.Delete(common.MPetKey)
	mpetsMap := make(map[string]interface{})
	for _, v := range *MPets {
		mpetsMap[strconv.Itoa(v.ID)] = v
	}
	return RdbOperator.Hmset(common.MPetKey, mpetsMap).Error()
}

func GetMSkill(MSkillId int) (*models.MSkill, error) {
	mskill := models.MSkill{}
	err := RdbOperator.Hget(common.MSkillKey, MSkillId).Struct(&mskill)
	if err != nil {
		return &mskill, err
	}
	return nil, err
}

func SetMSkill(MSkills *[]models.MSkill) error {
	RdbOperator.Delete(common.MSkillKey)
	mskillMap := make(map[string]interface{})
	for _, v := range *MSkills {
		mskillMap[strconv.Itoa(v.ID)] = v
	}
	return RdbOperator.Hmset(common.MSkillKey, mskillMap).Error()
}

func GetExpByLv(level int) (int, error) {
	return RdbOperator.Hget(common.ExpListKey, strconv.Itoa(level)).Int()
}

func GetExps() ([]models.ExpList, error) {
	var ExpList []models.ExpList
	err := RdbOperator.Hget(common.ExpListKey, common.ExpListKey).Struct(&ExpList)
	return ExpList, err
}

func SetExps(ExpList *[]models.ExpList) {
	RdbOperator.Delete(common.ExpListKey)
	expList := make(map[string]interface{})

	for _, v := range *ExpList {
		expList[strconv.Itoa(v.Level)] = v.NextLvExp
	}
	expList[common.ExpListKey] = *ExpList
	RdbOperator.Hmset(common.ExpListKey, expList)
}
