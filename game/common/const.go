package common

const (
	LoginExpireTime = 60 * 60 * 1000 //登录过期时间，单位为秒，客户端每10分钟会发一个请求更新过期时间
	// Redis keys
	MWelcome        = "MWELCOME"
	MTIMECONFIG     = "MTIMECONFIG"
	MPetKey         = "MPET"
	MSkillKey       = "MSKILL"
	MPropKey        = "MPROP"
	MSERIESKey      = "MSERIES"
	MMapKey         = "MMAP"
	MGpcKey         = "MGPC"
	MGpcGroupKey    = "MGPCGroup"
	MZbAttributeKey = "MZBATTRIBUTE"
	MFightTimeKey   = "MFIGHTTIME"
	ExpListKey      = "EXPLIST"

	//
	// Game setting
	MaxLevel = 130
	MaxJinbi = 10000 * 10000 * 10
	WxJing   = 1
	WxMu     = 2
	WxShui   = 3
	WxHuo    = 4
	WxTu     = 5
	WxShen   = 6
	WxSS     = 7
)

const (
	LenUserNameMin = 5
	LenUserNameMax = 12
	LenPasswordMax = 16
	Layout         = "2006-01-02 15:04:05"
	LenDesc        = 150
	LenAddr        = 50
	LenLimit       = 10
)
const (
	ArticleKey   = "article"
	ExipreSecond = 360
	IPKey        = "blog::ip"
)

//提示信息
const (
	MsgLoginSucc          = "登录成功"
	MsgResistSucc         = "注册成功"
	MsgGetUserInfoSucc    = "获取用户信息成功"
	MsgUpdateUserInfoSucc = "修改用户信息成功"
	MsgDelUserSucc        = "删除用户成功"
	MsgGetLinksSucc       = "获取用户友链成功"
	MsgCreateLinksSucc    = "创建用户友链成功"
	MsgUpdateLinksSucc    = "修改用户友链成功"
	MsgDelLinksSucc       = "删除用户友链成功"
	MsgCreateCategorySucc = "创建分类成功"
	MsgUpdateCategorySucc = "修改分类成功"
	MsgDelCategorySucc    = "删除分类成功"
	MsgGetCategorySucc    = "获取分类成功"
	MsgCreateArticleSucc  = "发布文章成功"
	MsgSaveArticleSucc    = "保存草稿成功"
	MsgGetArticleSucc     = "获取文章成功"
	MsgDelArticleSucc     = "删除文章成功"
	MsgUpdateArticleSucc  = "修改文章成功"
	MsgSaveImageSucc      = "保存图片成功"
	MsgGetExtSucc         = "获取扩展信息成功"
)

var DifficulMapId = []int{16, 100, 103, 106, 109, 112, 115, 118, 121, 128, 131, 134, 137, 140, 145, 148}
