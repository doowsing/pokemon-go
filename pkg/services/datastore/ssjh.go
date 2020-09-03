package datastore

import "pokemon/pkg/models"

var SSjh = &ssjhStore{newStore()}

type ssjhStore struct {
	*store
}

func (s *ssjhStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.SSjhRule
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.PetId, &v1)
	}
}

func (s *ssjhStore) Get(id int) *models.SSjhRule {
	if data, exist := s.store.data.Load(id); exist {
		return data.(*models.SSjhRule)
	}
	return nil
}
