package models

import (
	"pokemon/pkg/utils"
)

type User struct {
	ID        int    `gorm:"column:id;primary_key"`
	Account   string `gorm:"column:name;unique_index:account;not null"`
	PasswdMd5 string `gorm:"column:secret;not null"`
	Nickname  string `gorm:"column:nickname;unique_index:nickname;not null"`
	Sex       string `gorm:"column:sex;not null"`
	Regtime   int    `gorm:"column:regtime;not null"`
	Lastvtime int    `gorm:"column:lastvtime"` //最后访问时间
	Money     int    `gorm:"column:money;default:0;not null"`
	Yb        int    `gorm:"column:yb;default:0;not null"`
	Headimg   string `gorm:"column:headimg;default:0;not null"`
	TaskStr   string `gorm:"column:task;not null"`
	Mbid      int    `gorm:"column:mbid"`
	Password  string `gorm:"column:password"` // 是否被封号
	//BanTalkTime     time.Time `gorm:"column:ban_talk_time"`                  //禁言时间:单位：秒
	//BanTime         time.Time `gorm:"column:ban_time"`                       //禁号时间:单位：秒
	Maxdblexptime    int    `gorm:"column:maxdblexptime;default:0"` // 剩余多倍禁言时间
	Dbexpflag        int    `gorm:"column:dbexpflag;default:0"`     //是否处于自动战斗
	AutofitStartTime int    `gorm:"column:dblstime;default:0"`      //是否处于自动战斗
	Fightbb          int    `gorm:"column:fightbb"`
	HeartTime        int    `gorm:"column:heart_time"`              //不知道是什么时间，和登录时间有点相似
	AutoFightTimeM   int    `gorm:"column:sysautosum;default:800"`  //金币战斗次数
	Sysautotime      int    `gorm:"column:sysautotime;default:800"` //多倍经验到期时间
	BagPlace         int    `gorm:"column:maxbag;default:30;not null"`
	BasePlace        int    `gorm:"column:maxbase;default:40;not null"`
	McPlace          int    `gorm:"column:maxmc;default:10;not null"`
	TgPlace          int    `gorm:"column:tgmax;default:1"`       //托管空间
	InMap            int    `gorm:"column:inmap;default:1"`       //所在地图ID
	OpenMap          string `gorm:"column:openmap;default:'1'"`   //所有开启的地图，用,分割
	Autofitflag      int    `gorm:"column:autofitflag;default:0"` //是否处于自动战斗
	Secid            int    `gorm:"column:secid"`                 //封号
	BotMapId         int    `gorm:"column:bot_map_id"`
	Useyb            int    `gorm:"column:useyb;default:0"`           //本月使用元宝
	Score            int    `gorm:"column:score;default:0"`           //积分
	Vip              int    `gorm:"column:vip;default:0"`             //vip积分
	VipYb            int    `gorm:"column:vipyb;default:0"`           //用于计算vip等级
	TgTime           int    `gorm:"column:tgtime;default:0"`          //托管时间
	VipLast          int    `gorm:"column:viplast;default:0"`         //上个月的vip积分
	ChallengeRecord  string `gorm:"column:fighttop;default:'0:0'"`    // 挑战记录
	AutoFightTimeYb  int    `gorm:"column:maxautofitsum;default:800"` //元宝自动战斗次数
	Prestige         int    `gorm:"column:prestige;default:0"`        //威望
	Jprestige        int    `gorm:"column:jprestige;default:0"`       //贵族威望
	ActiveScore      int    `gorm:"column:active_score"`              // 不知道什么的积分
	TaskLog          string `gorm:"column:tasklog"`
	LoginTime        int    `gorm:"column:logintime"`
	OnlineTime       int    `gorm:"column:onlinetime"`
	Paihang          int    `gorm:"column:Paihang"`
	PaiMoney         int    `gorm:"column:paimoney;default:0"`
	FromType         int    `gorm:"column:fromtype"`
	Wg               int    `gorm:"column:wg"`       //不知道是啥
	McPwd            string `gorm:"column:fieldpwd"` //牧场密码
	CkPwd            string `gorm:"column:ckpwd"`    //仓库密码
	Phone            string `gorm:"column:phone"`
	AllRmb           int    `gorm:"column:all_rmb;default:0"`
	Code             string `gorm:"column:code"`
	//
	//Sj              int       `gorm:"column:sj;default:0;not null"`
	//PaiSj          int       `gorm:"column:pai_sj;default:0"`
	//PaiYb          int       `gorm:"column:pai_yb;default:0"`
	//InvitationCode string    `gorm:"column:invitation_code;default:0"`
	//ShowTimes      int       `gorm:"column:show_times;default:5"`  //宠物展示次数
	//MergeTimes     int       `gorm:"column:merge_times;default:0"` //合成失败次数
	//Tgt            int       `gorm:"column:tgt;default:0"`         //现TT怪物数
	//TgtTime        int       `gorm:"column:tgt_time"`              //上次TT重新开启时间
	//WelfareTime    time.Time `gorm:"column:welfare_time"`          //上次领取帮会福利时间
	//GuideStep      int       `gorm:"column:guide_step;default:0"`  //新人引导路线ID
	//SSCzl          float64   `gorm:"column:ss_czl;default:0.0"`    //抽取的神圣成长
	//TeamAutoTimes    int    `gorm:"column:team_auto_times;default:0"` //组队战斗次数

}

// TableName sets the insert table name for this struct type
func (user *User) TableName() string {
	return "player"
}

func (user *User) AfterFind() (err error) {
	//user.Nickname = utils.ToUtf8(user.Nickname)
	//user.Sex = utils.ToUtf8(user.Sex)
	return
}
func (user *User) BeforeSave() (err error) {
	//user.Nickname = utils.ToGbk(user.Nickname)
	return
}

func (user *User) IsValid() bool {
	return user.Regtime > 0
}

type UserInfo struct {
	ID                  int    `gorm:"column:id;primary_key"`
	Uid                 int    `gorm:"column:uid"`
	Bbshow              int    `gorm:"column:bbshow;default:0"`
	HechengNums         int    `gorm:"column:hecheng_nums;default:0"`
	OnlineTime          int    `gorm:"column:onlinetime"`
	LoginTime           int    `gorm:"column:logintime"`
	Sj                  int    `gorm:"column:sj;default:0"`
	Paisj               int    `gorm:"column:paisj;default:0"`
	ActiveLastvTime     int    `gorm:"column:active_lastvtime"`
	Ml                  int    `gorm:"column:ml;default:0"`
	Merge               int    `gorm:"column:merge;default:0"`
	RequestMerge        int    `gorm:"column:request_merge;default:0"`
	Request             int    `gorm:"column:request;default:0"`
	Tgt                 int    `gorm:"column:tgt;default:0"`
	TgtTime             int    `gorm:"column:tgttime;default:0"`
	TgtI                int    `gorm:"column:tgt_i;default:0"`
	GpcGroupId          int    `gorm:"column:gpc_group_id;default:0"`
	TgLastTime          int    `gorm:"column:tglasttime;default:0"`
	NoMergeTime         int    `gorm:"column:nomergetime"`
	Send                string `gorm:"column:send;default:0"`
	ChChengbb           int    `gorm:"column:chchengbb"`
	GetWelfareTime      string `gorm:"column:get_welfare_time"`
	PrizeItems          string `gorm:"column:prize_every_day"`
	GuildRequest        int    `gorm:"column:guild_request"`
	TeamAutoTimes       int    `gorm:"column:team_auto_times;default:0"`
	NewGuideStep        int    `gorm:"column:new_guide_step;default:0"`
	Consumption2expDay  string `gorm:"column:consumption2exp_day"`
	RegAddStr           string `gorm:"column:reg_add_str"`
	CzlSS               int    `gorm:"column:czl_ss;default:0"`
	ExpGotStep          int    `gorm:"column:exp_got_step;default:0"`
	OnlineTimeToday     int    `gorm:"column:onlinetime_today;default:0"`
	LastOnlineTime      int    `gorm:"column:last_onlinetime;default:0"`
	LastOnlineDay       int    `gorm:"column:last_online_day;default:0"`
	LastLoginTime       int    `gorm:"column:last_logintime;default:0"`
	TiaoZhan            int    `gorm:"column:tiaozhan;default:1"`
	BuffStatus          string `gorm:"column:buff_status"`
	ChouquChongwu       string `gorm:"column:chongqu_chongwu"`
	NowAchievementTitle string `gorm:"column:now_Achievement_title"`
	FHasGetPrize        string `gorm:"column:F_has_get_prize"`
	FUserCardInfo       string `gorm:"column:F_User_Card_Info"`
	FHasTitle           string `gorm:"column:F_Has_Title"`
	FMedicineBuff       string `gorm:"column:F_Medicine_Buff"`
	FSaoleiPoints       int    `gorm:"column:F_saolei_points;default:1"`
	Paiyb               int    `gorm:"column:paiyb;default:0"`
}

func (u *UserInfo) TableName() string {
	return "player_ext"
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
}

func (u *CardTitle) TableName() string {
	return "t_card_to_title"
}

func (c *CardTitle) AfterFind() (err error) {
	c.Name = utils.ToUtf8(c.Name)
	c.NeedCard = utils.ToUtf8(c.NeedCard)
	c.NeedDep = utils.ToUtf8(c.NeedDep)
	return
}

type AoyunPlayer struct {
	Id     int `gorm:"column:id;primary_key"`
	Uid    int `gorm:"column:uid;default:0"`
	Stime  int `gorm:"column:stime;default:0"`
	Tid    int `gorm:"column:tid;default:0"`
	Qsum   int `gorm:"column:qsums;default:0"`
	OkSum  int `gorm:"column:oksum;default:0"`
	Times  int `gorm:"column:times;default:0"`
	Result int `gorm:"column:result;default:0"`
}

func (u *AoyunPlayer) TableName() string {
	return "aoyun_player"
}
