package datastore

import (
	"pokemon/game/models"
)

var Growth = &growthStore{store: newStore(), data: make([]*models.Growth, 10)}

type growthStore struct {
	*store
	data  []*models.Growth
	_data []*models.Growth
}

func (s *growthStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make([]*models.Growth, 10)
	copy(s._data, s.data)

	s.updating = true
	s.data = make([]*models.Growth, 10)

	var all []*models.Growth
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.Wx] = v
	}
	s.updating = false
}

func (s *growthStore) getData() []*models.Growth {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *growthStore) Get(id int) *models.Growth {
	if id < len(s.getData()) {
		return s.getData()[id]
	}
	return nil
}
