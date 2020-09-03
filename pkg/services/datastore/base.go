package datastore

import (
	"pokemon/pkg/persistence"
	"sync"
	"time"
)

var db = persistence.GetOrm()

type store struct {
	lastUpTime int
	data       sync.Map
}

func newStore() *store {
	return &store{
		lastUpTime: 0,
		data:       sync.Map{},
	}
}

func (s *store) upTime() {
	s.lastUpTime = int(time.Now().Unix())
}

func (s *store) checkTime(newTime int) bool {
	return newTime > s.lastUpTime
}
