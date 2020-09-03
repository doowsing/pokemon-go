package models

import (
	"fmt"
	"time"
)

/*
日志结构体：
	游戏普通日志
	游戏任务日志
	游戏登录日志
	游戏系统日志
	游戏充值日志
	游戏元宝消费日志
*/
type GameLog struct {
	ID       uint   `gorm:"column:id;primary_key" json:"-"`     // 自增
	SUid     int    `gorm:"column:seller;not null" json:"suid"` // 主操作人
	BUid     int    `gorm:"column:buyer" json:"buid"`           // 副操作人，如交易物品时，主操作人负责卖出，副操作人负责买入
	Note     string `gorm:"column:pnote" json:"note"`
	Category int    `gorm:"column:vary;not null" json:"category"`
	Time     int    `gorm:"column:ptime;not null" json:"time"`
}

func (p *GameLog) TableName() string {
	return "gamelog"
}

func (p *GameLog) BeforeSave() (err error) {
	//p.Note = utils.ToGbk(p.Note)
	return
}

func (p *GameLog) String() string {
	return fmt.Sprintf("gamelog: suid=>%d, buid=>%d, category=>%d, time=>%d, note=>%s", p.SUid, p.BUid, p.Category, p.Time, p.Note)
}

type LoginLog struct {
	ID    uint   `gorm:"column:Id;primary_key"` // 自增
	UName string `gorm:"column:uname;not null"`
	IP    string `gorm:"column:uIP;not null"`
	Time  int    `gorm:"column:times;not null"`
}

func (p *LoginLog) TableName() string {
	return "logins"
}

func (p *LoginLog) String() string {
	return fmt.Sprintf("loginlog: uname=>%s, ip=>%s, time=>%d", p.UName, p.IP, p.Time)
}

func (p *LoginLog) BeforeSave() (err error) {
	//p.UName = utils.ToGbk(p.UName)
	return
}

type SysLog struct {
	ID   uint      `gorm:"primary_key"` // 自增
	Uid  int       `gorm:"column:uid;default=0"`
	Note string    `gorm:"column:note;type:text"`
	Time time.Time `gorm:"column:time;index:time"`
}

func (p *SysLog) TableName() string {
	return "sys_log"
}

type ChargeLog struct {
	ID    uint      `gorm:"primary_key"` // 自增
	Uid   int       `gorm:"column:uid;default=0;not null"`
	PayId string    `gorm:"column:pay_id;not null"`
	Money int       `gorm:"column:note;not null"`
	Yb    int       `gorm:"column:yb;not null"`
	Time  time.Time `gorm:"column:time;index:time;not null"`
}

func (p *ChargeLog) TableName() string {
	return "charge_log"
}

type ConspLog struct {
	ID     uint      `gorm:"primary_key"`                   // 自增
	Uid    int       `gorm:"column:uid;not null;index:uid"` //
	Pid    int       `gorm:"column:pid;not null"`
	Number int       `gorm:"column:number;not null"`
	UsedYb int       `gorm:"column:used_yb;not null"`
	LeftYb int       `gorm:"column:left_yb;not null"`
	Time   time.Time `gorm:"column:time;index:time;not null"`
}

func (p *ConspLog) TableName() string {
	return "consp_log"
}

type YbLog struct {
	ID      int    `gorm:"column:id;primary_key"`
	Pid     int    `gorm:"column:title"`
	Account string `gorm:"column:nickname"`
	UseYb   int    `gorm:"column:yb"`
	Btime   int    `gorm:"column:buytime"`
	Pnote   string `gorm:"column:pname"`
	Num     int    `gorm:"column:nums"`
}

func (p *YbLog) TableName() string {
	return "yblog"
}
