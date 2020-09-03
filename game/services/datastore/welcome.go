package datastore

import (
	"pokemon/game/models"
)

var Welcome = &welcomeStore{store: newStore(), data: make(map[string]*models.Welcome)}

type welcomeStore struct {
	*store
	data  map[string]*models.Welcome
	_data map[string]*models.Welcome
}

func (s *welcomeStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make(map[string]*models.Welcome)
	clone(s.data, s._data)
	s.updating = true

	s.data = make(map[string]*models.Welcome)

	var all []*models.Welcome
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.Code] = v
	}
	s.updating = false
}

func (s *welcomeStore) getData() map[string]*models.Welcome {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *welcomeStore) Get(code string) *models.Welcome {
	if data, exist := s.getData()[code]; exist {
		return data
	}
	return nil
}

var TimeConfig = &timeConfigStore{store: newStore(), data: make(map[string]*models.TimeConfig)}

type timeConfigStore struct {
	*store
	data  map[string]*models.TimeConfig
	_data map[string]*models.TimeConfig
}

func (s *timeConfigStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make(map[string]*models.TimeConfig)
	clone(s.data, s._data)
	s.updating = true

	s.data = make(map[string]*models.TimeConfig)

	var all []*models.TimeConfig
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.Title] = v
	}
	s.updating = false
}

func (s *timeConfigStore) getData() map[string]*models.TimeConfig {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *timeConfigStore) Get(code string) *models.TimeConfig {
	if data, exist := s.getData()[code]; exist {
		return data
	}
	return nil
}
