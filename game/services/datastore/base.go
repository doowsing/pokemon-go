package datastore

import (
	"bytes"
	"encoding/gob"
	"github.com/jinzhu/gorm"
	"pokemon/common/persistence"
	"time"
)

type store struct {
	updating   bool
	lastUpTime int
}

func GetDb() *gorm.DB {
	return persistence.GetOrm()
}

func newStore() *store {
	return &store{
		lastUpTime: 0,
	}
}

func (s *store) upTime() {
	s.lastUpTime = int(time.Now().Unix())
}

func (s *store) checkTime(newTime int) bool {

	return newTime > s.lastUpTime
}

// Clone 完整复制数据，需传地址
func clone(a, b interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(a); err != nil {
		return err
	}
	if err := dec.Decode(b); err != nil {
		return err
	}
	return nil
}
