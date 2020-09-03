package services

import (
	"github.com/jinzhu/gorm"
	"pokemon/common/persistence"
	"pokemon/game/utils"
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
	NpcSrv   *NpcServices
}

var OptServicePool = &sync.Pool{New: func() interface{} {
	opt := &OptService{}
	opt.FightSrv = NewFightService(opt)
	opt.PetSrv = NewPetService(opt)
	opt.PropSrv = NewPropService(opt)
	opt.UserSrv = NewUserService(opt)
	opt.SysSrv = NewSysService(opt)
	opt.TaskSrv = NewTaskServices(opt)
	opt.NpcSrv = NewNpcServices(opt)
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
	opt.NpcSrv = NewNpcServices(opt)
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

type BaseService struct {
	OptSvc *OptService
}

func (bsrc *BaseService) SetOptSrc(src *OptService) {
	bsrc.OptSvc = src
}

func (bsrc *BaseService) GetDb() *gorm.DB {
	if bsrc.OptSvc != nil {
		return bsrc.OptSvc.GetDb()
	}
	return persistence.GetOrm()
}

func (bsrc *BaseService) NowUnix() int {
	if bsrc.OptSvc != nil {
		return bsrc.OptSvc.NowUnix()
	}
	return utils.NowUnix()
}
