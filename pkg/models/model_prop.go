package models

import "pokemon/pkg/utils"

type MProp struct {
	ID           int    `gorm:"column:id;primary_key"`
	Name         string `gorm:"column:name"`
	Requires     string `gorm:"column:requires"`
	Usages       string `gorm:"column:usages"`
	Effect       string `gorm:"column:effect"`
	SellJb       int    `gorm:"column:sell;default=0"`
	Prestige     int    `gorm:"column:prestige;default=0"`
	BuyJb        int    `gorm:"column:buy;default=0"`
	BuyYb        int    `gorm:"column:yb;default=0"`
	BuySj        int    `gorm:"column:sj;default=0"`
	Stime        int    `gorm:"column:stime;default=0"`
	EndTime      int    `gorm:"column:endtime;default=0"`
	Img          string `gorm:"column:img"`
	Vary         int    `gorm:"column:vary;default=1"`     //是否可叠加，1为可叠加，2为不可叠加
	VaryName     int    `gorm:"column:varyname;default=0"` // 道具分类
	Position     int    `gorm:"column:postion;default=0"`  // 道具分类
	PlusEffect   string `gorm:"column:pluseffect"`
	PlusFlag     int    `gorm:"column:plusflag;default=0"`
	PlusPid      int    `gorm:"column:pluspid;default=0"`
	PlusGet      string `gorm:"column:plusget"`
	PlusNum      int    `gorm:"column:plusnum;default=0"`
	PropsColor   int    `gorm:"column:propscolor"`
	PropsLock    int    `gorm:"column:propslock;default=0"` // 是否可交易，0不可以，1可以
	Series       string `gorm:"column:series"`
	SeriesEffect string `gorm:"column:serieseffect"`
	Expire       int    `gorm:"column:expire"`
	Note         string `gorm:"column:note"`
	TimeLimit    string `gorm:"column:timelimit"`
	Merge        int    `gorm:"column:merge;default=0"`
	Vip          int    `gorm:"column:vip"`
	Honor        int    `gorm:"column:honor"`
	Contribution int    `gorm:"column:contribution"`
	GuildLevel   int    `gorm:"column:guild_level"`
	ZhekouYb     int    `gorm:"column:zhekouyb;default=0"`

	//ZbInfos     *MZbInfo `gorm:"-"`
	VaryNameStr string `gorm:"-"` // 道具类别的名称
	VaryStr     string `gorm:"-"` // 是否可叠加
}

func (mprop *MProp) TableName() string {
	return "props"
}

func (mprop *MProp) IsVary() bool {
	// 是否可叠加
	return mprop.Vary < 2
}

func (mprop *MProp) EnableDeal() bool {
	// 是否可交易
	return mprop.PropsLock > 0
}

func (mprop *MProp) AfterFind() (err error) {
	mprop.VaryNameStr = utils.GetVaryNameStr(mprop.VaryName)
	//mprop.Name = utils.ToUtf8(mprop.Name)
	//mprop.Usages = utils.ToUtf8(mprop.Usages)
	//mprop.Series = utils.ToUtf8(mprop.Series)
	return
}

type MSeries struct {
	ID     int    `gorm:"column:id;primary_key"`
	Name   string `gorm:"column:name;not null"`
	Zb     string `gorm:"column:zb"`
	Effect string `gorm:"column:effect"`
}

func (mprop *MSeries) TableName() string {
	return "model_series"
}

type MZbInfo struct {
	MainInfo         map[string]float64
	OtherInfo        map[string]float64
	EnableStrengthen bool
	StrengthenPid    int
	StrengthenEffect []float64
	Position         int
	Series           int
}
