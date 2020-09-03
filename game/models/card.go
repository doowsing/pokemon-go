package models

import (
	"fmt"
	"github.com/unknwon/com"
	"strings"
)

type CardRecord struct {
	ID      int `gorm:"column:id;primary_key"`
	Uid     int `gorm:"column:uid"`
	CardPid int `gorm:"column:card_pid"`
	Sum     int `gorm:"column:sum"`
}

func (u *CardRecord) TableName() string {
	return "t_card_user"
}

type CardTitle struct {
	ID        int    `gorm:"column:id;primary_key"`
	CodeName  string `gorm:"column:F_title_name"`
	Name      string `gorm:"column:F_title_Chinese"`
	Img       string `gorm:"column:F_title_img"`
	NeedCard  string `gorm:"column:F_title_must_card"`   //获得所需目标
	NeedDep   string `gorm:"column:F_title_get_methods"` //目标描述
	Hp        int    `gorm:"column:F_add_hp;default:0"`
	Mp        int    `gorm:"column:F_add_mc;default:0"`
	Ac        int    `gorm:"column:F_add_ac;default:0"`
	Mc        int    `gorm:"column:F_add_mc;default:0"`
	Hits      int    `gorm:"column:F_add_hits;default:0"`
	Miss      int    `gorm:"column:F_add_miss;default:0"`
	Speed     int    `gorm:"column:F_add_speed;default:0"`
	HpRate    int    `gorm:"column:F_add_hprate;default:0"`
	MpRate    int    `gorm:"column:F_add_mprate;default:0"`
	AcRate    int    `gorm:"column:F_add_acrate;default:0"`
	McRate    int    `gorm:"column:F_add_mcrate;default:0"`
	HitsRate  int    `gorm:"column:F_add_hitsrate;default:0"`
	MissRate  int    `gorm:"column:F_add_missrate;default:0"`
	SpeedRate int    `gorm:"column:F_add_speedrate;default:0"`
	Dxsh      int    `gorm:"column:F_dxsh;default:0"`
	HitsHp    int    `gorm:"column:F_hitshp;default:0"`
	HitsMp    int    `gorm:"column:F_hitsmp;default:0"`
	Shjs      int    `gorm:"column:F_shjs;default:0"`
	Sdmp      int    `gorm:"column:F_sdmp;default:0"`
	Szmp      int    `gorm:"column:F_szmp;default:0"`
	AddMoney  int    `gorm:"column:F_addmoney;default:0"`
	Time      int    `gorm:"column:F_time;default:0"`

	EffectDep string `gorm:"-"`
}

func (u *CardTitle) TableName() string {
	return "t_card_to_title"
}

func (c *CardTitle) AfterFind() (err error) {
	//c.Name = utils.ToUtf8(c.Name)
	//c.NeedCard = utils.ToUtf8(c.NeedCard)
	//c.NeedDep = utils.ToUtf8(c.NeedDep)
	var effect []string
	if c.Hp > 0 {
		effect = append(effect, fmt.Sprintf("增加HP:%d", c.Hp))
	}
	if c.Mp > 0 {
		effect = append(effect, fmt.Sprintf("增加MP:%d", c.Mp))
	}
	if c.Ac > 0 {
		effect = append(effect, fmt.Sprintf("增加攻击:%d", c.Ac))
	}
	if c.Mc > 0 {
		effect = append(effect, fmt.Sprintf("增加防御:%d", c.Mc))
	}
	if c.Hits > 0 {
		effect = append(effect, fmt.Sprintf("增加命中:%d", c.Hits))
	}
	if c.Miss > 0 {
		effect = append(effect, fmt.Sprintf("增加闪避:%d", c.Miss))
	}
	if c.Speed > 0 {
		effect = append(effect, fmt.Sprintf("增加速度:%d", c.Speed))
	}
	if c.HpRate > 0 {
		effect = append(effect, fmt.Sprintf("增加HP百分比:%d%%", c.HpRate))
	}
	if c.MpRate > 0 {
		effect = append(effect, fmt.Sprintf("增加MP百分比:%d%%", c.MpRate))
	}
	if c.AcRate > 0 {
		effect = append(effect, fmt.Sprintf("增加攻击百分比:%d%%", c.AcRate))
	}
	if c.McRate > 0 {
		effect = append(effect, fmt.Sprintf("增加防御百分比:%d%%", c.McRate))
	}
	if c.HitsRate > 0 {
		effect = append(effect, fmt.Sprintf("增加命中百分比:%d%%", c.HitsRate))
	}
	if c.MissRate > 0 {
		effect = append(effect, fmt.Sprintf("增加闪避百分比:%d%%", c.MissRate))
	}
	if c.SpeedRate > 0 {
		effect = append(effect, fmt.Sprintf("增加速度百分比:%d%%", c.SpeedRate))
	}
	if c.Dxsh > 0 {
		effect = append(effect, fmt.Sprintf("抵消伤害:%d%%", c.Dxsh))
	}
	if c.HitsHp > 0 {
		effect = append(effect, fmt.Sprintf("吸取伤害的%d%%转化为自身HP", c.HitsHp))
	}
	if c.HitsMp > 0 {
		effect = append(effect, fmt.Sprintf("吸取伤害的%d%%转化为自身MP", c.HitsMp))
	}
	if c.Shjs > 0 {
		effect = append(effect, fmt.Sprintf("对敌人造成的伤害增加:%d%%", c.Shjs))
	}
	if c.Sdmp > 0 {
		effect = append(effect, fmt.Sprintf("将受到伤害的%d%%已MP抵消", c.Sdmp))
	}
	if c.Szmp > 0 {
		effect = append(effect, fmt.Sprintf("将受到伤害的%d%%转化为MP", c.Szmp))
	}
	if c.AddMoney > 0 {
		effect = append(effect, fmt.Sprintf("战斗胜利获得金币增加%d", c.AddMoney))
	}
	if c.Time > 0 {
		effect = append(effect, fmt.Sprintf("自动攻击间隔时间减少%d", c.Time))
	}
	c.EffectDep = strings.Join(effect, " ")
	return
}

type CardSeries struct {
	ID   int    `gorm:"column:id;primary_key"`
	Name string `gorm:"column:F_Class_Name"`
	Card string `gorm:"column:F_Had_Card"`

	CardList []string `gorm:"-"`
}

func (c *CardSeries) TableName() string {
	return "t_card_type"
}

func (c *CardSeries) AfterFind() (err error) {
	//c.Name = utils.ToUtf8(c.Name)
	//c.NeedCard = utils.ToUtf8(c.NeedCard)
	//c.NeedDep = utils.ToUtf8(c.NeedDep)
	c.CardList = strings.Split(c.Card, ",")
	return
}

type CardPrize struct {
	ID    int    `gorm:"column:id;primary_key"`
	Need  string `gorm:"column:F_Satisfy_condition"`
	Prize string `gorm:"column:F_Prize"`
	Title string `gorm:"column:F_Prize_title"`

	NeedList []*struct {
		Name string
		Num  int
	} `gorm:"-"`

	PrizeList []*struct {
		Pid int
		Num int
	} `gorm:"-"`
}

func (c *CardPrize) TableName() string {
	return "t_card_prize"
}

func (c *CardPrize) AfterFind() (err error) {
	for _, s := range strings.Split(c.Need, ",") {
		items := strings.Split(s, ":")
		if len(items) > 1 {
			c.NeedList = append(c.NeedList, &struct {
				Name string
				Num  int
			}{Name: items[0], Num: com.StrTo(items[1]).MustInt()})
		}
	}

	for _, s := range strings.Split(c.Prize, ",") {
		items := strings.Split(s, ":")
		if len(items) > 1 {
			c.PrizeList = append(c.PrizeList, &struct {
				Pid int
				Num int
			}{Pid: com.StrTo(items[0]).MustInt(), Num: com.StrTo(items[1]).MustInt()})
		}
	}
	return
}
