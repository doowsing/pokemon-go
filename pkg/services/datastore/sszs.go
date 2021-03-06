package datastore

import "pokemon/pkg/models"

var SSzs = &sszsStore{newStore()}

type sszsStore struct {
	*store
}

func (s *sszsStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.SSzsRule
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.PetId, &v1)
	}
}

func (s *sszsStore) Get(id int) *models.SSzsRule {
	if data, exist := s.store.data.Load(id); exist {
		return data.(*models.SSzsRule)
	}
	return nil
}
