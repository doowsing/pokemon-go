package services

import (
	"github.com/jinzhu/gorm"
	"pokemon/pkg/persistence"
	"pokemon/pkg/utils"
	"sync"
)

type UpMap map[string]interface{}

type OptService struct {
	db       *gorm.DB
	inTs     bool
	onlyTime int

	FightSrv *FightService
	PetSrv   *PetService
	UserSrv  *UserService
	PropSrv  *PropService
	SysSrv   *SysService
	TaskSrv  *TaskServices
}

var OptServicePool = &sync.Pool{New: func() interface{} {
	opt := &OptService{}
	opt.FightSrv = NewFightService(opt)
	opt.PetSrv = NewPetService(opt)
	opt.PropSrv = NewPropService(opt)
	opt.UserSrv = NewUserService(opt)
	opt.SysSrv = NewSysService(opt)
	opt.TaskSrv = NewTaskServices(opt)
	return opt
}}

func NewOptService() *OptService {
	//opt := OptServicePool.Get().(*OptService)
	opt := &OptService{}
	opt.FightSrv = NewFightService(opt)
	opt.PetSrv = NewPetService(opt)
	opt.PropSrv = NewPropService(opt)
	opt.UserSrv = NewUserService(opt)
	opt.SysSrv = NewSysService(opt)
	opt.TaskSrv = NewTaskServices(opt)
	opt.db = persistence.GetOrm()
	return opt
}

func (ts *OptService) ReSet() {

	ts.onlyTime = utils.NowUnix()
	ts.inTs = false
}

func DropOptService(opt *OptService) {
	OptServicePool.Put(opt)
}

func (ts *OptService) GetDb() *gorm.DB {
	return ts.db
}

func (ts *OptService) Begin() {
	ts.db = ts.db.Begin()
	ts.inTs = true
}

func (ts *OptService) Rollback() {
	if ts.inTs {
		ts.db.Rollback()
		ts.inTs = false
		ts.db = persistence.GetOrm()
	}

}

func (ts *OptService) Commit() {
	if ts.inTs {
		ts.db.Commit()
		ts.inTs = false
		ts.db = persistence.GetOrm()
	}
}

func (ts *OptService) NowUnix() int {
	if ts.onlyTime != 0 {
		return ts.onlyTime
	}
	return utils.NowUnix()
}

type baseService struct {
	OptSrc *OptService
}

func (bsrc *baseService) SetOptSrc(src *OptService) {
	bsrc.OptSrc = src
}

func (bsrc *baseService) GetDb() *gorm.DB {
	if bsrc.OptSrc != nil {
		return bsrc.OptSrc.GetDb()
	}
	return persistence.GetOrm()
}

func (bsrc *baseService) NowUnix() int {
	if bsrc.OptSrc != nil {
		return bsrc.OptSrc.NowUnix()
	}
	return utils.NowUnix()
}
