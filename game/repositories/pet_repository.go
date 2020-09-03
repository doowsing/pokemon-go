package repositories

import (
	"github.com/jinzhu/gorm"
	"pokemon/game/models"
)

type PetRepositories struct {
	*BaseRepository
}

func NewPetRepositories() *PetRepositories {
	return &PetRepositories{BaseRepository: NewBaseRepository()}
}

func (pr *PetRepositories) GetPet(PetId int) *models.UPet {
	pet := &models.UPet{}
	ok := pr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ? ", PetId)
	}, pet)
	if ok {
		return pet
	} else {
		return nil
	}
}

func (pr *PetRepositories) GetPets(UserId int) []*models.UPet {
	return make([]*models.UPet, 0)
}

func (pr *PetRepositories) GetCarryPets(UserId int) *[]models.UPet {
	carryPets := make([]models.UPet, 0)
	_ = pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("uid = ? and location = ?", UserId, 0)
	}, &carryPets)
	return &carryPets
}

func (pr *PetRepositories) GetPetsCount(UserId int) int {
	return 0
}

func (pr *PetRepositories) GetCarryPetsCount(UserId int) int {
	_, count := pr.Count(&models.UPet{}, func(db *gorm.DB) *gorm.DB {
		return db.Where("uid = ? and location = ?", UserId, 0)
	})
	return count
}

//func (pr *PetRepositories) SavePetAttribute(pet *models.UPet) *gorm.DB {
//	return pr.Update(pet,
//		map[string]interface{}{"level": pet.Level, "now_exp": pet.NowExp, "attribute_hm": pet.AttributeHM, "attribute_other": pet.AttributeOt})
//}
//
//func (pr *PetRepositories) SavePetAttributeOther(pet *models.UPet) *gorm.DB {
//	return pr.Update(pet,
//		map[string]interface{}{"attribute_other": pet.AttributeOt})
//}
//
//func (pr *PetRepositories) SavePetAttributeHM(pet *models.UPet) *gorm.DB {
//	return pr.Update(pet,
//		map[string]interface{}{"attribute_hm": pet.AttributeHM})
//}
func (pr *PetRepositories) SavePetCzl(pet *models.UPet) *gorm.DB {
	return pr.Update(pet,
		map[string]interface{}{"czl": pet.Czl})
}

func (pr *PetRepositories) SetPetZbList(PetId int, ZbJson string) *gorm.DB {
	return pr.Update(&models.UPet{ID: PetId},
		map[string]interface{}{"equip": ZbJson})

}

// 元模型数据操作
// 宠物，技能，升级经验列表

func (pr *PetRepositories) GetALLMPetsFromMysql() (*[]models.MPet, bool) {
	MPets := make([]models.MPet, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &MPets)
	return &MPets, success
}

func (pr *PetRepositories) GetALLMSkillsFromMysql() (*[]models.MSkill, bool) {
	MSkills := make([]models.MSkill, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &MSkills)
	return &MSkills, success
}

func (pr *PetRepositories) GetALLExp2lvFromMysql() (*[]models.ExpList, bool) {
	ExpList := make([]models.ExpList, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &ExpList)
	return &ExpList, success
}

func (pr *PetRepositories) GetALLGrowthFromMysql() (*[]models.Growth, bool) {
	GrowthList := make([]models.Growth, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &GrowthList)
	return &GrowthList, success
}

func (pr *PetRepositories) GetALLSSZsFromMysql() (*[]models.SSzsRule, bool) {
	RoleList := make([]models.SSzsRule, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &RoleList)
	return &RoleList, success
}

func (pr *PetRepositories) GetALLSSJhFromMysql() (*[]models.SSjhRule, bool) {
	RoleList := make([]models.SSjhRule, 0)
	success := pr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &RoleList)
	return &RoleList, success
}
