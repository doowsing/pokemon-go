package datastore

import "pokemon/pkg/models"

var Growth = &growthStore{newStore()}

type growthStore struct {
	*store
}

func (s *growthStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.Growth
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.Wx, &v1)
	}
}

func (s *growthStore) Get(id int) *models.Growth {
	if data, exist := s.store.data.Load(id); exist {
		return data.(*models.Growth)
	}
	return nil
}
