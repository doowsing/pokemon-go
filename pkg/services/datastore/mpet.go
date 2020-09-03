package datastore

import "pokemon/pkg/models"

var Mpet = &mpetStore{newStore()}

type mpetStore struct {
	*store
}

func (s *mpetStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.MPet
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.ID, &v1)
		s.store.data.Store(v.Name, &v1)
	}
}

func (s *mpetStore) Get(id int) *models.MPet {
	if data, exist := s.store.data.Load(id); exist {
		return data.(*models.MPet)
	}
	return nil
}

func (s *mpetStore) GetByName(name string) *models.MPet {
	if data, exist := s.store.data.Load(name); exist {
		return data.(*models.MPet)
	}
	return nil
}
