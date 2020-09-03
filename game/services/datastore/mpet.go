package datastore

import "pokemon/game/models"

var Mpet = &mpetStore{store: newStore(), data: make([]*models.MPet, 1000)}

type mpetStore struct {
	*store
	data      []*models.MPet
	_data     []*models.MPet
	dataName  map[string]int
	_dataName map[string]int
}

func (s *mpetStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make([]*models.MPet, 1000)
	s._dataName = make(map[string]int)
	copy(s._data, s.data)
	clone(&s.dataName, &s._dataName)

	s.updating = true
	s.data = make([]*models.MPet, 1000)
	s.dataName = make(map[string]int)

	var all []*models.MPet
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.ID] = v
		s.dataName[v.Name] = v.ID
	}
	s.updating = false

}

func (s *mpetStore) getData() []*models.MPet {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *mpetStore) getDataName() map[string]int {
	if s.updating {
		return s._dataName
	}
	return s.dataName
}

func (s *mpetStore) Get(id int) *models.MPet {
	if id < len(s.getData()) {
		return s.getData()[id]
	}
	return nil
}

func (s *mpetStore) GetByName(name string) *models.MPet {

	if id, exist := s.getDataName()[name]; exist {
		return s.Get(id)
	}
	return nil
}
