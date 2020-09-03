package rcache

import (
	"pokemon/game/common"
	"pokemon/game/models"
	"strconv"
)

const (
	PMS              = "PMS"
	PMSCOOLTIME      = 2
	SMSHOP           = "SMSHOP"
	SMSHOPCOOLTIME   = 2
	PROPSHOW         = "PROPSHOW"
	PROPSHOWCOOLTIME = 2
)

var (
	PmsTimer      = NewTimeManager(PMS, PMSCOOLTIME)
	SmShopTimer   = NewTimeManager(SMSHOP, SMSHOPCOOLTIME)
	PropShowTimer = NewTimeManager(PROPSHOW, PROPSHOWCOOLTIME)
)

func SetMProps(MProps *[]models.MProp) error {
	RdbOperator.Delete(common.MPropKey)
	mpropMap := make(map[string]interface{})
	for _, v := range *MProps {
		mpropMap[strconv.Itoa(v.ID)] = v
	}
	return RdbOperator.Hmset(common.MPropKey, mpropMap).Error()
}

func GetMProp(MPropId int) (*models.MProp, error) {
	mprop := models.MProp{}
	err := RdbOperator.Hget(common.MPropKey, MPropId).Struct(&mprop)
	if err != nil {
		return nil, err
	}
	return &mprop, nil
}
