package datastore

import (
	"pokemon/game/models"
	"sync"
)

var Mprop = &mpropStore{store: newStore(), data: make([]*models.MProp, 10000)}

type mpropStore struct {
	*store
	data     []*models.MProp
	_data    []*models.MProp
	nameData sync.Map
}

func (s *mpropStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make([]*models.MProp, 10000)
	copy(s._data, s.data)

	s.updating = true
	s.data = make([]*models.MProp, 10000)

	var all []*models.MProp
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.ID] = v
		s.nameData.Store(v.Name, v.ID)
	}
	s.updating = false
}

func (s *mpropStore) getData() []*models.MProp {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *mpropStore) Get(id int) *models.MProp {
	if id < len(s.getData()) {
		return s.getData()[id]
	}
	return nil
}

func (s *mpropStore) GetIdsByName(names []string) (ids []int) {
	ids = []int{}
	for _, name := range names {
		id := s.GetIdByName(name)
		if id > 0 {
			ids = append(ids, id)
		}
	}
	return
}

func (s *mpropStore) GetIdByName(name string) (id int) {
	if data, ok := s.nameData.Load(name); ok {
		return data.(int)
	}
	return 0
}
