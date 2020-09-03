package datastore

import "pokemon/pkg/models"

var ExpList = &expListStore{newStore()}

type expListStore struct {
	*store
}

func (s *expListStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.ExpList
	db.Find(&all)
	for _, v := range all {
		s.store.data.Store(v.Level, v.NextLvExp)
	}
}

func (s *expListStore) Get(id int) int {
	if data, exist := s.store.data.Load(id); exist {
		return data.(int)
	}
	return 0
}
