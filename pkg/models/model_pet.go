package models

import (
	"pokemon/pkg/utils"
)

type MPet struct {
	ID          int    `gorm:"column:id;primary_key"`
	Name        string `gorm:"column:name;not null"`
	Wx          int    `gorm:"column:wx;default=1"`
	Ac          int    `gorm:"column:ac;default=0"`
	Mc          int    `gorm:"column:mc;default=0"`
	Hp          int    `gorm:"column:hp;default=0"`
	Mp          int    `gorm:"column:mp;ndefault=0"`
	Speed       int    `gorm:"column:speed;default=0"`
	Hits        int    `gorm:"column:hits;default=0"`
	Miss        int    `gorm:"column:miss;default=0"`
	ImgStand    string `gorm:"column:imgstand"`
	ImgAck      string `gorm:"column:imgack"`
	ImgDie      string `gorm:"column:imgdie"`
	SkillList   string `gorm:"column:skillist"`
	Czl         string `gorm:"column:czl"`
	Kx          string `gorm:"column:kx"`
	ReMakeLevel string `gorm:"column:remakelevel;default='0'"`
	ReMakeId    string `gorm:"column:remakeid"`
	ReMakePid   string `gorm:"column:remakepid;default='0'"`
	NowExp      int    `gorm:"column:wx;default=0"`
	LExp        int    `gorm:"column:wx;default=0"`
	Subyl       int    `gorm:"column:subyl;default=0"`
	Subsl       int    `gorm:"column:subsl;default=0"`
	Subxl       int    `gorm:"column:subxl;default=0"`
	Subdl       int    `gorm:"column:subdl;default=0"`
	Subfl       int    `gorm:"column:subfl;default=0"`
	Subhl       int    `gorm:"column:subhl;default=0"`
	Subkl       int    `gorm:"column:subkl;default=0"`
	ImgHead     string `gorm:"column:headimg"`
	ImgCard     string `gorm:"column:cardimg"`
	ImgEffect   string `gorm:"column:effectimg"`
	Bbdesc      string `gorm:"column:bbdesc"`
	WxName      string `gorm:"-"`

	//Attribute   string `gorm:"column:attribute;not null"` // ac,mc,hp,mp,speed,hits,miss,cz,kx，使用json保存
	//Class       int   `gorm:"column:class;not null"`
	//ImageId     int    `gorm:"column:image_id;not null"`
	//Skill       string `gorm:"column:skill;not null"`     //1,3,4,7  默认初始均为1
	//Evolution   string `gorm:"column:evolution;not null"` //1,3,4,7  默认初始均为1
	//Description string `gorm:"column:description;not null"`
}

// TableName sets the insert table name for this struct type
func (mpet *MPet) TableName() string {
	return "bb"
}

func (mpet *MPet) AfterFind() (err error) {
	mpet.WxName = utils.GetWxName(mpet.Wx)
	//mpet.Name = utils.ToUtf8(mpet.Name)
	return
}

type MSkill struct {
	ID       int    `gorm:"column:id;primary_key"`
	Pid      int    `gorm:"column:pid;default:0"`
	Name     string `gorm:"column:name"`
	Vary     string `gorm:"column:vary"`
	Wx       int    `gorm:"column:wx;default:0"`
	Img      string `gorm:"column:img"`
	AckValue string `gorm:"column:ackvalue"`
	Plus     string `gorm:"column:plus"`
	Requires string `gorm:"column:requires"`
	Uhp      string `gorm:"column:uhp"`
	Ump      string `gorm:"column:ump"`
	AckStyle int    `gorm:"column:ackstyle;default:1"`
	ImgEft   string `gorm:"column:imgeft"`
}

// TableName sets the insert table name for this struct type
func (mskill *MSkill) TableName() string {
	return "skillsys"
}
func (mskill *MSkill) AfterFind() (err error) {
	//mskill.Name = utils.ToUtf8(mskill.Name)
	return
}

type ExpList struct {
	ID        int `gorm:"column:id;primary_key"`
	Level     int `gorm:"column:level"`
	NextLvExp int `gorm:"column:nxtlvexp"`
}

func (expList *ExpList) TableName() string {
	return "exptolv"
}

type Attribute struct {
	Hp    int `gorm:"-"`
	Mp    int `gorm:"-"`
	Ac    int `gorm:"-"`
	Mc    int `gorm:"-"`
	Hits  int `gorm:"-"`
	Miss  int `gorm:"-"`
	Speed int `gorm:"-"`
}

type ZbAttribute struct {
	Hp                 int     `gorm:"-"`
	Mp                 int     `gorm:"-"`
	Ac                 int     `gorm:"-"`
	Mc                 int     `gorm:"-"`
	Hits               int     `gorm:"-"`
	Miss               int     `gorm:"-"`
	Speed              int     `gorm:"-"`
	HpRate             float64 `gorm:"-"`
	MpRate             float64 `gorm:"-"`
	AcRate             float64 `gorm:"-"`
	McRate             float64 `gorm:"-"`
	HitsRate           float64 `gorm:"-"`
	MissRate           float64 `gorm:"-"`
	CriticalRate       float64 `gorm:"-"` // 暴击率
	CriticalEffectRate float64 `gorm:"-"` // 额外暴击伤害比率
	AbsorbRate         float64 `gorm:"-"` // 吸血比率
	ExtraRate          float64 `gorm:"-"` // 伤害加深比率
	ReduceRate         float64 `gorm:"-"` // 伤害抵消比率
}

type MergeRule struct {
	Id     int
	Aid    int
	Bid    int
	Maid   int
	Mbid   int
	Limits string
}

func (this *MergeRule) TableName() string {
	return "merge"
}

type ZsRule struct {
	Id  int
	Aid int
	Bid int
	Mid int
}

func (this *ZsRule) TableName() string {
	return "zs"
}
