package repositories

import (
	"github.com/jinzhu/gorm"
	"pokemon/game/models"
)

type PropRepository struct {
	*BaseRepository
}

func NewPropRepository() *PropRepository {
	return &PropRepository{BaseRepository: NewBaseRepository()}
}

func (pr *PropRepository) GetAllMPropFromMysql() (*[]models.MProp, bool) {
	MProps := make([]models.MProp, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &MProps)
	return &MProps, success
}

func (pr *PropRepository) GetAllMSeriesFromMysql() (*[]models.MSeries, bool) {
	MSeries := make([]models.MSeries, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &MSeries)
	return &MSeries, success
}

func (pr *PropRepository) GetBpProps(UserId int) (*[]models.UProp, bool) {
	Props := make([]models.UProp, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("uid = ? and bp_sum>0", UserId)
	}, &Props)
	return &Props, success
}

func (pr *PropRepository) GetBpPropsCount(UserId int) int {
	_, count := pr.Count(&models.UProp{}, func(db *gorm.DB) *gorm.DB {
		return db.Where("uid = ? and bp_sum>0", UserId)
	})
	return count
}

func (pr *PropRepository) GetPropById(PropId int) *models.UProp {
	prop := &models.UProp{}
	ok := pr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ? ", PropId)
	}, prop)
	if ok {
		return prop
	} else {
		return nil
	}
}

func (pr *PropRepository) GetPropByPId(UserId, PropId int) *models.UProp {
	prop := &models.UProp{}
	ok := pr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("pid = ? and uid = ?", PropId, UserId)
	}, prop)
	if ok {
		return prop
	} else {
		return nil
	}
}

func (pr *PropRepository) GetOnePropByPIdList(UserId int, PropIdList string) *models.UProp {
	prop := &models.UProp{}
	_ = pr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("pid in (?) and uid = ?", PropIdList, UserId)
	}, prop)
	return prop
}

func (pr *PropRepository) CheckPropByPid(UserId, PropId, Num int) bool {
	prop := &models.UProp{}
	ok := pr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("pid = ? and uid = ? and bp_sum > ?", PropId, UserId, Num)
	}, prop)
	return ok
}

func (pr *PropRepository) UsePropByPid(UserId, PropId, Num int) *gorm.DB {
	return pr.UpdateWhere(func(db *gorm.DB) *gorm.DB {
		return db.Where("bp_sum > ?", Num)
	}, &models.UProp{Pid: PropId, Uid: UserId}, map[string]interface{}{"bp_sum": gorm.Expr("bp_sum - ?", Num)})

}

func (pr *PropRepository) SetPropSum(PropId, sum int) *gorm.DB {
	return pr.Update(&models.UProp{ID: PropId},
		map[string]interface{}{"bp_sum": sum})
}

func (pr *PropRepository) OffZb(ZbId int) *gorm.DB {
	return pr.Update(&models.UProp{ID: ZbId},
		map[string]interface{}{"bp_sum": 1, "zbpet": nil})
}

func (pr *PropRepository) SetZbPet(ZbId, PetId int) *gorm.DB {
	return pr.Update(&models.UProp{ID: ZbId},
		map[string]interface{}{"zb_pet": PetId})
}

func (pr *PropRepository) UseProp(UserId, PropID, Num int) *gorm.DB {
	result := pr.UpdateWhere(func(db *gorm.DB) *gorm.DB {
		return db.Where("bp_sum >= ?", Num)
	}, &models.UProp{ID: PropID, Uid: UserId},
		map[string]interface{}{"bp_sum": gorm.Expr("bp_sum-?", Num)})
	return result
}
