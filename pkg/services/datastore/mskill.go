package datastore

import (
	"pokemon/pkg/models"
	"strconv"
)

var Mskill = &mskillStore{newStore()}

type mskillStore struct {
	*store
}

func (s *mskillStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.MSkill
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.ID, &v1)
		if v.Pid > 0 {
			s.store.data.Store("pid"+strconv.Itoa(v.Pid), &v1)
		}
	}
}

func (s *mskillStore) Get(id int) *models.MSkill {
	if data, exist := s.store.data.Load(id); exist {
		return data.(*models.MSkill)
	}
	return nil
}

func (s *mskillStore) GetByPid(pid int) *models.MSkill {
	if data, exist := s.store.data.Load("pid" + strconv.Itoa(pid)); exist {
		return data.(*models.MSkill)
	}
	return nil
}
