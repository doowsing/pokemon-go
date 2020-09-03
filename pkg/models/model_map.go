package models

import (
	"github.com/unknwon/com"
	"pokemon/pkg/utils"
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

	Gpcs     []*Gpc                 `gorm:"-"`
	GpcNames []string               `gorm:"-"`
	Drops    [][]int                `gorm:"-"`
	Needs    map[string]interface{} `gorm:"-"`
	Levels   []int                  `gorm:"-"`
}

func (m *Map) TableName() string {
	return "map"
}

func (m *Map) AfterFind() (err error) {
	m.Name = utils.ToUtf8(m.Name)
	m.Description = utils.ToUtf8(m.Description)
	m.GpcList = utils.ToUtf8(m.GpcList)
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
	Wx        int8   `gorm:"column:wx;default:1"`
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
	Skills     *[][]int   `gorm:"-"`
	Drops      *[][]int   `gorm:"-"`
}

func (m *Gpc) TableName() string {
	return "gpc"
}

func (m *Gpc) AfterFind() (err error) {
	m.Name = utils.ToUtf8(m.Name)
	return
}

type GpcGroup struct {
	ID       int  `gorm:"column:id;primary_key"`
	Category int8 `gorm:"column:category;not null"` // TT层数，0则为遗忘副本

	SpeInfo  string `gorm:"column:spe_info"` //  遗忘地图ID，层数，批次
	GpcList  string `gorm:"column:gpc_list"`
	DropList string `gorm:"column:drop_list"`

	// 反序列化后的数据
	Gpcs  *[]int   `gorm:"-"`
	Drops *[][]int `gorm:"-"`
}

func (m *GpcGroup) TableName() string {
	return "model_gpc_group"
}
