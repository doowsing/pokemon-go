package repositories

import (
	"github.com/jinzhu/gorm"
	"pokemon/pkg/models"
)

type FightRepositories struct {
	*BaseRepository
}

func NewFightRepositories() *FightRepositories {
	return &FightRepositories{BaseRepository: NewBaseRepository()}
}

func (fr *FightRepositories) GetAllMapFromMysql() (*[]models.Map, bool) {
	Maps := make([]models.Map, 0)
	success := fr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &Maps)
	return &Maps, success
}

func (fr *FightRepositories) GetAllGpcFromMysql() (*[]models.Gpc, bool) {
	MGPCs := make([]models.Gpc, 0)
	success := fr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &MGPCs)
	return &MGPCs, success
}

func (fr *FightRepositories) GetAllGpcGroupFromMysql() (*[]models.GpcGroup, bool) {
	MGPCGroup := make([]models.GpcGroup, 0)
	success := fr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}, &MGPCGroup)
	return &MGPCGroup, success
}

func (fr *FightRepositories) GetFbRecord(UserId, MapId int) *models.RecordFb {
	record := models.RecordFb{}
	ok := fr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("uid = ? and map_id =?", UserId, MapId)
	}, &record)
	if !ok {
		return nil
	}
	return &record
}

func (fr *FightRepositories) GetBossRecord(MapId int) *models.RecordBoss {
	record := models.RecordBoss{}
	ok := fr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("map_id =?", MapId)
	}, &record)
	if !ok {
		return nil
	}
	return &record
}
