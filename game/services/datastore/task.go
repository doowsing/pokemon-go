package datastore

import "pokemon/game/models"

var Task = &taskStore{store: newStore(), data: make([]*models.Task, 1000)}

type taskStore struct {
	*store
	data  []*models.Task
	_data []*models.Task
}

func (s *taskStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make([]*models.Task, 1000)
	copy(s._data, s.data)

	s.updating = true

	s.data = make([]*models.Task, 10000)
	var all []*models.Task
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.Id] = v
	}
	s.updating = false
}

func (s *taskStore) getData() []*models.Task {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *taskStore) Get(id int) *models.Task {
	if id < len(s.getData()) {
		return s.getData()[id]
	}
	return s.data[id]
}
