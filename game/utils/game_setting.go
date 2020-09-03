package utils

import (
	"encoding/json"
	"io/ioutil"
)

const (
	EQTITLECOLOR   = "#ED9037" // 装备显示头部颜色,增加默认。
	EQBASECOLOR    = "#FEFDFA" // 装备显示基础属性颜色
	EQPLUSCOLOR    = "#0067CB" // 装备显示附加属性颜色
	EQSPECIALCOLOR = "#14FD10" // 装备显示特殊属性颜色
	EQGLODCOLOR    = "#FED625" //金色
	EQGREENCOLOR   = "#A8A7A4" //灰色
	MaxJinBi       = 1000000000
)

type FbSetStruct struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	CoolTime int    `json:"time"`
	Level    int    `json:"lv"`
	Gpcs     []int  `json:"gpcs"`
}

type DtStruct struct {
	FbSet  []*FbSetStruct `json:"fb_set"`
	MapSet struct {
		NeedKeyMap  []int `json:"need_key_map"`
		NormalMap   []int `json:"normal_map"`
		FbMap       []int `json:"fb_map"`
		MultipleMap []int `json:"multiple_map"`
		SsMap       []int `json:"ss_map"`
	} `json:"map_set"`
	WxName              []string          `json:"wx_name"`
	ZbPositionName      map[int]string    `json:"zb_position_name"`
	ZbEffectName        map[string]string `json:"zb_effect_name"`
	ZbSpecialEffectName map[string]string `json:"zb_special_effect_name"`
	PropTypeName        map[int]string    `json:"prop_type_name"`
	PetAttrName         map[string]string `json:"pet_attr_name"`
	PropColor           map[int]string    `json:"prop_color"`
}

var setFileDir string
var gameSettings *DtStruct

//var wxNames = []string{"所有", "金", "木", "水", "火", "土", "神", "神圣"}
//
//// 非普通副本地图
//var normalMapIds = []int{
//	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 16, //玛雅大陆地图
//	20,                                     // 圣兽云殿
//	100, 103, 106, 109, 112, 115, 118, 121, // 新大陆地图
//	128,                          //遗忘的宫殿
//	131, 134, 137, 140, 145, 148, // 神圣大陆地图
//}
//
//// 普通副本地图
//var fbMapIds = []int{11, 12, 13, 14, 50, 124, 127, 143, 144, 151}
//
//// 多难度地图
//var multipleMapIds = []int{
//	100, 103, 106, 109, 112, 115, 118, 121, // 新大陆地图
//	128,                          //遗忘的宫殿
//	131, 134, 137, 140, 145, 148, // 神圣大陆地图
//}
//
//// 神圣副本地图
//var ssMapIds = []int{131, 134, 137, 140, 145, 148, 143, 144}
//
//var fbSet = map[int]FbSetStruct{}
//
//var zbPositions = map[int]string{
//	1:  "头部",
//	2:  "身体",
//	3:  "脚部",
//	4:  "武器",
//	5:  "项链",
//	6:  "戒指",
//	7:  "翅膀",
//	8:  "手镯",
//	9:  "宝石",
//	10: "道具",
//	11: "卡片左",
//	12: "卡片右",
//}
//
//var zbEffect2Name = map[string]string{
//	"openpet": "获得一个宠物",
//	"mc":      "防御",
//	"openmap": "开启一个地图",
//	"hp":      "生命",
//	"mp":      "魔法",
//	"hits":    "命中",
//	"miss":    "闪避",
//	//"kx":      "抗性",
//	"speed":     "速度",
//	"ac":        "攻击",
//	"hprate":    "生命",
//	"mprate":    "魔法",
//	"acrate":    "攻击",
//	"mcrate":    "防御",
//	"hitsrate":  "命中",
//	"missrate":  "闪避",
//	"speedrate": "速度",
//}
//
//var zbSpecialEffect2Name = map[string]string{
//	"dxsh": "伤害抵消",
//	"shjs": "伤害加深",
//	"shft": "反弹伤害",
//}
//
//var propTypeNames = map[int]string{
//	1:  "辅助类",
//	2:  "增益类",
//	3:  "捕捉类",
//	4:  "收集类",
//	5:  "技能书类",
//	6:  "卡片类",
//	7:  "进化类",
//	8:  "合体类",
//	9:  "装备类",
//	10: "精炼类",
//	11: "精炼类",
//	12: "礼包类",
//	13: "特殊类",
//	14: "功能类",
//	15: "宠物卵",
//	16: "合成类",
//	17: "水晶类",
//	18: "特殊回复类",
//	19: "涅槃加成",
//	22: "魔法石",
//	23: "神圣转生道具",
//	24: "卡片类",
//	25: "宝石类",
//	26: "洗练石",
//	27: "合成保低石类",
//	28: "刮刮卡类",
//	29: "奇石类",
//	30: "扫雷道具类",
//	31: "扫雷道具类",
//	32: "扫雷道具类",
//	50: "魔塔回复类",
//	51: "魔塔复活类",
//	52: "魔塔解密类",
//	53: "魔塔杀伤类",
//	54: "魔塔BUFF",
//	55: "魔塔洗点类",
//	56: "魔塔洗点类",
//	57: "魔塔出战卷",
//	58: " 魔塔增益类",
//}
//
//var attrName = map[string]string{
//
//	"name":        "宝贝名字",
//	"wx":          "五行",
//	"ac":          "攻击",
//	"mc":          "防御",
//	"hp":          "生命",
//	"srchp":       "生命",
//	"srcmp":       "魔法",
//	"speed":       "速度",
//	"hits":        "命中",
//	"miss":        "闪避",
//	"imgstand":    "站立图片名",
//	"imgack":      "攻击图片名",
//	"imgdie":      "施法图片名",
//	"skillist":    "技能列表",
//	"czl":         "成长率",
//	"kx":          "抗性",
//	"remakelevel": "进化等级",
//	"remakeid":    "进化后的宝贝ID",
//	"remakepid":   "进化需要道具ID",
//	"nowexp":      "当前经验",
//	"lexp":        "升级经验",
//	"subyl":       "减晕",
//	"subsl":       "减睡",
//	"subdl":       "减毒",
//	"subxl":       "减虚",
//	"subfl":       "减防",
//	"subhl":       "减缓",
//	"subkl":       "减抗",
//	"headimg":     "头像图片",
//	"cardimg":     "卡片图",
//	"effectimg":   "展示图",
//	"bbdesc":      "宝宝介绍",
//}
//
//var propColor = map[int]string{
//	1: "#FEFDFA",
//	2: "#0067CB",
//	3: "#9833DC",
//	4: "#14FD10",
//	5: "#FED625",
//	6: "#ED9037",
//}

func UpdateSetting(fileDir string) error {
	if fileDir != "" {
		setFileDir = fileDir
	}
	data, err := ioutil.ReadFile(setFileDir)
	if err != nil {
		return err
	}
	dt := DtStruct{}
	err = json.Unmarshal(data, &dt)
	if err != nil {
		return err
	}
	gameSettings = &dt
	return nil
}

func GetWxName(wx int) string {
	wxName := "未知"
	if wx <= len(gameSettings.WxName) && wx > -1 {
		wxName = gameSettings.WxName[wx]
	}
	return wxName
}

func GetZbPositionName(p int) string {
	return gameSettings.ZbPositionName[p]
}

func GetVaryNameStr(varyname int) string {
	name, ok := gameSettings.PropTypeName[varyname]
	if ok {
		return name
	}
	return "未知"
}

func GetAttrName(attr string) string {
	name, ok := gameSettings.PetAttrName[attr]
	if ok {
		return name
	}
	return ""
}

func GetPropColor(colorType int) string {
	color, ok := gameSettings.PropColor[colorType]
	if !ok {
		color = EQTITLECOLOR
	}
	return color
}

func GetZbEffectName(effect string) string {
	name, ok := gameSettings.ZbEffectName[effect]
	if !ok {
		return ""
	}
	return name
}

func GetZbSpecialEffectName(effect string) string {
	name, ok := gameSettings.ZbSpecialEffectName[effect]
	if !ok {
		return ""
	}
	return name
}

func IsNormalMap(mapId int) bool {
	for _, i := range gameSettings.MapSet.NormalMap {
		if i == mapId {
			return true
		}
	}
	return false
}

func IsFbMap(mapId int) bool {
	for _, i := range gameSettings.MapSet.FbMap {
		if i == mapId {
			return true
		}
	}
	return false
}

func IsMultipleMap(mapId int) bool {
	for _, m := range gameSettings.MapSet.MultipleMap {
		if m == mapId {
			return true
		}
	}
	return false
}

func IsHuPoWuMap(mapId int) bool {
	return mapId == 125
}

func IsTTMap(mapId int) bool {
	return mapId == 126
}

func IsYiWangMap(mapId int) bool {
	return mapId == 128
}

func IsNeedKeyMap(mapId int) bool {
	for _, m := range gameSettings.MapSet.NeedKeyMap {
		if m == mapId {
			return true
		}
	}
	return false
}

func IsSSZhanchangMap(mapId int) bool {
	return mapId == -1
}

func IsFamilyMap(mapId int) bool {
	return mapId == -2
}

func IsSsMap(mapId int) bool {
	for _, m := range gameSettings.MapSet.SsMap {
		if m == mapId {
			return true
		}
	}
	return false
}

func GetFbSet(mapId int) *FbSetStruct {
	for _, m := range gameSettings.FbSet {
		if m.Id == mapId {
			return m
		}
	}
	return nil
}
