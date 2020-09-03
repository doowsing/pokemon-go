package rcache

import (
	"encoding/json"
	"fmt"
	"pokemon/pkg/common"
	"pokemon/pkg/models"
	"reflect"
	"strconv"
)

type PropRedisRepository struct {
	BaseRedisRepository
}

func NewPropRedisRepository() *PropRedisRepository {
	return &PropRedisRepository{BaseRedisRepository: BaseRedisRepository{}}
}

func (rr *PropRedisRepository) SetMProps(MProps *[]models.MProp) (bool, error) {
	_, _ = Delete(common.MPropKey)
	mpropMap := make(map[string]interface{})
	for _, v := range *MProps {
		if v.VaryName == 9 {
			//v.ZbInfos, _ = rr.EncodeZbInfo(v.ZbInfo)
		}
		mpropMap[strconv.Itoa(v.ID)] = v
	}
	return Hmset(common.MPropKey, &mpropMap)
}

func (rr *PropRedisRepository) SetMSeries(MSeries *[]models.MSeries) (bool, error) {
	_, _ = Delete(common.MSERIESKey)
	mseriesMap := make(map[string]interface{})
	for _, v := range *MSeries {
		mseriesMap[strconv.Itoa(v.ID)] = v
	}
	return Hmset(common.MSERIESKey, &mseriesMap)
}

func (rr *PropRedisRepository) GetMProp(MPropId int) (*models.MProp, error) {
	mpropStr, err := Hget(common.MPropKey, strconv.Itoa(MPropId))
	if err != nil {
		return nil, err
	}
	var mprop models.MProp
	err = json.Unmarshal(mpropStr, &mprop)
	if err != nil {
		fmt.Println("还原结构体失败！", err)
		return nil, err
	}
	return &mprop, err
}

func (rr *PropRedisRepository) GetMSeries(MSeriesId int) (*models.MSeries, error) {
	mseriesStr, err := Hget(common.MSERIESKey, strconv.Itoa(MSeriesId))
	if err != nil {
		return nil, err
	}
	var mseries models.MSeries
	err = json.Unmarshal(mseriesStr, &mseries)
	if err != nil {
		fmt.Println("还原结构体失败！", err)
		return nil, err
	}
	return &mseries, err
}

func (rr *PropRedisRepository) EncodeZbInfo(zbInfo string) (*models.MZbInfo, error) {
	var ZbInfo models.MZbInfo
	var NewInfo map[string]interface{}
	err := json.Unmarshal([]byte(zbInfo), &NewInfo)
	fmt.Print(NewInfo)
	if NewInfo["main_info"] != nil {
		NewFormat := map[string]float64{}
		for k, v := range NewInfo["main_info"].(map[string]interface{}) {
			if reflect.TypeOf(v).Kind() == reflect.Int {
				NewFormat[k] = float64(v.(int))
			} else {
				NewFormat[k] = v.(float64)
			}
		}
		ZbInfo.MainInfo = NewFormat
	}
	if NewInfo["other_info"] != nil {
		NewFormat := map[string]float64{}
		for k, v := range NewInfo["other_info"].(map[string]interface{}) {
			if reflect.TypeOf(v).Kind() == reflect.Int {
				NewFormat[k] = float64(v.(int))
			} else {
				NewFormat[k] = v.(float64)
			}
		}
		ZbInfo.OtherInfo = NewFormat
	}
	if NewInfo["enable_strengthen"] != nil {
		if NewInfo["enable_strengthen"].(float64) > 0 {
			ZbInfo.EnableStrengthen = true
		} else {
			ZbInfo.EnableStrengthen = false
		}
	}
	if NewInfo["strengthen_pid"] != nil && NewInfo["strengthen_pid"] != 0 {
		if NewInfo["strengthen_pid"].(float64) > 0 {
			ZbInfo.StrengthenPid = int(NewInfo["strengthen_pid"].(float64))
		}
	}
	if NewInfo["strengthen_effect"] != nil {
		NewFormat := make([]float64, len(NewInfo["strengthen_effect"].([]interface{})))
		for k, v := range NewInfo["strengthen_effect"].([]interface{}) {
			if reflect.TypeOf(v).Kind() == reflect.Int {
				NewFormat[k] = float64(v.(int))
			} else {
				NewFormat[k] = v.(float64)
			}
		}
		ZbInfo.StrengthenEffect = NewFormat
	}
	if NewInfo["series"] != nil && NewInfo["series"].(float64) != 0 {
		ZbInfo.Series = int(NewInfo["series"].(float64))
	}
	if NewInfo["position"] != nil && NewInfo["position"].(float64) != 0 {
		ZbInfo.Position = int(NewInfo["position"].(float64))
	}
	return &ZbInfo, err
}
