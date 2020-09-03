package datastore

import "pokemon/game/models"

var ExpList = &expListStore{store: newStore(), data: make(map[int]int)}

type expListStore struct {
	*store
	data  map[int]int
	_data map[int]int
}

func (s *expListStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make(map[int]int)
	clone(&s.data, &s._data)

	s.updating = true
	s.data = make(map[int]int)

	var all []*models.ExpList
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.Level] = v.NextLvExp
	}
	s.updating = false
}

func (s *expListStore) getData() map[int]int {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *expListStore) Get(id int) int {
	if data, exist := s.getData()[id]; exist {
		return data
	}
	return 0
}
