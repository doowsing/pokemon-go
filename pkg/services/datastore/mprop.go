package datastore

import "pokemon/pkg/models"

var Mprop = &mpropStore{newStore()}

type mpropStore struct {
	*store
}

func (s *mpropStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.MProp
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.ID, &v1)
	}
}

func (s *mpropStore) Get(id int) *models.MProp {
	if data, exist := s.store.data.Load(id); exist {
		return data.(*models.MProp)
	}
	return nil
}
