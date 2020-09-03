package models

import "time"

type RecordFb struct {
	ID       int `gorm:"Id"`
	Uid      int `gorm:"uid"`
	GpcId    int `gorm:"column:gwid"`
	LeftTime int `gorm:"column:lttime"`
	InMap    int `gorm:"column:inmap;default:0"`
	SrcTime  int `gorm:"column:srctime;default:0"` // 刷新周期，单位：秒
}

func (m *RecordFb) TableName() string {
	return "fuben"
}

type RecordBoss struct {
	ID       int       `gorm:"id"`
	MapId    int       `gorm:"map_id"`
	KillTime time.Time `gorm:"kill_time"`
	MeetTime time.Time `gorm:"meet_time"`
}

func (m *RecordBoss) TableName() string {
	return "record_boss"
}

type FightStatus struct {
	// 保存战斗状态使用的类
	UUID     int
	LastTime time.Time

	// 宠物状态，考虑是否保存装备的属性
	// 怪物状态
	Gpc *Gpc
}
