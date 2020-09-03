package datastore

import (
	"pokemon/pkg/models"
)

var CardTitle = &cardTitleStore{newStore()}

type cardTitleStore struct {
	*store
}

func (s *cardTitleStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.CardTitle
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.CodeName, &v1)
	}
}

func (s *cardTitleStore) Get(CodeName string) *models.CardTitle {
	if data, exist := s.store.data.Load(CodeName); exist {
		return data.(*models.CardTitle)
	}
	return nil
}
