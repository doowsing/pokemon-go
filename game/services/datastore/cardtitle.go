package datastore

import (
	"pokemon/game/models"
	"sort"
	"sync"
)

var CardTitle = &cardTitleStore{store: newStore(), data: make(map[string]*models.CardTitle)}
var CardSerie = &cardSeiresStore{store: newStore()}
var CardPrize = &cardPrizeStore{store: newStore()}

type cardTitleStore struct {
	*store
	data  map[string]*models.CardTitle
	_data map[string]*models.CardTitle
}

func (s *cardTitleStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()

	s._data = make(map[string]*models.CardTitle)
	clone(&s.data, &s._data)

	s.updating = true
	s.data = make(map[string]*models.CardTitle)

	var all []*models.CardTitle
	GetDb().Find(&all)
	for _, v := range all {
		s.data[v.CodeName] = v
	}
	s.updating = false
}

func (s *cardTitleStore) getData() map[string]*models.CardTitle {
	if s.updating {
		return s._data
	}
	return s.data
}

func (s *cardTitleStore) Get(CodeName string) *models.CardTitle {
	if data, exist := s.getData()[CodeName]; exist {
		return data
	}
	return nil
}

func (s *cardTitleStore) GetById(id int) *models.CardTitle {
	for _, title := range s.getData() {
		if title.ID == id {
			return title
		}
	}
	return nil
}

func (s *cardTitleStore) GetAll() []*models.CardTitle {
	datas := []*models.CardTitle{}
	for _, t := range s.getData() {
		datas = append(datas, t)
	}
	sort.Slice(datas, func(i, j int) bool {
		return datas[i].ID < datas[j].ID
	})
	return datas
}

type cardSeiresStore struct {
	*store
	data sync.Map
}

func (s *cardSeiresStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var ss []*models.CardSeries
	GetDb().Find(&ss)
	for _, se := range ss {
		s.data.Store(se.ID, se)
	}
}

func (s *cardSeiresStore) Get(id int) *models.CardSeries {
	if data, ok := s.data.Load(id); ok {
		return data.(*models.CardSeries)
	}
	return nil
}

func (s *cardSeiresStore) GetAll() []*models.CardSeries {
	var ss []*models.CardSeries
	s.data.Range(func(key, value interface{}) bool {
		ss = append(ss, value.(*models.CardSeries))
		return true
	})
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].ID < ss[j].ID
	})
	return ss
}

type cardPrizeStore struct {
	*store
	data sync.Map
}

func (s *cardPrizeStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var ss []*models.CardPrize
	GetDb().Find(&ss)
	for _, se := range ss {
		s.data.Store(se.ID, se)
	}
}

func (s *cardPrizeStore) Get(id int) *models.CardPrize {
	if data, ok := s.data.Load(id); ok {
		return data.(*models.CardPrize)
	}
	return nil
}

func (s *cardPrizeStore) GetAll() []*models.CardPrize {
	var ss []*models.CardPrize
	s.data.Range(func(key, value interface{}) bool {
		ss = append(ss, value.(*models.CardPrize))
		return true
	})
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].ID < ss[j].ID
	})
	return ss
}
