package datastore

import "pokemon/pkg/models"

var taskSlice = make([]*models.Task, 1000)

var Task = &taskStore{store: newStore()}

type taskStore struct {
	*store
	sliceStore []*models.Task
}

func (s *taskStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	s.sliceStore = make([]*models.Task, 10000)
	var all []*models.Task
	db.Find(&all)
	for _, v := range all {
		s.sliceStore[v.Id] = v
	}
}

func (s *taskStore) Get(id int) *models.Task {
	if len(s.sliceStore) >= id {
		return nil
	}
	return s.sliceStore[id]
}
