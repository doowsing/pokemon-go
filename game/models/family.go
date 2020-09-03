package models

type Family struct {
	Id             int    `gorm:"column:id;primary"`
	Name           string `gorm:"column:name"`
	Info           string `gorm:"column:info"`
	CreatorIdStr   string `gorm:"column:creator_id"`
	PresidentId    int    `gorm:"column:president_id"`
	Honor          int    `gorm:"column:honor"`
	Level          int    `gorm:"column:level"`
	ShopLevel      int    `gorm:"column:shop_level"`
	NumberOfMember int    `gorm:"column:number_of_member"`
	CreateTime     int    `gorm:"column:create_time"`
	VictoryTimes   int    `gorm:"column:victory_times"`
	FailedTimes    int    `gorm:"column:failed_times"`
}

func (f *Family) TableName() string {
	return "guild"
}

type FamilyBag struct {
	Id       int `gorm:"column:id;primary"`
	FamilyId int `gorm:"column:guild_id"`
	Pid      int `gorm:"column:pid"`
	Sums     int `gorm:"column:sums"`
}

func (f *FamilyBag) TableName() string {
	return "guild_bag"
}

type FamilySetting struct {
	Id               int    `gorm:"column:id;primary"`
	Level            int    `gorm:"column:level"`
	NeedHonor        int    `gorm:"column:need_honor"`
	NeedProps        string `gorm:"column:need_props"`
	NeedMemberNumber int    `gorm:"column:need_member_number"`
	NeedItemsForShop string `gorm:"column:need_items_for_shop"`
	MaxMemberNumber  int    `gorm:"column:max_member_number"`
	Welfare          string `gorm:"column:welfare"`
}

func (f *FamilySetting) TableName() string {
	return "guild_settings"
}

type FamilyMember struct {
	UserId       int `gorm:"column:member_id"`
	FamilyId     int `gorm:"column:guild_id"`
	JoinTime     int `gorm:"column:join_time"`
	Authority    int `gorm:"column:priv"`
	Contribution int `gorm:"column:contribution"`
	Honor        int `gorm:"column:honor"`
}

func (f *FamilyMember) TableName() string {
	return "guild_members"
}

func (f *FamilyMember) IsPresident() bool {
	return f.Authority == 3
}

func (f *FamilyMember) IsElder() bool {
	return f.Authority == 2
}

func (f *FamilyMember) IsNormalMember() bool {
	return f.Authority == 1
}

type FamilyChallenge struct {
	Id              int    `gorm:"column:id;primary"`
	ChallengerId    int    `gorm:"column:challenger_id"`    // 挑战者
	DefenserId      int    `gorm:"column:defenser_id"`      // 防守者
	ChallengeMsg    string `gorm:"column:challenge_msg"`    // 战术内容
	CreateTime      int    `gorm:"column:create_time"`      // 发出时间
	ChallengerScore int    `gorm:"column:challenger_score"` // 挑战者积分
	DefenserScore   int    `gorm:"column:defenser_score"`   // 防守者积分
	Flags           int    `gorm:"column:flags"`            // 是否接受
}

func (f *FamilyChallenge) TableName() string {
	return "guild_challenges"
}
