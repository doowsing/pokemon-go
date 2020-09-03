package common

import "fmt"

//
//const (
//	Ok_Code           = 200 //get put ok
//	PostOk_Code       = 201
//	DelOk_Code        = 204
//	RequsetErr_Code   = 400 //参数错误等
//	UnAuthorized_Code = 401
//	Forbidden_Code    = 403
//	NotFound_Code     = 404
//	InternalErr_Code  = 500
//)

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (err Err) Error() string {
	return fmt.Sprintf("ErrCode:%d , ErrMsg:%s", err.Code, err.Msg)
}

//返回结果中的错误码表示了用户调用API 的结果。其中，code 为公共错误码，其适用于所有模块的 API 接口。
//若 code 为 0，表示调用成功，否则，表示调用失败。当调用失败后，用户可以根据下表确定错误原因并采取相应措施。
var (
	ErrClientParams   = Err{Code: 4000, Msg: "缺少必要参数，或者参数值/路径格式不正确。"}
	ErrAuthorized     = Err{Code: 4100, Msg: "签名鉴权失败，请参考文档中鉴权部分。"}
	ErrAccoutDeny     = Err{Code: 4200, Msg: "帐号被封禁，或者不在接口针对的用户范围内等。"}
	ErrNotFound       = Err{Code: 4300, Msg: "资源不存在，或者访问了其他用户的资源。"}
	ErrMethodNotAllow = Err{Code: 4400, Msg: "协议不支持，请参考文档说明。"}
	ErrSignParams     = Err{Code: 4500, Msg: "签名错误，请参考文档说明。"}
	ErrInternal       = Err{Code: 5000, Msg: "服务器内部出现错误，请稍后重试。"}
	ErrApiClose       = Err{Code: 6200, Msg: "当前接口处于停服维护状态，请稍后重试。"}
)

//模块错误码
var (
	ErrUserExist          = Err{Code: 10000, Msg: "用户已存在,请修改后重试。"}
	ErrUserLogin          = Err{Code: 10010, Msg: "用户名或密码错误,请检查后重试。"}
	ErrUserNameUnExist    = Err{Code: 10020, Msg: "用户名不存在,请检查后重试。"}
	ErrUserNoExist        = Err{Code: 10030, Msg: "用户不存在,请检查后重试。"}
	ErrUserNameFormat     = Err{Code: 10040, Msg: fmt.Sprintf("用户名需字母开头长度(%d~%d)位字母数字_。", LenUserNameMin, LenUserNameMax)}
	ErrUserPwdFormat      = Err{Code: 10050, Msg: fmt.Sprintf("密码需字母开头长度(%d~%d)位字母数字_。", LenUserNameMin, LenPasswordMax)}
	ErrUserDescLen        = Err{Code: 10060, Msg: fmt.Sprintf("描述长度不能超过%d位,请改正后重试。", LenDesc)}
	ErrUserAddrLen        = Err{Code: 10070, Msg: fmt.Sprintf("地址长度不能超过%d位,请改正后重试。", LenAddr)}
	ErrUserEmailFormat    = Err{Code: 10080, Msg: "邮箱格式不正确,请检查后重试。"}
	ErrUserNickNameFormat = Err{Code: 10090, Msg: fmt.Sprintf("昵称长度在(%d~%d)之间,请改正后重试。", LenUserNameMin, LenUserNameMax)}
	ErrUpdateParams       = Err{Code: 10100, Msg: "修改用户信息,参数必填其一。"}
	ErrUserLinksNoExist   = Err{Code: 10110, Msg: "用户友链数据不存在,请检查后重试。"}
	ErrCategoryNoExist    = Err{Code: 10120, Msg: "该分类不存在,请检查后重试。"}
	ErrCategoryExist      = Err{Code: 10130, Msg: "该分类已存在,请修改后重试。"}
	ErrArticleNoExist     = Err{Code: 10140, Msg: "该文章不存在,请修改后重试。"}
	ErrUploadLenNotAllow  = Err{Code: 10150, Msg: "图片上传个数不允许,请修改后重试。"}
	ErrUploadExtNotAllow  = Err{Code: 10160, Msg: "仅支持jpg,jpeg,png格式图片,请修改后重试。"}
	ErrCollectSource      = Err{Code: 10170, Msg: "文章采集失败,请检查采集Url。"}
)

// 道具使用状态码

var (
	ErrUsePropNoExist         = NewErr(600, "没有相关的道具！")
	ErrUsePropNoModel         = NewErr(601, "该道具设定已被删除！")
	ErrUsePropNoLevel         = NewErr(602, "使用等级不足！")
	ErrUsePropOnlySS          = NewErr(603, "只有神圣宠物可以使用！")
	ErrUsePropPetErr          = NewErr(604, "宠物资料出现错误！")
	ErrUsePropZbSuccess       = NewErr(605, "装备成功！")
	ErrUsePropSetErr          = NewErr(606, "道具配置出错，请联系管理员！")
	ErrUsePropTgKcNeed        = NewErr(607, "托管扩充等级需求不对！")
	ErrUsePropTgKcSuccess     = NewErr(608, "托管空间扩充成功！")
	ErrUsePropKcBbErr         = NewErr(609, "扩充背包失败，请检查使用条件！")
	ErrUsePropKcBbSuccess     = NewErr(610, "扩充背包成功！")
	ErrUsePropKcCkErr         = NewErr(611, "扩充仓库失败，请检查使用条件！")
	ErrUsePropKcCkSuccess     = NewErr(612, "扩充仓库成功！")
	ErrUsePropKcTgTimeSuccess = NewErr(613, "扩充托管时间成功！")
	ErrUsePropMapOpened       = NewErr(614, "地图已经开启了！")
	ErrUsePropMapOpenSuccess  = NewErr(615, "开启地图成功！")
	ErrUsePropAutoTimeSuccess = NewErr(616, "添加次数成功！")
	ErrUsePropNoBgPlace       = NewErr(617, "背包空间不足，请至少预留3个位置！")
	ErrUsePropNeedKey         = NewErr(618, "使用礼包需要消耗钥匙！")
	ErrUsePropSuccess         = NewErr(619, "使用成功！")
	ErrUsePropNoSS            = NewErr(620, "神圣宠物不可以使用该道具！")
	ErrUsePropMaxLevel        = NewErr(621, "宠物已到最高级，无法再继续使用道具经验！")
	ErrUsePropAddExpSuccess   = NewErr(622, "宠物使用道具经验成功！")
	ErrUsePropHcNoNeed        = NewErr(623, "合成图纸所需道具不足！")
	ErrUsePropCZNoNeed        = NewErr(623, "没有可重铸的物品！")
	ErrUsePropCarryPCNT       = NewErr(624, "您只能携带3个宝宝,使用道具失败！请将宝宝放入牧场再使用")
	ErrUsePropCannotUse       = NewErr(625, "该道具不可以直接使用！")
	ErrUsePropDBErr           = NewErr(626, "数据库连接出错！")
	ErrUsePropAddZsSuccess    = NewErr(627, "增加宠物展示次数成功！")
)

var (
	// 添加道具响应码
	ErrAddPropFullBP   = NewErr(700, "背包已满!")
	ErrAddPropMNoExist = NewErr(701, "道具设定不存在！")
	ErrAddPropSuccess  = NewErr(702, "道具添加成功！")
	ErrAddPropError    = NewErr(703, "道具添加出错")
)

var (
	// 战斗响应码
	ErrFightEnterMapSuccess = NewErr(800, "进入地图正常！")
	ErrFightNoMap           = NewErr(801, "地图无效！")
	ErrFightStartSuccess    = NewErr(802, "成功进入战斗！")
	ErrFightNoCzl           = NewErr(803, "成长不足以进入地图！")
	ErrFightNoLevel         = NewErr(804, "等级不足以进入地图！")
	ErrFightOnlySS          = NewErr(805, "该地图只有神圣宠物可进入！")
	ErrFightSetErr          = NewErr(806, "地图设定出错！")
	ErrFightNoGpc           = NewErr(807, "怪物设定不存在！")
)

func ErrNoMsg(code int) Err {
	return Err{Code: code, Msg: ""}
}

func NewErr(code int, msg string) Err {
	return Err{Code: code, Msg: msg}
}
