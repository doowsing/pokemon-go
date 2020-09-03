package models

import "time"

type RecordFb struct {
	ID       int `gorm:"column:Id;primary_key"`
	Uid      int `gorm:"column:uid"`
	GpcId    int `gorm:"column:gwid"`
	LeftTime int `gorm:"column:lttime"`
	InMap    int `gorm:"column:inmap;default:0"`
	SrcTime  int `gorm:"column:srctime;default:0"` // 刷新周期，单位：秒
}

func (m *RecordFb) TableName() string {
	return "fuben"
}

type RecordBoss struct {
	ID       int       `gorm:"column:id;primary_key"`
	MapId    int       `gorm:"column:map_id"`
	KillTime time.Time `gorm:"column:kill_time"`
	MeetTime time.Time `gorm:"column:meet_time"`
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

type TTRecord struct {
	GpcGroupId int
	Index      int
}

type TarotCard struct {
	Id      int    `gorm:"column:id;primary_key"`
	Name    string `gorm:"column:name"`
	Sj      int    `gorm:"column:sj"`
	Effect  string `gorm:"column:effect"`
	Flag    int    `gorm:"column:flag"`
	Content string `gorm:"column:content"`
	Boss    string `gorm:"column:boss"`
	Img     string `gorm:"column:img"`
	MapId   int    `gorm:"column:mapid"`

	Multiple int `gorm:"-"`
}

func (m *TarotCard) AfterFind() (err error) {
	m.Multiple = m.MapId - 128
	if m.Multiple < 0 {
		m.Multiple = 0
	}
	return
}

func (m *TarotCard) TableName() string {
	return "tarot"
}

type SSBattleUser struct {
	Id           int    `gorm:"column:id;primary_key"`
	Uid          int    `gorm:"column:uid"`
	Pos          int    `gorm:"column:pos"`
	Bid          int    `gorm:"column:bid"`
	JgValue      int    `gorm:"column:jgvalue"`
	Levels       string `gorm:"column:levels"`
	AddJgValue   int    `gorm:"column:addjgvalue"`
	AckValue     int    `gorm:"column:ackvalue"`
	FailJgValue  int    `gorm:"column:failjgvalue"`
	FailAckValue int    `gorm:"column:failackvalue"`
	LastVTime    int    `gorm:"column:lastvtime"`
	DoubleJg     int    `gorm:"column:doublejg"`
	Tops         int    `gorm:"column:tops"`
	CurJgValue   int    `gorm:"column:curjgvalue"`
	BoxNum       int    `gorm:"column:boxnum"`
	Nscf         int    `gorm:"column:nscf"`
	SubHp        int    `gorm:"column:subhp"`
	AddHp        int    `gorm:"column:addhp"`
}

func (m *SSBattleUser) TableName() string {
	return "battlefield_user"
}

type SSBattle struct {
	Id           int    `gorm:"column:id;primary_key"`
	PosName      string `gorm:"column:posname"`
	SrcHp        int    `gorm:"column:srchp"`
	Hp           int    `gorm:"column:hp"`
	MaxUser      int    `gorm:"column:maxuser"`
	LevelGet     string `gorm:"column:level_get"`
	BfDate       int    `gorm:"column:bfdate"`
	StartTime    int    `gorm:"column:start_time"`
	EndTime      int    `gorm:"column:end_time"`
	TipsTime     int    `gorm:"column:tips_time"`
	BfMlNum      int    `gorm:"column:bf_ml_num"`
	BfLevelLimit int    `gorm:"column:bf_level_limit"`
	CountF       int    `gorm:"column:countf"`
	StartF       int    `gorm:"column:startf"`
	Success      int    `gorm:"column:success"`
	Ends         int    `gorm:"column:ends"`
}

func (m *SSBattle) TableName() string {
	return "battlefield"
}

type SSBattleGood struct {
	Id    int `gorm:"column:id;primary_key"`
	Pid   int `gorm:"column:pid"`
	Price int `gorm:"column:need"`
}

func (m *SSBattleGood) TableName() string {
	return "battlefield_props"
}

type SSBattleLog struct {
	Id    int    `gorm:"column:id;primary_key"`
	Uid   int    `gorm:"column:uid"`
	UseJg int    `gorm:"column:usejg"`
	Type  string `gorm:"column:type"`
	Num   string `gorm:"column:num"`
	Pid   string `gorm:"column:pid"`
	Times int    `gorm:"column:times"`
}

func (m *SSBattleLog) TableName() string {
	return "jg_log"
}

type BossRecord struct {
	Id       int `gorm:"column:id;primary_key"`
	GpcId    int `gorm:"column:gid;default:0"`
	Rtime    int `gorm:"column:rtime;default:0"`
	FightUid int `gorm:"column:fightuid;default:0"`
	Dtime    int `gorm:"column:dtime;default:0"`
	Glock    int `gorm:"column:glock;default:0"`
}

func (u *BossRecord) TableName() string {
	return "boss_refresh"
}
