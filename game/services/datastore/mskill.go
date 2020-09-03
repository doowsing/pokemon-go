package datastore

import (
	"pokemon/game/models"
)

var Mskill = &mskillStore{store: newStore(), data: make([]*models.MSkill, 10000), dataPid: make(map[int]*models.MSkill)}

type mskillStore struct {
	*store
	data     []*models.MSkill
	_data    []*models.MSkill
	dataPid  map[int]*models.MSkill
	_dataPid map[int]*models.MSkill
}

func (s *mskillStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make([]*models.MSkill, 10000)
	s._dataPid = make(map[int]*models.MSkill)
	copy(s._data, s.data)
	clone(&s.dataPid, &s._dataPid)

	s.updating = true
	s.data = make([]*models.MSkill, 10000)
	s.dataPid = make(map[int]*models.MSkill)

	var all []*models.MSkill
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.ID] = v
		if v.Pid > 0 {
			s.dataPid[v.Pid] = v
		}
	}
	s.updating = false
}

func (s *mskillStore) getData() []*models.MSkill {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *mskillStore) getDataPid() map[int]*models.MSkill {
	if s.updating {
		return s._dataPid
	}
	return s.dataPid
}

func (s *mskillStore) Get(id int) *models.MSkill {
	if id < len(s.getData()) {
		return s.getData()[id]
	}
	return nil
}

func (s *mskillStore) GetByPid(pid int) *models.MSkill {
	if data, exist := s.getDataPid()[pid]; exist {
		return data
	}
	return nil
}
