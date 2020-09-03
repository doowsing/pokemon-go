package models

var getMProp = DefaultGetMProp

type UProp struct {
	ID         int    `gorm:"column:id;primary_key"`
	Pid        int    `gorm:"column:pid;not null"`
	Uid        int    `gorm:"column:uid;not null"`
	Sell       int    `gorm:"column:sell;default:0"`
	Vary       int    `gorm:"column:vary;default:1"`
	Sums       int    `gorm:"column:sums;default:0"`
	Stime      int    `gorm:"column:stime;default:0"`
	Psell      int    `gorm:"column:psell;default:0"`
	Pstime     int    `gorm:"column:pstime;default:0"`
	Petime     int    `gorm:"column:petime;default:0"`
	Psum       int    `gorm:"column:psum;default:0"`
	Psj        int    `gorm:"column:psj;default:0"`
	Bsum       int    `gorm:"column:bsum;default:0"`
	Zbing      int    `gorm:"column:zbing;default:0"`
	Zbpets     int    `gorm:"column:zbpets;default:0"`
	BuyCode    int    `gorm:"column:buycode;default:0"`
	PlusTmsEft string `gorm:"column:plus_tms_eft"`
	CanTrade   int    `gorm:"column:cantrade;default:0"`
	FHoleInfo  string `gorm:"column:F_item_hole_info"`
	Pyb        int    `gorm:"column:pyb;default:0"`

	//EnableDeal int       `gorm:"column:enable_deal;default:0"` // 是否可交易
	//BPSum      int       `gorm:"column:bp_sum;default:0"`      // 背包数量
	//CKSum      int       `gorm:"column:ck_sum;default:0"`      // 仓库数量
	//SellInfo   string    `gorm:"column:sell_info"`             // 拍卖所信息 {"jb":{"price":10,"expire":"2019 06 17", "yb":...}
	//ZbPet      int       `gorm:"column:zb_pet"`                // 装备宠物ID
	//ZbInfo     string    `gorm:"column:zb_info"`               // 装备信息: {"strengthen":{"level":1,"effect":{"acrate":0.5}},"Inser":[{"acrate":0.5},{"acrate":0.25}]}
	//Expire     time.Time `gorm:"column:expire"`

	MModel     *MProp `gorm:"-"`
	PmTimeStr  string `gorm:"-"` // 拍卖时间字符串
	PmMoneyStr string `gorm:"-"` // 拍卖货币字符串
}

func (up *UProp) TableName() string {
	return "userbag"
}
func (up *UProp) GetM() *MProp {
	if up.MModel == nil {
		up.MModel = getMProp(up.Pid)
	}
	return up.MModel
}
func (up *UProp) AfterFind() (err error) {
	return
}

func (up *UProp) AbleTrade() bool {
	if up.ID == 0 {
		return false
	}
	return (up.CanTrade == 0 && up.GetM().PropsLock != 0) || up.CanTrade == 1
}

func DefaultGetMProp(pid int) *MProp {
	return nil
}

func SetMPropFunc(f func(pid int) *MProp) {
	getMProp = f
}
