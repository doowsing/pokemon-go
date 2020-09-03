package datastore

import (
	"pokemon/pkg/models"
)

var Welcome = &welcomeStore{newStore()}

type welcomeStore struct {
	*store
}

func (s *welcomeStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.Welcome
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.Code, &v1)
	}
}

func (s *welcomeStore) Get(code string) *models.Welcome {
	if data, exist := s.store.data.Load(code); exist {
		return data.(*models.Welcome)
	}
	return nil
}

var TimeConfig = &timeConfigStore{newStore()}

type timeConfigStore struct {
	*store
}

func (s *timeConfigStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.TimeConfig
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.Title, &v1)
	}
}

func (s *timeConfigStore) Get(code string) *models.TimeConfig {
	if data, exist := s.store.data.Load(code); exist {
		return data.(*models.TimeConfig)
	}
	return nil
}
