package datastore

import "pokemon/game/models"

var SSjh = &ssjhStore{store: newStore(), data: make([]*models.SSjhRule, 1000)}

type ssjhStore struct {
	*store
	data  []*models.SSjhRule
	_data []*models.SSjhRule
}

func (s *ssjhStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make([]*models.SSjhRule, 1000)
	copy(s._data, s.data)

	s.updating = true
	s.data = make([]*models.SSjhRule, 1000)

	var all []*models.SSjhRule
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.PetId] = v
	}
	s.updating = false
}
func (s *ssjhStore) getData() []*models.SSjhRule {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *ssjhStore) Get(id int) *models.SSjhRule {
	if id < len(s.getData()) {
		return s.getData()[id]
	}
	return nil
}
