package models

import (
	"github.com/unknwon/com"
	"pokemon/game/utils"
	"strings"
)

var getMPet = DefaultGetMPet
var getMPetByName = DefaultGetMPetByName
var getMSkill = DefaultGetMSkill

type UPet struct {
	ID           int    `gorm:"column:id;primary_key"`
	Bid          int    `gorm:"column:bid;not null;default:-1"`
	Name         string `gorm:"column:name"`
	Uid          int    `gorm:"column:uid;not null"`
	UserName     string `gorm:"column:username"`
	Level        int    `gorm:"column:level;default:1"`
	Wx           int    `gorm:"column:wx;default:1"`
	Ac           int    `gorm:"column:ac;default:0"`
	Mc           int    `gorm:"column:mc;default:0"`
	SrcHp        int    `gorm:"column:srchp;default:0"`
	AddHp        int    `gorm:"column:addhp;default:0"`
	Hp           int    `gorm:"column:hp;default:0"`
	Mp           int    `gorm:"column:mp;default:0"`
	SrcMp        int    `gorm:"column:srcmp;default:0"`
	AddMp        int    `gorm:"column:addmp;default:0"`
	SkillList    string `gorm:"column:skillist"`
	Stime        int    `gorm:"column:stime;default:0"`
	NowExp       int    `gorm:"column:nowexp;default:0"`
	LExp         int    `gorm:"column:lexp;default:0"`
	ImgStand     string `gorm:"column:imgstand"`
	ImgAck       string `gorm:"column:imgack"`
	ImgDie       string `gorm:"column:imgdie"`
	ImgHead      string `gorm:"column:headimg"`
	ImgCard      string `gorm:"column:cardimg"`
	ImgEffect    string `gorm:"column:effectimg"`
	Hits         int    `gorm:"column:hits;default:0"`
	Miss         int    `gorm:"column:miss;default:0"`
	Speed        int    `gorm:"column:speed;default:0"`
	Kx           string `gorm:"column:kx;default:'1'"`
	Subyl        int    `gorm:"column:subyl;default:0"`
	Subsl        int    `gorm:"column:subsl;default:0"`
	Subdl        int    `gorm:"column:subdl;default:0"`
	Subxl        int    `gorm:"column:subxl;default:0"`
	Subhl        int    `gorm:"column:subhl;default:0"`
	Subfl        int    `gorm:"column:subfl;default:0"`
	Subkl        int    `gorm:"column:subkl;default:0"`
	Fatting      int    `gorm:"column:fatting;default:0"`
	RemakeLevel  string `gorm:"column:remakelevel;default:0"`
	RemakeId     string `gorm:"column:remakeid"`
	RemakePid    string `gorm:"column:remakepid"`
	Czl          string `gorm:"column:czl"`
	Psell        int    `gorm:"column:psell;default:0"`
	Pstime       int    `gorm:"column:pstime;default:0"`
	Petime       int    `gorm:"column:petime;default:0"`
	Muchang      int    `gorm:"column:muchang;default:0"`
	Zb           string `gorm:"column:zb"`
	ReMakeTimes  int    `gorm:"column:remaketimes;default:0"`
	TgFlag       int    `gorm:"column:tgflag;default:0"`
	TgStime      int    `gorm:"column:tgstime;default:0"`
	TgTime       int    `gorm:"column:tgtime;default:0"`
	TgMes        int    `gorm:"column:tgmes;default:0"`
	ChChengbb    int    `gorm:"column:chchengbb"`
	ChChengcz    string `gorm:"column:chchengcz"`
	ChChengtime  int    `gorm:"column:chchengtime"`
	ChChengWp    string `gorm:"column:chchengwp"`
	AddSx        string `gorm:"column:addsx"`
	ChChengColor string `gorm:"column:chchengcolor"`
	ChChengSx    string `gorm:"column:chchengsx"`
	OldBid       int    `gorm:"column:old_bid"`
	CqFlag       int    `gorm:"column:cqflag;default=0"`
	MModel       *MPet  `gorm:"-"`

	CC     float64    `gorm:"-"` // czl的float形式
	ZbAttr *PetZbAttr `gorm:"-"`

	//Attributes *Attribute `gorm:"-"`
	//Skills     *[][]int   `gorm:"-"`
}

// TableName sets the insert table name for this struct type
func (up *UPet) TableName() string {
	return "userbb"
}

func (up *UPet) GetM() *MPet {
	if up.MModel == nil {
		if up.Bid < 1 {
			up.MModel = getMPetByName(up.Name)
		} else {
			up.MModel = getMPet(up.Bid)
		}
	}
	return up.MModel
}

func (up *UPet) WxName() string {
	return utils.GetWxName(up.GetM().Wx)

}
func (up *UPet) InMuchang() bool {
	return up.Muchang > 0
}
func (up *UPet) AfterFind() (err error) {
	//up.Name = utils.ToUtf8(up.Name)
	//up.UserName = utils.ToUtf8(up.UserName)
	up.CC = com.StrTo(up.Czl).MustFloat64()
	up.ZbAttr = newPetZbAttr()
	return
}

func (up *UPet) BeforeSave() (err error) {
	//up.Name = utils.ToGbk(up.Name)
	//up.UserName = utils.ToGbk(up.UserName)
	return
}
func DefaultGetMPet(bid int) *MPet {
	return nil
}
func DefaultGetMPetByName(name string) *MPet {
	return nil
}

func SetMPetFunc(f1 func(bid int) *MPet, f2 func(name string) *MPet) {
	getMPet = f1
	getMPetByName = f2
}

type Uskill struct {
	ID    int    `gorm:"column:id;primary_key"`
	Bid   int    `gorm:"column:bid"`
	Sid   int    `gorm:"column:sid"`
	Name  string `gorm:"column:name"`
	Level int    `gorm:"column:level;default:'0'"`
	Vary  string `gorm:"column:vary"`
	Wx    int    `gorm:"column:wx"`
	Value string `gorm:"column:value;default:'0'"`
	Plus  string `gorm:"column:plus;default:'0'"`
	Img   string `gorm:"column:img;default:'0'"`
	Uhp   int    `gorm:"column:uhp"`
	Ump   int    `gorm:"column:ump"`

	EnbleUp     bool    `gorm:"-"`
	MModel      *MSkill `gorm:"-"`
	AckValue    int     `gorm:"-"`
	PlusValue   int     `gorm:"-"`
	EffectValue *struct {
		Key   string
		Value float64
	} `gorm:"-"`
}

func (us *Uskill) GetM() *MSkill {
	if us.MModel == nil {
		us.MModel = getMSkill(us.Sid)
	}
	return us.MModel
}

func (us *Uskill) TableName() string {
	return "skill"
}
func (us *Uskill) AfterFind() (err error) {
	//us.Name = utils.ToUtf8(us.Name)
	us.AckValue = com.StrTo(us.Value).MustInt()
	us.PlusValue = com.StrTo(us.Plus).MustInt()
	if items := strings.Split(us.Img, ":"); len(items) > 1 {
		us.EffectValue = &struct {
			Key   string
			Value float64
		}{Key: items[0], Value: com.StrTo(strings.ReplaceAll(items[1], "%", "")).MustFloat64() * 0.01}

	}
	return
}

func (up *Uskill) BeforeSave() (err error) {
	//up.Name = utils.ToGbk(up.Name)
	return
}

func DefaultGetMSkill(bid int) *MSkill {
	return nil
}

func SetMSkillFunc(f func(bid int) *MSkill) {
	getMSkill = f
}

//func (up *UPet) DecodeAttribute() error {
//	var attributeHM map[string]int
//	err := json.Unmarshal([]byte(up.AttributeHM), &attributeHM)
//	if err != nil {
//		return err
//	}
//
//	up.Attributes.Hp = attributeHM["hp"]
//	up.Attributes.Mp = attributeHM["mp"]
//	fmt.Println(up.Attributes.Hp, attributeHM["hp"])
//	var attributeOt map[string]int
//	err = json.Unmarshal([]byte(up.AttributeOt), &attributeOt)
//	if err != nil {
//		return err
//	}
//	up.Attributes.Ac = attributeOt["ac"]
//	up.Attributes.Mc = attributeOt["mc"]
//	up.Attributes.Hits = attributeOt["hits"]
//	up.Attributes.Miss = attributeOt["miss"]
//	up.Attributes.Speed = attributeOt["speed"]
//	fmt.Println(attributeOt)
//	err = json.Unmarshal([]byte(up.Skill), up.Skills)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (up *UPet) EncodeAttribute() error {
//	attributeHM := make(map[string]int)
//	attributeHM["hp"] = up.Attributes.Hp
//	attributeHM["mp"] = up.Attributes.Mp
//	attributeOt := make(map[string]int)
//	attributeOt["ac"] = up.Attributes.Ac
//	attributeOt["mc"] = up.Attributes.Mc
//	attributeOt["hits"] = up.Attributes.Hits
//	attributeOt["miss"] = up.Attributes.Miss
//	attributeOt["speed"] = up.Attributes.Speed
//	attributeHMJSON, err := json.Marshal(attributeHM)
//	if err != nil {
//		return err
//	}
//	attributeOtJSON, err := json.Marshal(attributeOt)
//	if err != nil {
//		return err
//	}
//	up.AttributeHM = string(attributeHMJSON)
//	up.AttributeOt = string(attributeOtJSON)
//	return nil
//}

type Growth struct {
	ID    int `gorm:"column:id;primary_key"`
	Wx    int `gorm:"column:wx"`
	Hp    int `gorm:"column:hp"`
	Mp    int `gorm:"column:mp"`
	Ac    int `gorm:"column:ac"`
	Mc    int `gorm:"column:mc"`
	Hits  int `gorm:"column:hits"`
	Miss  int `gorm:"column:miss"`
	Speed int `gorm:"column:speed"`
	J     int `gorm:"column:j"`
	M     int `gorm:"column:m"`
	S     int `gorm:"column:s"`
	H     int `gorm:"column:h"`
	T     int `gorm:"column:t"`
}

// TableName sets the insert table name for this struct type
func (up *Growth) TableName() string {
	return "wx"
}

type SSjhRule struct {
	ID         int    `gorm:"column:id;primary_key"`
	PetId      int    `gorm:"column:pet_id"`
	NeedLevels string `gorm:"column:need_levels"`
	NeedProps  string `gorm:"column:need_props"`
	MaxCzl     int    `gorm:"column:max_czl"`
	ZsProgress int    `gorm:"column:zs_progress"`
	ZsLine     int    `gorm:"column:zs_line"`
	MaxLevel   int    `gorm:"column:max_level"`
}

// TableName sets the insert table name for this struct type
func (up *SSjhRule) TableName() string {
	return "super_jh"
}

type SSzsRule struct {
	ID              int    `gorm:"column:id;primary_key"`
	PetId           int    `gorm:"column:cur_pet_id"`
	NeedLevel       int    `gorm:"column:need_level"`
	NeedCzl         int    `gorm:"column:need_czl"`
	NeedProps       string `gorm:"column:need_props"`
	BaseSuccessRate int    `gorm:"column:base_success_rate"`
	FailedCzlPecent int    `gorm:"column:failed_czl_percent"`
	NextPetId       int    `gorm:"column:next_pet_id"`
}

// TableName sets the insert table name for this struct type
func (up *SSzsRule) TableName() string {
	return "super_zs"
}

type PetZbAttr struct {
	Hp    int
	Mp    int
	Ac    int
	Mc    int
	Hits  int
	Miss  int
	Speed int

	Special map[string]float64
}

func newPetZbAttr() *PetZbAttr {
	return &PetZbAttr{Special: map[string]float64{"crit": 0.05}}
}
