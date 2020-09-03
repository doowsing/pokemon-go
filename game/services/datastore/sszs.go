package datastore

import "pokemon/game/models"

var SSzs = &sszsStore{store: newStore(), data: make([]*models.SSzsRule, 1000)}

type sszsStore struct {
	*store
	data  []*models.SSzsRule
	_data []*models.SSzsRule
}

func (s *sszsStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make([]*models.SSzsRule, 1000)
	copy(s._data, s.data)

	s.updating = true
	s.data = make([]*models.SSzsRule, 1000)

	var all []*models.SSzsRule
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.PetId] = v
	}
	s.updating = false
}
func (s *sszsStore) getData() []*models.SSzsRule {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *sszsStore) Get(id int) *models.SSzsRule {
	if id < len(s.getData()) {
		return s.getData()[id]
	}
	return nil
}
