package models

import (
	"github.com/unknwon/com"
	"strings"
)

type Map struct {
	ID          int    `gorm:"column:id;primary_key"`
	Name        string `gorm:"column:name;not null"`
	Description string `gorm:"column:descs"`
	GpcList     string `gorm:"column:gpclist"`
	Level       string `gorm:"column:level"`
	//GpcList     string `gorm:"column:gpc_list"`
	DropList      string `gorm:"column:drop_list"`
	Img           string `gorm:"column:img"`
	Need          string `gorm:"column:needs;default:'0'"`
	MultiMonsters string `gorm:"column:multi_monsters"`
	CzlProp       string `gorm:"column:czlprops"`

	GpcIds    []int                  `gorm:"-"`
	Lv2GpcIds map[int][]int          `gorm:"-"`
	GpcNames  []string               `gorm:"-"`
	Drops     [][]int                `gorm:"-"`
	Needs     map[string]interface{} `gorm:"-"`
	Levels    []int                  `gorm:"-"`
}

func (m *Map) TableName() string {
	return "map"
}

func (m *Map) AfterFind() (err error) {
	//m.Name = utils.ToUtf8(m.Name)
	//m.Description = utils.ToUtf8(m.Description)
	//m.GpcList = utils.ToUtf8(m.GpcList)
	m.GpcNames = strings.Split(m.GpcList, ",")
	for _, l := range strings.Split(m.Level, ",") {
		m.Levels = append(m.Levels, com.StrTo(l).MustInt())
	}
	m.Level = strings.Join(strings.Split(m.Level, ","), "-")
	return
}

type Gpc struct {
	ID        int    `gorm:"column:id;primary_key"`
	Name      string `gorm:"column:name;not null"`
	Wx        int    `gorm:"column:wx;default:1"`
	Level     int    `gorm:"column:level;default:1"`
	Hp        int    `gorm:"column:hp;default:1"`
	Mp        int    `gorm:"column:mp;default:0"`
	Ac        int    `gorm:"column:ac;default:0"`
	Mc        int    `gorm:"column:mc;default:0"`
	Hits      int    `gorm:"column:hits;default:0"`
	Speed     int    `gorm:"column:speed;default:0"`
	Miss      int    `gorm:"column:miss;default:0"`
	CatchRate int    `gorm:"column:catchv"`
	CatchBid  int    `gorm:"column:catchid"`
	Skill     string `gorm:"column:skill"`

	ImgStand string `gorm:"column:imgstand"`
	ImgAck   string `gorm:"column:imgack"`
	ImgDie   string `gorm:"column:imgdie"`
	DropList string `gorm:"column:droplist"`
	Exp      int    `gorm:"column:exps"`
	Money    int    `gorm:"column:money;default:1"`

	Boss           int    `gorm:"column:boss;default:0"`
	Kx             string `gorm:"column:kx"`
	ActiveDropList string `gorm:"column:activedroplist"`

	Attributes *Attribute `gorm:"-"`
	Skills     []*struct {
		Sid   int
		Level int
	} `gorm:"-"`
	Drops []*RatePid `gorm:"-"`
}

func (m *Gpc) TableName() string {
	return "gpc"
}

func (m *Gpc) AfterFind() (err error) {
	//m.Name = utils.ToUtf8(m.Name)
	var sid, lv, pid, rate int

	m.Skills = []*struct {
		Sid   int
		Level int
	}{}
	for _, s := range strings.Split(m.Skill, ",") {
		item := strings.Split(s, ":")
		if len(item) > 1 {
			sid = com.StrTo(item[0]).MustInt()
			lv = com.StrTo(item[1]).MustInt()
			if sid > 0 && lv > 0 {
				m.Skills = append(m.Skills, &struct {
					Sid   int
					Level int
				}{Sid: sid, Level: lv})
			}
		}
	}

	m.Drops = []*RatePid{}
	for _, s := range strings.Split(m.DropList, ",") {
		item := strings.Split(s, ":")
		if len(item) > 1 {
			pid = com.StrTo(item[0]).MustInt()
			rate = com.StrTo(item[1]).MustInt()
			if pid > 0 && rate > 0 {
				m.Drops = append(m.Drops, &RatePid{Pid: pid, Rate: rate})
			}
		}
	}
	return
}

type GpcGroup struct {
	ID      int    `gorm:"column:id;primary_key"`
	Gpcs    string `gorm:"column:gpc"`
	Level   int    `gorm:"column:boss"`
	Drop    string `gorm:"column:drops"`
	MapId   int    `gorm:"column:map_id"`
	StepId  int    `gorm:"column:step_id"`
	GroupId int    `gorm:"column:group_id"`

	GpcList  []int      `gorm:"-"`
	DropList []*RatePid `gorm:"-"`
}

func (m *GpcGroup) TableName() string {
	return "c_gpc"
}

func (m *GpcGroup) AfterFind() (err error) {
	list := []int{}
	for _, id := range strings.Split(m.Gpcs, ",") {
		list = append(list, com.StrTo(id).MustInt())
	}
	m.GpcList = list
	DropList := []*RatePid{}
	for _, s := range strings.Split(m.Drop, ",") {
		item := strings.Split(s, ":")
		if len(item) > 1 {
			DropList = append(DropList, &RatePid{Pid: com.StrTo(item[0]).MustInt(), Rate: com.StrTo(item[1]).MustInt()})
		}
	}
	m.DropList = DropList
	return
}

type RatePid struct {
	Pid  int
	Rate int
}
