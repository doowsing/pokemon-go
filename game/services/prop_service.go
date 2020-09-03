package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
	"math/rand"
	"pokemon/common/persistence"
	"pokemon/common/rcache"
	"pokemon/game/models"
	"pokemon/game/repositories"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PropService struct {
	BaseService
	Repo *repositories.PropRepository
}

func NewPropService(osrc *OptService) *PropService {
	us := &PropService{Repo: repositories.NewPropRepository()}
	us.SetOptSrc(osrc)
	return us
}
func (ps *PropService) GetMProp(MPropId int) (*models.MProp, error) {
	return common.GetMProp(MPropId), nil
}

// 业务代码

func (ps *PropService) GetProps(UserId int) []*models.UProp {
	props := []*models.UProp{}
	ps.GetDb().Where("uid=?", UserId).Find(&props)
	return props
}
func (ps *PropService) GetProp(UserId, PropId int, lock bool) *models.UProp {
	prop := &models.UProp{}
	rs := ps.GetDb()
	if lock {
		rs = rs.Set("gorm:query_option", "FOR UPDATE")
	}
	if rs.Where("uid =? and Id = ?", UserId, PropId).First(prop).RowsAffected > 0 {
		return prop
	}
	return nil
}
func (ps *PropService) GetPropById(PropId int, lock bool) *models.UProp {
	prop := &models.UProp{}
	rs := ps.GetDb()
	if lock {
		rs = rs.Set("gorm:query_option", "FOR UPDATE")
	}
	if rs.Where("Id = ?", PropId).First(prop).RowsAffected > 0 {
		return prop
	}
	return nil
}

// 用于获取其他玩家的道具信息，购买的时候需要
func (ps *PropService) GetOtherPropById(PropId int, lock bool) *models.UProp {

	prop := &models.UProp{}
	rs := ps.GetDb()
	if lock {
		rs = rs.Set("gorm:query_option", "FOR UPDATE")
	}
	if rs.Where("Id = ?", PropId).First(prop).RowsAffected > 0 {
		return prop
	}
	return nil
}

func (ps *PropService) GetPropByPid(UserId, PropId int, lock bool) *models.UProp {
	prop := &models.UProp{}
	rs := ps.GetDb()
	if lock {
		rs = rs.Set("gorm:query_option", "FOR UPDATE")
	}
	if rs.Where("uid =? and pid = ?", UserId, PropId).First(prop).RowsAffected > 0 {
		return prop
	}
	return nil
}

func (ps *PropService) AddPropSums(PropId, num int) bool {
	return ps.GetDb().Model(&models.UProp{ID: PropId}).Update(UpMap{"sums": gorm.Expr("sums+?", num)}).RowsAffected > 0
}

func (ps *PropService) AddProp(userId, propId, num int, checkPlace bool) bool {
	mprop := common.GetMProp(propId)
	if mprop == nil {
		return false
	}
	addFlag := false
	if mprop.IsVary() {

		prop := ps.GetPropByPid(userId, propId, false)
		if prop != nil && ps.AddPropSums(prop.ID, num) {
			addFlag = true
		}
	} else {
		num = 1
	}
	if !addFlag {
		if checkPlace && ps.BagLeftPlace(userId) == 0 {
			fmt.Printf("create prop error:no enough place:%d\n", checkPlace)
			return false
		}
		prop := &models.UProp{
			Pid:   propId,
			Uid:   userId,
			Sums:  num,
			Sell:  mprop.SellJb,
			Stime: utils.NowUnix(),
		}
		Result := ps.GetDb().Create(prop)
		if Result.Error != nil {
			fmt.Printf("create prop error:%s\n", Result.Error)
			return false
		}
	}
	return true
}

func (ps *PropService) AddOrCreateProp(userId, propId, num int, checkPlace bool) (*models.UProp, bool) {
	mprop := common.GetMProp(propId)
	addFlag := false
	var prop *models.UProp
	if mprop.IsVary() {
		prop = ps.GetPropByPid(userId, propId, false)
		if prop != nil && ps.AddPropSums(prop.ID, num) {
			addFlag = true
		}
	} else {
		num = 1
	}
	if !addFlag {
		if checkPlace && ps.BagLeftPlace(userId) == 0 {
			fmt.Printf("create prop error:no enough place:%d\n", checkPlace)
			return nil, false
		}
		prop = &models.UProp{
			Pid:   propId,
			Uid:   userId,
			Sums:  num,
			Sell:  mprop.SellJb,
			Stime: utils.NowUnix(),
		}
		Result := ps.GetDb().Create(prop)
		if Result.Error != nil {
			fmt.Printf("create prop error:%s\n", Result.Error)
			return nil, false
		}
	}
	return prop, true
}

func (ps *PropService) UsePropByPid(UserId, Pid int) string {
	prop := ps.GetPropByPid(UserId, Pid, true)
	if prop == nil || prop.Sums == 0 {
		return "无相关魔法石，无法满足释放魔法需要的魔力T_T下次再来吧。"
	}
	prop.GetM()
	return ps.UseProp(UserId, prop)
}

func (ps *PropService) UsePropById(UserId, upid int) string {
	prop := ps.GetProp(UserId, upid, true)
	if prop == nil || prop.Sums == 0 {
		return "物品不存在!"
	}
	prop.GetM()
	if prop.MModel.VaryName == 22 {
		return `占卜石,请在占卜屋中占卜时使用!<br/><span style="cursor:pointer;color:#ff0000;font-size:14px;font-weight:bold" onclick="(\'gw\').contentWindow.location=\'/function/zhanbuwu.php\'"><strong>点击这里前往“占卜屋”！</strong></span>`
	}
	return ps.UseProp(UserId, prop)
}

func (ps *PropService) UseProp(UserId int, prop *models.UProp) string {
	user := ps.OptSvc.UserSrv.GetUserById(UserId)
	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	switch prop.MModel.VaryName {
	case 9:
		// 装备类
		mainPet := ps.OptSvc.PetSrv.GetPet(user.ID, user.Mbid)
		if mainPet == nil {
			return "您还没有设置主战宝宝，不能进行装备！"
		}
		mainPet.GetM()
		if prop.MModel.Requires != "" {
			requireList := strings.Split(prop.MModel.Requires, ",")
			for _, rqr := range requireList {
				if rqr == "" {
					continue
				}
				if items := strings.Split(rqr, ":"); len(items) > 1 {
					item0, err := strconv.Atoi(items[1])
					if items[0] == "lv" && err != nil && mainPet.Level < item0 {
						return "宝宝五行不匹配!"
					} else if items[0] == "wx" && err != nil && mainPet.MModel.Wx != item0 {
						return "宝宝等级不够!"
					}
				}
			}
		}
		newZbs := []string{}
		find := false
		for _, v := range strings.Split(mainPet.Zb, ",") {
			if items := strings.Split(v, ":"); len(items) > 1 {
				if com.StrTo(items[0]).MustInt() == prop.MModel.Position {
					newZbs = append(newZbs, fmt.Sprintf("%s:%d", items[0], prop.ID))
					ps.OffZb(com.StrTo(items[1]).MustInt())
					find = true
				} else {
					newZbs = append(newZbs, v)
				}
			}
		}
		if !find {
			newZbs = append(newZbs, fmt.Sprintf("%d:%d", prop.MModel.Position, prop.ID))
		}
		if !ps.EquipPet(mainPet.ID, prop.ID) {
			return "装备失败！1"
		}
		if !ps.SetPetEquips(mainPet.ID, strings.Join(newZbs, ",")) {
			return "装备失败！2"
		}
		ps.OptSvc.FightSrv.DelZbAttr(mainPet.ID)
		ps.OptSvc.Commit()
		return "装备成功！"
	case 13:
		// 空间扩充类
		if !ps.DecrPropById(prop.ID, 1) {
			return "物品不存在!"
		}
		if prop.Pid == 1203 {
			if user.TgPlace >= 2 {
				return "您只能使用此卷扩充一次托管所！"
			} else if ps.OptSvc.UserSrv.AddTgPlace(UserId) {
				ps.OptSvc.Commit()
				return "使用托管所扩充卷轴（一）成功!"
			} else {
				return "用户不存在!"
			}
		} else if prop.Pid == 1204 {
			if user.TgPlace == 1 {
				return "请先使用托管所扩充卷（一）扩充您的托管所!"
			} else if user.TgPlace >= 3 {
				return "您只能使用此卷扩充一次托管所！"
			} else if ps.OptSvc.UserSrv.AddTgPlace(UserId) {
				ps.OptSvc.Commit()
				return "使用托管所扩充卷轴（一）成功!"
			} else {
				return "用户不存在!"
			}
		}
		EffectItems := strings.SplitN(prop.MModel.Effect, ":", 2)
		if EffectItems[0] == "zhanshi" {
			// 宠物展示卷
			if ps.OptSvc.UserSrv.AddShowTimes(user.ID, com.StrTo(EffectItems[1]).MustInt()) {
				ps.OptSvc.Commit()
				return "恭喜您使用宠物展示卷成功增加" + EffectItems[1] + "次展示机会！"
			} else {
				return "使用宠物展示卷失败"
			}
		} else if EffectItems[0] == "addsj" {
			// 水晶卡
			if ps.OptSvc.UserSrv.AddSj(UserId, com.StrTo(EffectItems[1]).MustInt()) {
				ps.OptSvc.Commit()
				return "恭喜您得到了" + EffectItems[1] + "个水晶！"
			} else {
				return "用户不存在!"
			}
		} else if EffectItems[0] == "addyb" {
			// 元宝卡
			if ps.OptSvc.UserSrv.AddYb(UserId, com.StrTo(EffectItems[1]).MustInt()) {
				ps.OptSvc.Commit()
				return "恭喜您得到了" + EffectItems[1] + "个元宝！"
			} else {
				return "用户不存在!"
			}
		} else if EffectItems[0] == "addbag" {
			// 实星背包升级卷轴
			num := com.StrTo(EffectItems[1]).MustInt()
			if user.BagPlace+num > 200 {
				num = 200 - user.BagPlace
			}
			if user.BagPlace < 150 {
				return "您的背包没有达到150，不能使用此道具扩展！"
			} else if user.BagPlace >= 200 {
				return "您的背包已经有200格了，不能再使用此道具扩展！"
			} else if ps.OptSvc.UserSrv.AddBagPlace(UserId, num, 200) {
				ps.OptSvc.Commit()
				return "恭喜您背包格子扩充了" + EffectItems[1] + "格！"
			} else {
				return "您的背包已经有200格了，不能再使用此道具扩展！"
			}
		} else if EffectItems[0] == "addck" {
			// 实星仓库升级卷轴
			num := com.StrTo(EffectItems[1]).MustInt()
			if user.BasePlace+num > 200 {
				num = 200 - user.BasePlace
			}
			if user.BasePlace < 150 {
				return "您的仓库没有达到150，不能使用此道具扩展！"
			} else if user.BasePlace >= 200 {
				return "您的仓库已经有200格了，不能再使用此道具扩展！"
			} else if ps.OptSvc.UserSrv.AddCkPlace(UserId, num, 200) {
				ps.OptSvc.Commit()
				return "恭喜您仓库格子扩充了" + EffectItems[1] + "格！"
			} else {
				return "您的仓库已经有200格了，不能再使用此道具扩展！"
			}
		} else if EffectItems[0] == "addbag1" {
			// 空星背包升级卷轴
			num := com.StrTo(EffectItems[1]).MustInt()
			if user.BagPlace+num > 300 {
				num = 300 - user.BagPlace
			}
			if user.BagPlace < 200 {
				return "您的背包没有达到200，不能使用此道具扩展！"
			} else if user.BagPlace >= 300 {
				return "您的背包已经有300格了，不能再使用此道具扩展！"
			} else if ps.OptSvc.UserSrv.AddBagPlace(UserId, num, 300) {
				ps.OptSvc.Commit()
				return "恭喜您背包格子扩充了" + EffectItems[1] + "格！"
			} else {
				return "您的背包已经有300格了，不能再使用此道具扩展！"
			}
		} else if EffectItems[0] == "addck1" {
			// 空星仓库升级卷轴
			num := com.StrTo(EffectItems[1]).MustInt()
			if user.BasePlace+num > 300 {
				num = 300 - user.BasePlace
			}
			if user.BasePlace < 200 {
				return "您的仓库没有达到200，不能使用此道具扩展！"
			} else if user.BasePlace >= 300 {
				return "您的仓库已经有300格了，不能再使用此道具扩展！"
			} else if ps.OptSvc.UserSrv.AddCkPlace(UserId, num, 300) {
				ps.OptSvc.Commit()
				return "恭喜您仓库格子扩充了" + EffectItems[1] + "格！"
			} else {
				return "您的仓库已经有300格了，不能再使用此道具扩展！"
			}
		} else if EffectItems[0] == "addmcplace" {
			if len(EffectItems) < 4 {
				return "道具配置出错！"
			}
			baseNum := com.StrTo(EffectItems[1]).MustInt()
			maxNum := com.StrTo(EffectItems[2]).MustInt()
			addNum := com.StrTo(EffectItems[3]).MustInt()
			if user.McPlace < baseNum {
				return fmt.Sprintf("您的牧场没有达到%d，不能使用此道具扩展！", baseNum)
			}
			if user.McPlace >= maxNum {
				return fmt.Sprintf("您的牧场已经达到%d，不能使用此道具扩展！", maxNum)
			}
			if ps.OptSvc.UserSrv.AddMcPlace(UserId, addNum, maxNum) {
				ps.OptSvc.Commit()
				return fmt.Sprintf("恭喜您牧场格子扩充了%d格！", addNum)
			} else {
				return fmt.Sprintf("您的牧场已经达到%d，不能使用此道具扩展！", maxNum)
			}
		} else if EffectItems[0] == "tuoguan" {
			// 托管卷，增加宠物托管时间
			if ps.OptSvc.UserSrv.AddTgTime(UserId, com.StrTo(EffectItems[1]).MustInt()) {
				ps.OptSvc.Commit()
				return "使用" + EffectItems[1] + "小时托管卷成功!"
			} else {
				return "使用失败！"
			}

		} else if EffectItems[0] == "openmap" {
			NewMapId := EffectItems[1]
			OpendMaps := strings.Split(user.OpenMap, ",")
			IsFind := false
			for _, v := range OpendMaps {
				if v == NewMapId {
					IsFind = true
					break
				}
			}
			if IsFind {
				return prop.MModel.Name + " 对应的地图已经打开了!"
			} else {
				OpendMaps = append(OpendMaps, NewMapId)
				if ps.OptSvc.UserSrv.SetMaps(UserId, strings.Join(OpendMaps, ",")) {
					ps.OptSvc.Commit()
					return prop.MModel.Name + " 对应的地图打开成功!"
				} else {
					return "地图打开失败，请确认包裹中有打开该地图对应的钥匙!"
				}
			}
		} else if (prop.Pid >= 200 && prop.Pid <= 202) || (prop.Pid >= 1342 && prop.Pid <= 1344) {
			// 普通背包、仓库、牧场升级卷、高级牧场升级卷
			num := 6
			switch prop.MModel.Name {
			case "背包升级卷轴":
				if user.BagPlace >= 96 {
					num = 0
				} else {
					if user.BagPlace+num >= 96 {
						num = 96 - user.BagPlace
					}
					if !ps.OptSvc.UserSrv.AddBagPlace(UserId, num, 96) {
						num = 0
					}
				}
				break
			case "仓库升级卷轴":
				if user.BasePlace >= 96 {
					num = 0
				} else {
					if user.BasePlace+num >= 96 {
						num = 96 - user.BasePlace
					}
					if !ps.OptSvc.UserSrv.AddCkPlace(UserId, num, 96) {
						num = 0
					}
				}
				break
			case "牧场升级卷轴":
				if user.BasePlace >= 40 {
					num = 0
				} else {
					if user.BasePlace+num >= 40 {
						num = 40 - user.BasePlace
					}
					if !ps.OptSvc.UserSrv.AddMcPlace(UserId, num, 40) {
						num = 0
					}
				}
				break
			case "高级背包升级卷轴":
				num = 1
				if user.BagPlace < 96 {
					return "您的背包格子还没扩展到96格，请先买其它道具扩展到96格才能再用此道具扩展!"
				} else if user.BagPlace >= 96 {
					num = 0
				} else {
					if user.BagPlace+num >= 96 {
						num = 96 - user.BagPlace
					}
					if !ps.OptSvc.UserSrv.AddBagPlace(UserId, num, 150) {
						num = 0
					}
				}
				break
			case "高级仓库升级卷轴":
				num = 1
				if user.BasePlace < 96 {
					return "您的背包格子还没扩展到96格，请先买其它道具扩展到96格才能再用此道具扩展!"
				} else if user.BasePlace >= 96 {
					num = 0
				} else {
					if user.BasePlace+num >= 96 {
						num = 96 - user.BasePlace
					}
					if !ps.OptSvc.UserSrv.AddBagPlace(UserId, num, 150) {
						num = 0
					}
				}
				break
			case "高级牧场升级卷轴":
				num = 1
				if user.McPlace < 40 {
					return "您的牧场格子还没扩展到40格，请先买其它道具扩展到40格才能再用此道具扩展!"
				} else if user.McPlace >= 40 {
					num = 0
				} else {
					if user.McPlace+num >= 40 {
						num = 40 - user.McPlace
					}
					if !ps.OptSvc.UserSrv.AddMcPlace(UserId, num, 40) {
						num = 0
					}
				}
				break
			default:
				return "参数出错！"
			}
			if num == 0 {
				return "已经扩展到极限，如需再扩展请买其它道具!"
			} else {
				ps.OptSvc.Commit()
				return "使用道具 " + prop.MModel.Name + "成功！"
			}
		} else if EffectItems[0] == "exp" {
			// 多倍经验卷轴使用 exp:3:3600
			// 这个暂时不做
		} else if EffectItems[0] == "autofree" || EffectItems[0] == "autoteam" || EffectItems[0] == "auto" {
			// 自动战斗卷轴使用 金币版，元宝版，团队版
			// 这里需要改！！！20191106
			num := com.StrTo(EffectItems[1]).MustInt()
			msg := ""
			if EffectItems[0] == "autofree" {
				if ps.OptSvc.UserSrv.AddJbAuto(UserId, num) {
					msg = "增加金钱版自动战斗次数 " + EffectItems[1]
				}

			} else if EffectItems[0] == "auto" {
				if ps.OptSvc.UserSrv.AddYbAuto(UserId, num) {
					msg = "增加元宝版自动战斗次数 " + EffectItems[1]
				}
			} else if EffectItems[0] == "autoteam" {
				if ps.OptSvc.UserSrv.AddTeamAuto(UserId, num) {
					msg = "增加组队自动战斗次数 " + EffectItems[1]
				}
			} else {
				return ""
			}
			if msg != "" {
				ps.OptSvc.Commit()
				return msg
			}
			return "使用道具失败!"
		}
		break
	case 4:
		// 号码开奖类
		break
	case 12, 22:
		// 礼包类
		if !ps.DecrPropById(prop.ID, 1) {
			return "物品不存在!"
		}
		bagCnt := ps.GetCarryPropsCnt(UserId)
		if user.BagPlace-bagCnt < 3 {
			return "请留至少三个空格子！"
		}
		// 检查使用要求
		if prop.MModel.Requires != "" {
			if requires := strings.Split(prop.MModel.Requires, ":"); len(requires) > 1 && requires[0] == "lv" {
				mainpet := ps.OptSvc.PetSrv.GetPet(user.ID, user.Mbid)
				if mainpet.Level < com.StrTo(requires[1]).MustInt() {
					if prop.MModel.VaryName == 12 {
						return "您没有达到相应的等级，不能开启该宝箱！"
					} else {
						return "您没有达到相应的等级，不能进行占卜！"
					}
				}
			}
		}
		EffectItems := strings.Split(prop.MModel.Effect, ",")
		for _, v := range EffectItems {
			items := strings.SplitN(v, ":", 2)
			if len(items) < 1 {
				continue
			}
			switch items[0] {
			case "needkey":
				if !ps.DecrPropByPid(UserId, com.StrTo(items[1]).MustInt(), 1) {
					if prop.MModel.VaryName == 12 {
						return "您没有开启宝箱的钥匙!"
					} else {
						return "您没有占卜的钥匙!"
					}
				}
				break
			case "giveitems":
				giveitems := append(strings.Split(items[1], ","), EffectItems[1:]...)
				if len(giveitems)+bagCnt > user.BagPlace {
					return "背包空间不足！"
				}
				getPnames := []string{}
				for _, str := range giveitems {
					if pitems := strings.Split(str, ":"); len(pitems) > 1 {
						if ps.AddProp(UserId, com.StrTo(pitems[0]).MustInt(), com.StrTo(pitems[1]).MustInt(), false) {
							getPnames = append(getPnames, common.GetMProp(com.StrTo(pitems[0]).MustInt()).Name+"*"+pitems[1])
						} else {
							return "使用道具出错！"
						}
					}
				}
				if len(getPnames) > 0 {
					ps.OptSvc.Commit()
					return "获得道具 " + strings.Join(getPnames, ",")
				} else {
					return "使用道具失败"
				}
			case "randitem":
				randitems := strings.Split(items[1], "|")
				for _, str := range randitems {
					if pitems := strings.Split(str, ":"); len(pitems) > 3 {
						pid := com.StrTo(pitems[0]).MustInt()
						num := com.StrTo(pitems[1]).MustInt()
						rateNum := com.StrTo(pitems[2]).MustInt()
						gonggao := com.StrTo(pitems[3]).MustInt()
						if rand.Intn(rateNum)+1 == 1 {
							if ps.AddProp(UserId, pid, num, false) {
								if gonggao == 2 {
									// 发公告

									if prop.MModel.VaryName == 12 {
										AnnounceAll(user.Nickname, fmt.Sprintf("使用 %s ,幸运地得到自然女神的祝福,获得了 %d 个%s", prop.MModel.Name, num, common.GetMProp(pid).Name))
									} else {
										AnnounceAll(user.Nickname, fmt.Sprintf("使用 %s ,虔诚的占卜感动了自然女神,获得了 %d 个%s", prop.MModel.Name, num, common.GetMProp(pid).Name))
									}
								}
								ps.OptSvc.Commit()
								return fmt.Sprintf("获得道具 %s %d 个", common.GetMProp(pid).Name, num)
							} else {
								return "使用道具出错！"
							}
						}
					}
				}
				return "使用道具出错！"
			default:
				continue
			}
		}
		break
	case 2:
		// 增益类
		mainpet := ps.OptSvc.PetSrv.GetPet(user.ID, user.Mbid)
		mainpet.GetM()
		if mainpet == nil {
			return "请设置主宠后再使用！"
		}
		if mainpet.MModel.Wx == 7 && prop.MModel.Requires != "__SS__" {
			return "神圣宠物无法使用此类物品！"
		}
		if mainpet.MModel.Wx != 7 && prop.MModel.Requires == "__SS__" {
			return "非神圣宠物无法使用此类物品！"
		}
		if !ps.DecrPropById(prop.ID, 1) {
			return "物品不存在!"
		}
		EffectItems := strings.Split(prop.MModel.Effect, ":")
		switch EffectItems[0] {
		case "addexp":
			EffectItems[1] = regexp.MustCompile(`[(|)]`).ReplaceAllString(EffectItems[1], "")
			items := strings.Split(EffectItems[1], ",")
			getexp := 0
			if len(items) > 1 {
				minexp := com.StrTo(items[0]).MustInt()
				maxexp := com.StrTo(items[1]).MustInt()
				if minexp == maxexp {
					getexp = minexp
				} else {
					getexp = minexp + rand.Intn(maxexp-minexp+1)
				}
				fmt.Printf("getexp:%s to %d, %s to %d\n", items[0], minexp, items[1], maxexp)
			} else {
				getexp = com.StrTo(items[0]).MustInt()
			}
			if ps.OptSvc.PetSrv.IncreaseExp2Pet(mainpet, getexp) {
				ps.OptSvc.Commit()
				return "获得经验 " + utils.IntToStr(getexp)
			} else {
				return "宠物已经满级，不能再升级了！"
			}
			break
		case "addczl":
			if mainpet.CqFlag > 0 {
				return "这个宠物抽取过成长,不能再使用这个道具!"
			}
			num := com.StrTo(EffectItems[1]).MustFloat64()
			oldCzl := com.StrTo(mainpet.Czl).MustFloat64()
			newCzl := oldCzl + num
			if mainpet.MModel.Wx == 7 {
				if float64(common.GetSSJhRule(mainpet.Bid).MaxCzl) < newCzl {
					newCzl = float64(common.GetSSJhRule(mainpet.Bid).MaxCzl)
				}
			}
			if ps.OptSvc.PetSrv.SetPetCzl(mainpet.ID, utils.CzlStr(newCzl)) {
				ps.OptSvc.Commit()
				return fmt.Sprintf("主宠物永久增加 %s 成长！", utils.CzlStr(newCzl-oldCzl))
			} else {
				return "使用道具失败！"
			}
		case "addac", "addmc", "addhits", "addmiss", "addhp", "addspeed", "addmp":
			attr := strings.ReplaceAll(EffectItems[0], "add", "")
			if attr == "hp" || attr == "mp" {
				attr = "src" + attr
			}
			if ps.OptSvc.PetSrv.AddPetAttribute(mainpet.ID, attr, com.StrTo(EffectItems[1]).MustInt()) {
				ps.OptSvc.Commit()
				return fmt.Sprintf("主宠物永久增加 %s %s！", EffectItems[1], utils.GetAttrName(attr))
			} else {
				return "使用道具失败！"
			}
			break
		case "weiwang":
			attr := "prestige"
			if ps.OptSvc.PetSrv.AddPetAttribute(mainpet.ID, attr, com.StrTo(EffectItems[1]).MustInt()) {
				ps.OptSvc.Commit()
				return fmt.Sprintf("增加 %s %s！", EffectItems[1], utils.GetAttrName(attr))
			} else {
				return "使用道具失败！"
			}
		case "add_cq_czl":
			if ps.OptSvc.UserSrv.AddSsCzl(UserId, com.StrTo(EffectItems[1]).MustInt()) {
				ps.OptSvc.Commit()
				return "获得抽取成长" + EffectItems[1] + "点！"
			} else {
				return "使用道具失败！"
			}
		case "add_zc_jifen":
			if ps.OptSvc.UserSrv.AddZcScore(UserId, EffectItems[1]) {
				ps.OptSvc.Commit()
				return "操作成功！"
			} else {
				return "使用道具失败！"
			}
		}
	case 24:
		// 卡片类
		if !ps.DecrPropById(prop.ID, 1) {
			return "物品不存在!"
		}
		var msg string
		record := &models.CardRecord{}
		ps.GetDb().Where("uid=? and card_pid=?", UserId, prop.MModel.ID).First(record)
		if record.ID > 0 {
			ps.GetDb().Model(record).Update(gin.H{"sum": gorm.Expr("sum+1")})
			msg = "使用过的卡片使用成功"
		} else {
			record = &models.CardRecord{
				Uid:     UserId,
				CardPid: prop.MModel.ID,
				Sum:     1,
			}
			ps.GetDb().Create(record)
			msg = "恭喜您第一次使用卡片成功"
		}
		ps.OptSvc.NpcSrv.UpdateCardUserTitles(user)
		ps.OptSvc.Commit()
		return msg

	case 16:
		// 图纸合成类
		if !ps.DecrPropById(prop.ID, 1) {
			return "物品不存在!"
		}
		bagCnt := ps.GetCarryPropsCnt(UserId)
		if user.BagPlace-bagCnt < 3 {
			return "请留至少三个空格子！"
		}
		EffectItems := strings.SplitN(prop.MModel.Effect, ":", 2)
		switch EffectItems[0] {
		case "hecheng":
			//图纸合成 格式：hecheng:(956:10|957:10|958:10|1026:1|1055:2):1012:1
			strItems := strings.Split(EffectItems[1], "):")
			requireStr := strings.ReplaceAll(strItems[0], "(", "") //要求的东西
			getpropStr := strItems[1]                              //获得的东西
			for _, v := range strings.Split(requireStr, "|") {
				if items := strings.Split(v, ":"); len(items) > 1 {
					p := ps.GetPropByPid(UserId, com.StrTo(items[0]).MustInt(), false)
					needNum := com.StrTo(items[1]).MustInt()
					if p == nil || p.Sums < needNum || p.Zbpets != 0 || !ps.DecrPropById(p.ID, needNum) {
						return "你的材料不足，无法制作！"
					}
				}
			}
			getpropList := []string{}
			for _, v := range strings.Split(getpropStr, "|") {
				if items := strings.Split(v, ":"); len(items) > 1 {
					pid := com.StrTo(items[0]).MustInt()
					if ps.AddProp(UserId, pid, com.StrTo(items[1]).MustInt(), false) {
						p := common.GetMProp(pid)
						getpropList = append(getpropList, fmt.Sprintf("%s %s个", p.Name, items[1]))
					}
				}
			}
			if len(getpropList) > 0 {
				ps.OptSvc.Commit()
				return "恭喜您,制作成功!获得道具 " + strings.Join(getpropList, "，")
			} else {
				return "制作失败！"
			}
		case "chongzhu":
			break
		case "random_combine":
			break
		}
		return ""
	case 15:
		// 宠物卵
		if cnt := ps.OptSvc.PetSrv.GetCarryPetCnt(UserId); cnt >= 3 {
			return "您只能携带3个宝宝,使用道具失败！<br/>[系统推荐]：您可以把身上携带的宝宝放入到牧场！"
		}
		if !ps.DecrPropById(prop.ID, 1) {
			return "物品不存在!"
		}
		EffectItems := strings.Split(prop.MModel.Effect, ":")
		if EffectItems[0] == "openpet" && len(EffectItems) > 1 {
			if ok, _ := ps.OptSvc.PetSrv.CreatPetById(user, com.StrTo(EffectItems[1]).MustInt()); ok {
				ps.OptSvc.Commit()
				return "使用道具成功"
			}
		}
		return "使用道具失败！参数错误"
	case 14:
		// 军功令
		ssUser := &models.SSBattleUser{}
		ps.GetDb().Where("uid=?", UserId).First(ssUser)
		if ssUser.Id == 0 {
			return "您目前没有参加战场活动，不能使用此道具！"
		}
		if !ps.DecrPropById(prop.ID, 1) {
			return "物品不存在!"
		}
		effectItems := strings.Split(prop.MModel.Effect, ":")
		if len(effectItems) <= 1 {
			return "道具配置出错！"
		}
		num := com.StrTo(effectItems[1]).MustInt()
		ps.GetDb().Model(ssUser).Update(gin.H{"jgvalue": gorm.Expr("jgvalue+?", num)})
		ps.OptSvc.Commit()
		return fmt.Sprintf("恭喜您，使用道具成功，您获得了 %d 点军功！", num)
	default:
		return ""
	}
	return ""
}

func (ps *PropService) EquipPet(petId, zbId int) bool {
	return ps.GetDb().Model(&models.UProp{}).Where("Id=?", zbId).Update(UpMap{"zbpets": petId, "zbing": 1}).RowsAffected > 0
}

func (ps *PropService) SetPetEquips(petId int, zbs string) bool {
	return ps.GetDb().Model(&models.UPet{ID: petId}).Update(UpMap{"zb": zbs}).RowsAffected > 0
}

func (ps *PropService) DecrPropById(PropID, num int) bool {
	return ps.GetDb().Model(&models.UProp{}).Where("Id = ? and sums>= ?", PropID, num).Update(UpMap{"sums": gorm.Expr("sums - ?", num)}).RowsAffected > 0
}

func (ps *PropService) DecrProp(userId, PropID, num int) bool {
	return ps.GetDb().Model(&models.UProp{}).Where("Id = ? and sums>= ? and uid=?", PropID, num, userId).Update(UpMap{"sums": gorm.Expr("sums - ?", num)}).RowsAffected > 0
}

func (ps *PropService) DecrPropByPid(userId, PropID, num int) bool {
	return ps.GetDb().Model(&models.UProp{}).Where("pid = ? and uid = ? and sums>=?", PropID, userId, num).Update(UpMap{"sums": gorm.Expr("sums - ?", num)}).RowsAffected > 0
}

func (ps *PropService) GetPZbs(petid int) []models.UProp {
	zbs := []models.UProp{}
	ps.GetDb().Where("zbpets = ?", petid).Find(&zbs)
	return zbs
}

func (ps *PropService) GetCarryProps(userId int, clean bool) []*models.UProp {
	_props := []models.UProp{}
	ps.GetDb().Where("uid = ? and sums>0 and zbing!=1", userId).Order("pid").Find(&_props)
	props := []*models.UProp{}
	for _, prop := range _props {
		prop1 := prop
		prop1.GetM()
		props = append(props, &prop1)
	}

	if clean {
		sort.Slice(props, func(i, j int) bool {
			return props[i].MModel.VaryName < props[j].MModel.VaryName
		})
	}
	return props
}

func (ps *PropService) GetCarryPropsCnt(userId int) int {
	cnt := 0
	ps.GetDb().Model(&models.UProp{}).Where("uid = ? and sums>0 and zbing=0", userId).Count(&cnt)
	return cnt
}

func (ps *PropService) GetCarryPropsByVaryName(userId int, clean bool, varyname int) []*models.UProp {
	orderStr := "ub.pid"
	if clean {
		orderStr = "p.varyname, " + orderStr
	}
	where := "ub.uid = ? and ub.sums>0 and ub.zbing=0"
	if varyname > 0 {
		where += fmt.Sprintf(" and p.varyname=%d", varyname)
	}
	results, err := ps.GetDb().Table("userbag ub").Select(`ub.Id, ub.pid, ub.sums, p.varyname, p.Name, p.sell`).Joins("inner join props p on p.Id=ub.pid").Where(where, userId).Order(orderStr).Rows()
	if err != nil {
		fmt.Println("error1 : ", err)
		return nil
	}
	defer results.Close()
	ups := []*models.UProp{}
	for results.Next() {
		up := models.UProp{}

		var id, pid, sums, varyname, name, sell string

		_ = results.Scan(&id, &pid, &sums, &varyname, &name, &sell)
		up.ID = com.StrTo(id).MustInt()
		up.Pid = com.StrTo(pid).MustInt()
		up.Sums = com.StrTo(sums).MustInt()
		up.MModel = &models.MProp{
			Name:     name,
			VaryName: com.StrTo(varyname).MustInt(),
			SellJb:   com.StrTo(sell).MustInt(),
		}
		_ = up.MModel.AfterFind()
		up1 := up
		ups = append(ups, &up1)
	}
	return ups
}

func (ps *PropService) GetPropInfoJson(userId, upId, pid int) gin.H {
	var prop *models.UProp
	if upId != 0 {
		prop = ps.GetPropById(upId, false)
		if prop == nil {
			return nil
		}
		if prop.GetM(); prop.MModel == nil {
			return nil
		}
	} else {
		prop = &models.UProp{}
		prop.Pid = pid
		if prop.GetM(); prop.MModel == nil {
			return nil
		}
	}
	propMap := gin.H{
		"Id":       prop.ID,
		"Name":     prop.MModel.Name,
		"varyname": prop.MModel.VaryName,
		"color":    prop.MModel.PropsColor,
	}
	trade := "不可交易"
	if (prop.CanTrade == 0 && prop.MModel.PropsLock != 0) || prop.CanTrade == 1 {
		trade = "可交易"
	}
	propMap["can_trade"] = trade
	if prop.ID > 0 {
		isExpire := "过期"
		if prop.MModel.Expire == 0 {
			isExpire = "永久"
		} else {
			now := ps.NowUnix()
			expireTime := prop.Stime + prop.MModel.Expire
			if expireTime > now {
				isExpire = time.Unix(int64(expireTime), 0).Format("2006-01-02 15:04:05")
			} else {
				isExpire = "过期"
			}
		}
		propMap["expire"] = isExpire

	}
	propMap["pid"] = prop.MModel.ID
	switch prop.MModel.VaryName {
	case 9:
		zbInfo := gin.H{}
		plus := "不可强化"
		if prop.MModel.PlusFlag == 1 {
			plus = "可强化"
		}
		zbInfo["can_strength"] = plus
		zbInfo["position"] = utils.GetZbPositionName(prop.MModel.Position)
		if prop.ID > 0 {
			items := strings.Split(prop.PlusTmsEft, ",")
			if len(items) > 1 {
				zbInfo["strength"] = gin.H{
					"Level": com.StrTo(items[0]).MustInt() + 1,
					"value": items[1],
				}
			}
		}

		// 主属性
		effectItems := strings.Split(prop.MModel.Effect, ":")
		zbInfo["main_effect"] = fmt.Sprintf("+%s %s", effectItems[1], utils.GetZbEffectName(effectItems[0]))

		// 装备条件
		requiresInfo := []string{}
		if prop.MModel.Requires != "" {
			if requires := strings.Split(prop.MModel.Requires, ","); len(requires) > 1 {

				if lvReq := strings.Split(requires[0], ":"); len(lvReq) > 1 {
					requiresInfo = append(requiresInfo, fmt.Sprintf("等级需求：%s级", lvReq[1]))
				}
				if wxReq := strings.Split(requires[1], ":"); len(wxReq) > 1 {
					requiresInfo = append(requiresInfo, fmt.Sprintf("五行需求：%s系", utils.GetWxName(com.StrTo(wxReq[1]).MustInt())))
				}
			}
		}
		zbInfo["requires"] = requiresInfo

		// 副属性部分
		plusEffectItems := strings.Split(prop.MModel.PlusEffect, ",")
		effectName := ""
		otherEffect := []string{}
		for _, effect := range plusEffectItems {
			if items := strings.Split(effect, ":"); len(items) > 1 {
				if effectName = utils.GetZbEffectName(items[0]); effectName != "" {
					otherEffect = append(otherEffect, fmt.Sprintf("+%s %s", items[1], effectName))
				} else if effectName = utils.GetZbSpecialEffectName(items[0]); effectName != "" {
					otherEffect = append(otherEffect, fmt.Sprintf("%s %s", effectName, items[1]))
				} else {
					effectDep := ""
					switch items[0] {
					case "szmp":
						effectDep = "伤害的%s转化为MP"
						break
					case "sdmp":
						effectDep = "伤害的%s以MP抵消"
						break
					case "addmoney":
						effectDep = "战斗胜利获得金币增加%s点"
						break
					case "hitshp":
						effectDep = "偷取伤害的%s转化为生命"
						break
					case "hitsmp":
						effectDep = "偷取伤害的%s转化为魔法"
						break
					case "time":
						effectDep = "战斗等待时间减少%s秒"
						break
					}
					if effectDep != "" {
						otherEffect = append(otherEffect, fmt.Sprintf(effectDep, items[1]))
					}

				}
			}
		}
		zbInfo["other_effect"] = otherEffect

		// 宝石卡槽和属性
		if prop.MModel.PlusNum == 0 {
			zbInfo["hole_num"] = "无卡槽"
		} else if prop.MModel.PlusNum > 0 && prop.FHoleInfo != "" {
			effectItems = strings.Split(prop.FHoleInfo, ",")
			zbInfo["hole_num"] = fmt.Sprintf("卡槽数：%d/%d", len(effectItems), prop.MModel.PlusNum)
			for _, v := range effectItems {
				items := strings.Split(v, ":")
				effectDep := ""
				switch items[0] {
				case "ac":
					effectDep = "增加攻击%s"
					break
				case "mc":
					effectDep = "增加防御%s"
					break
				case "hits":
					effectDep = "增加命中%s"
					break
				case "miss":
					effectDep = "增加闪避%s"
					break
				case "hp":
					effectDep = "增加HP上限%s"
					break
				case "mp":
					effectDep = "增加MP上限%s"
					break
				case "speed":
					effectDep = "增加速度%s"
					break
				case "sdmp":
					effectDep = "将受到伤害的%s以MP抵消"
					break
				case "szmp":
					effectDep = "将受到伤害的%s转化为MP"
					break
				case "hitshp":
					effectDep = "命中吸取伤害的%s转化为自身HP"
					break
				case "hitsmp":
					effectDep = "命中吸取伤害的%s转化为自身MP"
					break
				case "dxsh":
					effectDep = "伤害抵销%s"
					break
				case "shjs":
					effectDep = "对敌人造成的伤害增加%s"
					break
				case "crit":
					effectDep = "会心一击率增加%s"
					break
				}
				if effectDep != "" {
					zbInfo["hole_value"] = fmt.Sprintf("宝石效果："+effectDep, items[1])
				}
			}
		} else {
			zbInfo["hole_num"] = fmt.Sprintf("卡槽数：0/%d", prop.MModel.PlusNum)
		}

		// 套装效果
		if prop.MModel.Series != "" && prop.MModel.Series != "0" {
			series := gin.H{}
			serieItems := strings.Split(prop.MModel.Series, ":")
			serieList := strings.Split(serieItems[1], "|")
			var petzbids []string
			if prop.Zbpets > 0 {
				for _, v := range ps.GetPZbs(prop.Zbpets) {
					v.GetM()
					if com.IsSliceContainsStr(serieList, strconv.Itoa(v.Pid)) {
						petzbids = append(petzbids, strconv.Itoa(v.Pid))
					}
				}
			}

			zbList := []gin.H{}
			for _, v := range serieList {
				mprop := common.GetMProp(com.StrTo(v).MustInt())
				if com.IsSliceContainsStr(petzbids, strconv.Itoa(mprop.ID)) {
					zbList = append(zbList, gin.H{"Name": mprop.Name, "have": true})
				} else {
					zbList = append(zbList, gin.H{"Name": mprop.Name, "have": false})
				}
			}
			series["zb_list"] = zbList
			series["Name"] = fmt.Sprintf("%s(%d/%d)", serieItems[0], len(petzbids), len(serieList))

			effectList := []gin.H{}
			serieEffects := strings.Split(prop.MModel.SeriesEffect, ",")
			for i, v := range serieEffects {
				if v != "" && v != "0" {
					effectItems := strings.Split(v, ":")
					if len(effectItems) < 2 {
						continue
					}
					effectDep := ""
					effectName := ""
					if effectName = utils.GetZbEffectName(effectItems[0]); effectName != "" {
						effectDep = fmt.Sprintf(`+%s %s`, effectItems[1], effectName)
					} else if effectName = utils.GetZbSpecialEffectName(effectItems[0]); effectName != "" {
						effectDep = fmt.Sprintf(`%s %s`, effectName, effectItems[1])
					} else {
						switch effectItems[0] {
						case "szmp":
							effectDep = "伤害的%s转化为MP"
							break
						case "sdmp":
							effectDep = "伤害的%s以MP抵消"
							break
						case "addmoney":
							effectDep = "战斗胜利获得金币增加%s点"
							break
						case "hitshp":
							effectDep = "偷取伤害的%s转化为生命"
							break
						case "hitsmp":
							effectDep = "偷取伤害的%s转化为魔法"
							break
						case "time":
							effectDep = "战斗等待时间减少%s秒"
							break
						}
						if effectDep != "" {
							effectDep = fmt.Sprintf(effectDep, effectItems[1])
						}
					}
					if effectDep != "" {
						if i < len(petzbids) {
							effectList = append(effectList, gin.H{"value": fmt.Sprintf("(%d)套装：%s", i+1, effectDep), "have": true})
						} else {
							effectList = append(effectList, gin.H{"value": fmt.Sprintf("(%d)套装：%s", i+1, effectDep), "have": false})
						}
					}
				}
			}
			series["effect_list"] = effectList
			zbInfo["series"] = series
		}
		propMap["zb_info"] = zbInfo
		break
	case 25:
		// 宝石
		require := "镶嵌部位:"
		if prop.MModel.Requires != "" {
			for _, v := range strings.Split(prop.MModel.Requires, ",") {
				if items := strings.Split(v, ":"); len(items) > 1 {
					if items[0] == "postion" {

						for _, p := range strings.Split(items[1], "|") {
							require += utils.GetZbPositionName(com.StrTo(p).MustInt()) + " "
						}
					} else if items[0] == "color" {
						color2names := map[string]string{
							"2": "蓝色装备",
							"3": "紫色装备",
							"4": "绿色装备",
							"5": "黄色装备",
							"6": "橙色装备",
						}
						require += "只能镶嵌" + color2names[items[1]]
					}
				}
			}
		} else {
			require = "无需求"
		}
		propMap["hole_require"] = require
		break
	default:
		break
	}
	propMap["usage"] = prop.MModel.Usages
	return propMap
}

func (ps *PropService) GeneratePropInfo(userId, upId, pid int) string {
	var prop *models.UProp
	if upId != 0 {
		prop = ps.GetProp(userId, upId, false)
		if prop == nil {
			return ""
		}
		if prop.GetM(); prop.MModel == nil {
			return ""
		}
	} else {
		prop = &models.UProp{}
		prop.MModel = common.GetMProp(pid)
	}

	infoHtml := fmt.Sprintf(`<font color="%s"><b>%s</b></font><br/>`, utils.GetPropColor(prop.MModel.PropsColor), prop.MModel.Name)

	// 是否可交易
	trade := "不可交易"
	if (prop.CanTrade == 0 && prop.MModel.PropsLock != 0) || prop.CanTrade == 1 {
		trade = "可交易"
	}
	infoHtml += fmt.Sprintf(`<font color=%s>%s</font><br/>`, utils.EQGREENCOLOR, trade)

	// 过期时间
	expireHtml := ""
	if prop.ID > 0 {
		isExpire := "过期"
		if prop.MModel.Expire == 0 {
			isExpire = "永久"
		} else {
			now := ps.NowUnix()
			expireTime := prop.Stime + prop.MModel.Expire
			if expireTime > now {
				isExpire = time.Unix(int64(expireTime), 0).Format("2006-01-02 15:04:05")
			} else {
				isExpire = "过期"
			}
		}
		expireHtml = fmt.Sprintf(`<font color=%s>%s</font><br />`, utils.EQBASECOLOR, isExpire)
	}
	infoHtml += expireHtml
	switch prop.MModel.VaryName {
	case 9:
		strongValue := ""

		infoHtml = fmt.Sprintf(`<font color="%s"><b>%s</b></font><br/>`, utils.GetPropColor(prop.MModel.PropsColor), prop.MModel.Name)
		if prop.ID > 0 {
			items := strings.Split(prop.PlusTmsEft, ",")
			if len(items) > 1 {
				strongValue = fmt.Sprintf(`<font color=red>+%s</font>`, items[1])
				infoHtml = fmt.Sprintf(`<font color="%s"><b>%s&nbsp;+%d</b></font><br/>`, utils.GetPropColor(prop.MModel.PropsColor), prop.MModel.Name, com.StrTo(items[0]).MustInt()+1)
			}
		}

		// 是否可交易
		trade := "不可交易"
		if (prop.CanTrade == 0 && prop.MModel.PropsLock > 0) || prop.CanTrade == 1 {
			trade = "可交易"
		}
		infoHtml += fmt.Sprintf(`<font color=%s>%s</font><br/>`, utils.EQGREENCOLOR, trade)
		infoHtml += expireHtml

		// 是否可强化
		plus := "不可强化"
		if prop.MModel.PlusFlag == 1 {
			plus = "可强化"
		}
		infoHtml += fmt.Sprintf(`<font color=%s>%s装备&nbsp(%s)</font><br/>`, utils.EQBASECOLOR, utils.GetZbPositionName(prop.MModel.Position), plus)

		// 主属性部分
		effectItems := strings.Split(prop.MModel.Effect, ":")
		infoHtml += fmt.Sprintf(`<font color=%s class="line">+%s %s %s</font><br/>`, utils.EQBASECOLOR, effectItems[1], utils.GetZbEffectName(effectItems[0]), strongValue)

		// 装备条件
		if prop.MModel.Requires != "" {
			if requires := strings.Split(prop.MModel.Requires, ","); len(requires) > 1 {
				if lvReq := strings.Split(requires[0], ":"); len(lvReq) > 1 {
					infoHtml += fmt.Sprintf(`<font color=%s>等级需求：%s级</font><br/>`, utils.EQBASECOLOR, lvReq[1])
				}
				if wxReq := strings.Split(requires[1], ":"); len(wxReq) > 1 {
					infoHtml += fmt.Sprintf(`<font color=%s>五行需求：%s系</font><br/>`, utils.EQBASECOLOR, utils.GetWxName(com.StrTo(wxReq[1]).MustInt()))
				}
			}
		}

		// 副属性部分
		plusEffectItems := strings.Split(prop.MModel.PlusEffect, ",")
		effectName := ""
		for _, effect := range plusEffectItems {
			if items := strings.Split(effect, ":"); len(items) > 1 {
				if effectName = utils.GetZbEffectName(items[0]); effectName != "" {
					infoHtml += fmt.Sprintf(`<font color=%s>+%s %s</font><br/>`, utils.EQPLUSCOLOR, items[1], effectName)
				} else if effectName = utils.GetZbSpecialEffectName(items[0]); effectName != "" {
					infoHtml += fmt.Sprintf(`<font color=%s>%s %s</font><br/>`, utils.EQPLUSCOLOR, effectName, items[1])
				} else {
					effectDep := ""
					switch items[0] {
					case "szmp":
						effectDep = "伤害的%s转化为MP"
						break
					case "sdmp":
						effectDep = "伤害的%s以MP抵消"
						break
					case "addmoney":
						effectDep = "战斗胜利获得金币增加%s点"
						break
					case "hitshp":
						effectDep = "偷取伤害的%s转化为生命"
						break
					case "hitsmp":
						effectDep = "偷取伤害的%s转化为魔法"
						break
					case "time":
						effectDep = "战斗等待时间减少%s秒"
						break
					}
					if effectDep != "" {
						infoHtml += fmt.Sprintf(`<font color=%s>`+effectDep+`</font><br/>`, utils.EQPLUSCOLOR, items[1])
					}
				}
			}
		}

		// 宝石镶嵌部分
		if prop.MModel.PlusNum == 0 {
			infoHtml += fmt.Sprintf(`<font color=%s>无卡槽</font><br/>`, utils.EQSPECIALCOLOR)
		} else if prop.MModel.PlusNum > 0 && prop.FHoleInfo != "" {
			effectItems = strings.Split(prop.FHoleInfo, ",")
			infoHtml += fmt.Sprintf(`<font color=%s>卡槽数：%d/%d</font><br/>`, utils.EQSPECIALCOLOR, len(effectItems), prop.MModel.PlusNum)
			for _, v := range effectItems {
				items := strings.Split(v, ":")
				effectDep := ""
				switch items[0] {
				case "ac":
					effectDep = "增加攻击%s"
					break
				case "mc":
					effectDep = "增加防御%s"
					break
				case "hits":
					effectDep = "增加命中%s"
					break
				case "miss":
					effectDep = "增加闪避%s"
					break
				case "hp":
					effectDep = "增加HP上限%s"
					break
				case "mp":
					effectDep = "增加MP上限%s"
					break
				case "speed":
					effectDep = "增加速度%s"
					break
				case "sdmp":
					effectDep = "将受到伤害的%s以MP抵消"
					break
				case "szmp":
					effectDep = "将受到伤害的%s转化为MP"
					break
				case "hitshp":
					effectDep = "命中吸取伤害的%s转化为自身HP"
					break
				case "hitsmp":
					effectDep = "命中吸取伤害的%s转化为自身MP"
					break
				case "dxsh":
					effectDep = "伤害抵销%s"
					break
				case "shjs":
					effectDep = "对敌人造成的伤害增加%s"
					break
				case "crit":
					effectDep = "会心一击率增加%s"
					break
				}
				if effectDep != "" {
					infoHtml += fmt.Sprintf(`<font color="red">宝石效果：`+effectDep+`</font><br/>`, items[1])
				}
			}
		} else {
			infoHtml += fmt.Sprintf(`<font color=%s>卡槽数：0/%d</font><br/>`, utils.EQSPECIALCOLOR, prop.MModel.PlusNum)
		}

		// 套装效果
		if prop.MModel.Series != "" && prop.MModel.Series != "0" {
			serieItems := strings.Split(prop.MModel.Series, ":")
			serieList := strings.Split(serieItems[1], "|")
			var petzbids []string
			if prop.Zbpets > 0 {
				for _, v := range ps.GetPZbs(prop.Zbpets) {
					v.GetM()
					if com.IsSliceContainsStr(serieList, strconv.Itoa(v.Pid)) {
						petzbids = append(petzbids, strconv.Itoa(v.Pid))
					}
				}
			}

			petZbStr := ""
			for _, v := range serieList {
				mprop := common.GetMProp(com.StrTo(v).MustInt())
				if com.IsSliceContainsStr(petzbids, strconv.Itoa(mprop.ID)) {
					petZbStr += fmt.Sprintf(`<font color=%s>%s</font><br/>`, utils.EQSPECIALCOLOR, mprop.Name)
				} else {
					petZbStr += fmt.Sprintf(`<font color=%s>%s</font><br/>`, utils.EQGREENCOLOR, mprop.Name)
				}
			}
			infoHtml += fmt.Sprintf(`<font color=%s>%s(%d/%d)</font><br/>`, utils.EQGLODCOLOR, serieItems[0], len(petzbids), len(serieList)) + petZbStr
			serieEffects := strings.Split(prop.MModel.SeriesEffect, ",")
			for i, v := range serieEffects {
				if v != "" && v != "0" {
					effectItems := strings.Split(v, ":")
					if len(effectItems) < 2 {
						continue
					}
					effectDep := ""
					effectName := ""
					if effectName = utils.GetZbEffectName(effectItems[0]); effectName != "" {
						effectDep = fmt.Sprintf(`+%s %s`, effectItems[1], effectName)
					} else if effectName = utils.GetZbSpecialEffectName(effectItems[0]); effectName != "" {
						effectDep = fmt.Sprintf(`%s %s`, effectName, effectItems[1])
					} else {
						switch effectItems[0] {
						case "szmp":
							effectDep = "伤害的%s转化为MP"
							break
						case "sdmp":
							effectDep = "伤害的%s以MP抵消"
							break
						case "addmoney":
							effectDep = "战斗胜利获得金币增加%s点"
							break
						case "hitshp":
							effectDep = "偷取伤害的%s转化为生命"
							break
						case "hitsmp":
							effectDep = "偷取伤害的%s转化为魔法"
							break
						case "time":
							effectDep = "战斗等待时间减少%s秒"
							break
						}
						if effectDep != "" {
							effectDep = fmt.Sprintf(effectDep, effectItems[1])
						}
					}
					if effectDep != "" {
						if i < len(petzbids) {
							infoHtml += fmt.Sprintf(`<font color=%s>(%d)套装：%s</font><br/>`, utils.EQSPECIALCOLOR, i+1, effectDep)
						} else {
							infoHtml += fmt.Sprintf(`<font color=%s>(%d)套装：%s</font><br/>`, utils.EQGREENCOLOR, i+1, effectDep)
						}
					}

				}
			}
		}
		break
	case 25:
		// 宝石
		infoHtml += "<font color='red'>"
		if prop.MModel.Requires != "" {
			infoHtml += "镶嵌部位:"
			for _, v := range strings.Split(prop.MModel.Requires, ",") {
				if items := strings.Split(v, ":"); len(items) > 1 {
					if items[0] == "postion" {

						for _, p := range strings.Split(items[1], "|") {
							infoHtml += utils.GetZbPositionName(com.StrTo(p).MustInt()) + " "
						}
					} else if items[0] == "color" {
						color2names := map[string]string{
							"2": "蓝色装备",
							"3": "紫色装备",
							"4": "绿色装备",
							"5": "黄色装备",
							"6": "橙色装备",
						}
						infoHtml += "只能镶嵌" + color2names[items[1]]
					}
				}
			}
		} else {
			infoHtml += "无需求"
		}
		infoHtml += "</font><br>"
		break
	default:
		break
	}
	// 道具说明
	infoHtml += fmt.Sprintf(`<font color=%s>%s</font><br/>`, utils.EQBASECOLOR, prop.MModel.Usages)

	headHmtl := `<table style="font-size:12px;" width=185 cellpadding=0 cellspacing=0 border=0>
					<tr> <td background=../images/ui/tips/border4_tl.gif width=5 height=5></td>
					<td background=../images/ui/tips/border4_t.gif></td>
					<td background=../images/ui/tips/border4_tr.gif></td>
					</tr>
					<tr><td width=5 background=../images/ui/tips/border4_l.gif></td>
					<td   style="background:#1F1F30;filter:Alpha(opacity=90);" align=center></td>
					<td width=5 background=../images/ui/tips/border4_r.gif></td></tr><tr><td width=5 background=../images/ui/tips/border4_l.gif></td>
					<td style="background:#1F1F30;filter:Alpha(opacity=90);">`
	footHtml := `</td><td width=5 background=../images/ui/tips/border4_r.gif></td>
					</tr><tr><td background=../images/ui/tips/border4_bl.gif width=5 height=5></td><td background=../images/ui/tips/border4_b.gif></td>
					<td background=../images/ui/tips/border4_br.gif></td>
					</tr>
					</table>`
	return headHmtl + infoHtml + footHtml
}

func (ps *PropService) CountZbAttr(zbs []models.UProp) map[string]float64 {
	zbAttr := make(map[string]float64)
	seriesCnts := make(map[string]*struct {
		effects string
		cnt     int
	})
	for _, v := range zbs {
		v.GetM()
		// 主属性与强化
		mEffect := strings.Split(v.MModel.Effect, ":")
		if len(mEffect) > 1 {
			zbAttr[mEffect[0]] += utils.ToFloat64(mEffect[1])
			//fmt.Printf("%s 主属性：%s\n", v.MModel.Name, mEffect)
			if plusEffect := strings.Split(v.PlusTmsEft, ","); len(plusEffect) > 1 {
				zbAttr[mEffect[0]] += utils.ToFloat64(plusEffect[1])
			}
		}

		// 附加属性
		bEffect := strings.Split(v.MModel.PlusEffect, ",")
		for _, eff := range bEffect {
			if eItems := strings.Split(eff, ":"); len(eItems) > 1 {
				//fmt.Printf("%s 副属性：%s\n", v.MModel.Name, eItems)
				zbAttr[eItems[0]] += utils.ToFloat64(eItems[1])
			}
		}

		// 水晶属性
		hEffect := strings.Split(v.FHoleInfo, ",")
		for _, eff := range hEffect {
			if eItems := strings.Split(eff, ":"); len(eItems) > 1 {
				switch eItems[0] {
				case "ac", "mc", "hp", "mp", "speed", "hits", "miss":
					zbAttr[eItems[0]+"rate"] += utils.ToFloat64(eItems[1])
				case "sdmp", "szmp", "dxsh", "hitshp", "hitsmp", "shjs", "crit":
					zbAttr[eItems[0]] += utils.ToFloat64(eItems[1])
				}
			}
		}

		// 套装数量统计
		if seriesItems := strings.Split(v.MModel.Series, ":"); len(seriesItems) > 1 {
			if ids := strings.Split(seriesItems[1], "|"); com.IsSliceContainsStr(ids, strconv.Itoa(v.Pid)) {
				se, ok := seriesCnts[seriesItems[1]]
				if ok {
					se.cnt += 1
				} else {
					seriesCnts[seriesItems[1]] = &struct {
						effects string
						cnt     int
					}{effects: v.MModel.SeriesEffect, cnt: 1}
				}
			}
		}

		//for k1, v1 := range zbAttr {
		//	fmt.Printf("%s,all attrs=>%s:%s\n", v.MModel.Name, k1, v1)
		//}
	}

	for _, se := range seriesCnts {
		sEffect := strings.Split(se.effects, ",")
		for i, eff := range sEffect {
			if i < se.cnt {
				//fmt.Printf("%s 套装属性：%s\n", sName, eff)
				if eItems := strings.Split(eff, ":"); len(eItems) > 1 {
					zbAttr[eItems[0]] += utils.ToFloat64(eItems[1])
				}
			}
		}
	}
	//for k, v := range zbAttr {
	//	fmt.Printf("all attrs=>%s:%s\n", k, v)
	//}
	return zbAttr
}

func (ps *PropService) CheckPropExpire() {
	ubs := []struct {
		Id     int
		Uid    int
		Pid    int
		Zbpets int
		Sums   int
		Bsum   int
		Psum   int
		Name   string
	}{}
	ps.GetDb().Table("userbag ub").Select(`ub.id, ub.uid, ub.pid, ub.zbpets, ub.sums, ub.bsum, ub.psum, p.name`).Joins("inner join props p on p.Id=ub.pid").Where("p.expire>0 and ub.stime+p.expire<=? and (ub.sums>0 or ub.bsum>0 or ub.psum>0)", utils.NowUnix()).Limit(50).Scan(&ubs)

	delIds := []int{}
	for _, prop := range ubs {
		delIds = append(delIds, prop.Id)
		//name = utils.ToUtf8(name)
		SelfGameLog(prop.Uid, fmt.Sprintf("物品到期：%s * %d", prop.Name, (prop.Sums+prop.Bsum+prop.Psum)), 222)
		if prop.Zbpets > 0 {
			ps.OptSvc.FightSrv.DelZbAttr(prop.Zbpets)
		}
	}
	if len(delIds) > 0 {
		ps.GetDb().Where("id in (?)", delIds).Delete(&models.UProp{})
	}
}

func (ps *PropService) CheckPropValid() {
	ps.GetDb().Where("sums=0 and bsum=0 and psum=0").Delete(&models.UProp{})
}

func (ps *PropService) OffZb(upId int) bool {
	return ps.GetDb().Model(&models.UProp{ID: upId}).Update(UpMap{"zbing": 0, "zbpets": 0}).RowsAffected > 0
}

func (ps *PropService) OffZbBypid(userId, petId, mpId int) bool {
	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	prop := ps.GetPropByPid(userId, mpId, true)
	if prop.Zbpets != petId {
		return false
	}

	if ps.GetDb().Model(prop).Update(UpMap{"zbing": 0, "zbpets": 0}).RowsAffected > 0 {
		if pet := ps.OptSvc.PetSrv.GetPet(userId, petId); pet != nil {
			newzbs := []string{}
			for _, v := range strings.Split(pet.Zb, ",") {
				if items := strings.Split(v, ":"); len(items) > 1 {
					if items[1] != strconv.Itoa(prop.ID) {
						newzbs = append(newzbs, v)
					}
				}
			}
			if ps.GetDb().Model(pet).Update("zb", strings.Join(newzbs, ",")).RowsAffected > 0 {
				ps.OptSvc.Commit()
				ps.OptSvc.FightSrv.DelZbAttr(pet.ID)
				return true
			}
		}
	}
	return false
}

func (ps *PropService) GetPmProps() ([]*models.UProp, []*models.UProp, []*models.UProp) {
	props := []models.UProp{}
	jbProps := []*models.UProp{}
	sjProps := []*models.UProp{}
	ybProps := []*models.UProp{}
	if ps.GetDb().Where("psum>0 and petime>?", ps.NowUnix()).Find(&props).RowsAffected > 0 {
		for _, p := range props {
			p1 := p
			p1.GetM()
			if p.Psell > 0 {
				jbProps = append(jbProps, &p1)
			} else if p.Psj > 0 {
				sjProps = append(sjProps, &p1)
			} else {
				ybProps = append(ybProps, &p1)
			}
		}
	}
	//selfProps := []*models.UProp{}

	return jbProps, sjProps, ybProps
}

func (ps *PropService) GetSelfPmProps(userId int) []*models.UProp {
	props := []models.UProp{}
	selfPmProps := []*models.UProp{}
	if ps.GetDb().Where("uid=? and psum>0", userId).Find(&props).RowsAffected > 0 {
		for _, p := range props {
			p1 := p
			p1.GetM()
			if p1.Petime < ps.NowUnix() {
				p1.PmTimeStr = "已过期"
			} else {
				p1.PmTimeStr = time.Unix(int64(p1.Petime), 0).Format("15:04:05")
			}
			if p1.Psell > 0 {
				p1.PmMoneyStr = strconv.Itoa(p1.Psell) + "金币"
			} else if p1.Psj > 0 {
				p1.PmMoneyStr = strconv.Itoa(p1.Psj) + "水晶"
			} else {
				p1.PmMoneyStr = strconv.Itoa(p1.Pyb) + "元宝"
			}
			selfPmProps = append(selfPmProps, &p1)
		}
	}
	return selfPmProps
}

func (ps *PropService) PutIn(userId, upId, num int) (result gin.H, msg string) {
	// num 小于 0 则放入全部
	prop := ps.GetProp(userId, upId, false)
	result = gin.H{}
	result["result"] = false
	if prop == nil {
		msg = "道具不存在!"
		return
	}
	if num > 0 && prop.Sums < num {
		msg = "存入道具数量过多！"
		return
	}
	if prop.Bsum == 0 {
		user := ps.OptSvc.UserSrv.GetUserById(userId)
		if user.BasePlace <= ps.GetCkPropCnt(userId) {
			msg = "仓库已满！"
			return
		}
	}
	msg = "存入道具数量过多！"
	if num < 0 {
		if ps.GetDb().Model(prop).Update(gin.H{"bsum": gorm.Expr("bsum+sums"), "sums": 0}).RowsAffected > 0 {
			result["result"] = true
			msg = "放入成功！"
		}
	} else if ps.GetDb().Model(prop).Where("sums >= ?", num).Update(gin.H{"sums": gorm.Expr("sums-?", num), "bsum": gorm.Expr("bsum+?", num)}).RowsAffected > 0 {
		result["result"] = true
		msg = "放入成功！"
	}
	if result["result"].(bool) {
		prop := ps.GetProp(userId, upId, false)
		result["bsum"] = prop.Bsum
		result["sum"] = prop.Sums
		result["id"] = prop.ID
	}
	return
}

func (ps *PropService) PutOut(userId, upId, num int, inputPwd string) (result gin.H, msg string) {
	// num 小于 0 则取出全部
	prop := ps.GetProp(userId, upId, false)
	result = gin.H{}
	result["result"] = false
	if prop == nil {
		msg = "道具不存在!"
		return
	}
	if num > 0 && prop.Bsum < num {
		msg = "取出道具数量过多！"
		return
	}
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	need_pass := ps.OptSvc.UserSrv.CheckNeedPwd(inputPwd, user.CkPwd)

	result["need_pass"] = need_pass
	if need_pass {
		msg = "仓库密码错误！"
		return
	}
	if prop.Sums == 0 {
		if user.BagPlace <= ps.GetCarryPropsCnt(userId) {
			msg = "背包已满！"
			return
		}
	}
	msg = "取出道具数量过多！"
	if num < 0 {
		if ps.GetDb().Model(prop).Update(gin.H{"sums": gorm.Expr("bsum+sums"), "bsum": 0}).RowsAffected > 0 {
			result["result"] = true
			msg = "取出成功！"
		}
	}
	if ps.GetDb().Model(prop).Where("bsum >= ?", num).Update(gin.H{"bsum": gorm.Expr("bsum-?", num), "sums": gorm.Expr("sums+?", num)}).RowsAffected > 0 {
		result["result"] = true
		msg = "取出成功！"
	}
	if result["result"].(bool) {
		prop := ps.GetProp(userId, upId, false)
		result["bsum"] = prop.Bsum
		result["sum"] = prop.Sums
		result["id"] = prop.ID
	}
	return
}

func (ps *PropService) Throw(userId, upId int) bool {
	prop := ps.GetProp(userId, upId, false)
	if prop == nil {
		return false
	}
	if ps.GetDb().Model(prop).Where("sums >= 0").Update(gin.H{"sums": 0}).RowsAffected > 0 {
		return true
	}
	return false
}

func (ps *PropService) Auction(userId, upId, num, price int, currency, nickname string) (bool, string) {
	prop := ps.GetProp(userId, upId, false)
	if prop.Psum > 0 || prop.Psell > 0 || prop.Psj > 0 || prop.Pyb > 0 {
		return false, "该物品已在拍卖所中！"
	}
	if prop.Sums < num {
		return false, "拍卖数量过多！"
	}
	if prop.Zbing > 0 {
		return false, "道具不存在！"
	}
	pmProps := ps.GetSelfPmProps(userId)
	count := 0
	currencyName := "金币"
	for _, p := range pmProps {
		switch currency {
		case "jb":
			if p.Psell > 0 {
				count++
			}
			break
		case "sj":
			if p.Psj > 0 {
				count++
			}
		case "yb":
			if p.Pyb > 0 {
				count++
			}
		}

	}
	if count >= 4 {
		return false, "拍卖所拍卖道具数量已达上限"
	}
	if !prop.AbleTrade() {
		return false, "该道具不可交易！"
	}
	nowUnix := ps.NowUnix()
	if currency == "jb" {
		currency = "sell"
		currencyName = "金币"
	} else if currency == "sj" {
		if price < 10 {
			return false, "最低拍卖价格 10 水晶起！"
		}
		currencyName = "水晶"
	} else if currency == "yb" {
		if price < 10 {
			return false, "最低拍卖价格 10 元宝起！"
		}
		currencyName = "元宝"
	}
	nicknameCrc := 0
	if nickname != "" {
		nicknameCrc = utils.CRC32(nickname)
	}

	if ps.GetDb().Model(prop).Update(gin.H{
		"p" + currency: price,
		"pstime":       nowUnix,
		"petime":       nowUnix + 10800,
		"sums":         gorm.Expr("sums - ?", num),
		"psum":         gorm.Expr("psum + ?", num),
		"buycode":      nicknameCrc,
	}).RowsAffected > 0 {
		SelfGameLog(userId, fmt.Sprintf("拍卖道具：%s * %d，价格为 %d %s， 指定收货人：%s", prop.MModel.Name, num, price, currencyName, nickname), 155)
		return true, "拍卖成功！"
	} else {
		return false, "拍卖道具不存在！"
	}
}

func (ps *PropService) ReAuction(userId, id int) (bool, string) {
	prop := ps.GetProp(userId, id, false)
	if prop == nil {
		return false, "道具不存在！"
	}
	if prop.Petime >= utils.NowUnix() {
		return false, "道具正在拍卖中！"
	}
	nowUnix := ps.NowUnix()
	if ps.GetDb().Model(prop).Update(gin.H{
		"pstime": nowUnix,
		"petime": nowUnix + 10800,
	}).RowsAffected > 0 {
		return true, "续拍成功！"
	} else {
		return false, "道具不存在！"
	}
}

func (ps *PropService) BagLeftPlace(userId int) int {
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	return user.BagPlace - ps.GetCarryPropsCnt(userId)
}

func (ps *PropService) RollAuction(userId, id int) (bool, string) {
	prop := ps.GetProp(userId, id, false)
	if prop == nil {
		return false, "道具不存在！"
	}
	if prop.Psum == 0 {
		return false, "道具没有在拍卖中！"
	}
	if prop.Sums == 0 && ps.BagLeftPlace(userId) < 1 {
		return false, "背包空间不足，取回失败！"
	}
	if ps.GetDb().Model(prop).Update(gin.H{
		"sums":   gorm.Expr("sums + psum"),
		"psum":   0,
		"pstime": 0,
		"petime": 0,
		"psell":  0,
		"psj":    0,
		"pyb":    0,
	}).RowsAffected > 0 {
		return true, "取回成功！"
	} else {
		return false, "道具不存在！"
	}
}

func (ps *PropService) Purchase(userId, id, num int) (bool, string) {
	prop := ps.GetOtherPropById(id, false)
	if prop == nil || prop.Psum == 0 || prop.Psum < num || prop.Petime < utils.NowUnix() {
		return false, "您购买的数量太多！"
	}
	if prop.Uid == userId {
		return false, "不能购买自己的东西！"
	}
	now := utils.NowUnix()
	if rcache.PmsTimer.InCoolTime(userId, now) {
		return false, "操作频繁！"
	}
	currencyName := "jb"
	var payAll int
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	if prop.Psell > 0 {
		currencyName = "金币"
		payAll = prop.Psell * num
		if user.Money < payAll || ps.GetDb().Model(user).Where("money > ?", payAll).Update(gin.H{"money": gorm.Expr("money-?", payAll)}).RowsAffected == 0 {
			return false, "您的金币不足！"
		}
	} else if prop.Psj > 0 {
		currencyName = "水晶"
		payAll = prop.Psj * num
		userInfo := ps.OptSvc.UserSrv.GetUserInfoById(userId)
		if userInfo.Sj < payAll || ps.GetDb().Model(userInfo).Where("sj > ?", payAll).Update(gin.H{"sj": gorm.Expr("sj-?", payAll)}).RowsAffected == 0 {
			return false, "您的水晶不足！"
		}
	} else if prop.Pyb > 0 {
		currencyName = "元宝"
		payAll = prop.Pyb * num
		if user.Yb < payAll || ps.GetDb().Model(user).Where("yb > ?", payAll).Update(gin.H{"yb": gorm.Expr("yb-?", payAll)}).RowsAffected == 0 {
			return false, "您的金币不足！"
		}
	} else {
		return false, "您购买的数量太多！"
	}
	useCode := "否"
	if prop.BuyCode != 0 {
		useCode = "是"
	}
	if prop.BuyCode != 0 && utils.CRC32(user.Nickname) != prop.BuyCode {
		return false, "该物品不是卖给您的！"
	}
	note := "购买物品:"
	rcache.PmsTimer.Set(userId, now)
	if prop.MModel.Vary == 1 {
		myProp := ps.GetPropByPid(userId, prop.Pid, false)
		if myProp != nil {
			if ps.GetDb().Model(myProp).Update(gin.H{
				"sums": gorm.Expr("sums + ?", num),
			}).RowsAffected == 0 {
				return false, "服务器错误！"
			}
		} else {
			prop := &models.UProp{
				Pid:   prop.Pid,
				Uid:   userId,
				Sums:  num,
				Sell:  prop.MModel.SellJb,
				Stime: utils.NowUnix(),
			}
			if ps.GetDb().Create(prop).Error != nil {
				return false, "服务器错误！"
			}
		}
		if ps.GetDb().Model(prop).Update(gin.H{
			"psum":    0,
			"pstime":  0,
			"petime":  0,
			"psell":   0,
			"psj":     0,
			"pyb":     0,
			"buycode": 0,
		}).RowsAffected == 0 {
			return false, "服务器错误！"
		}
		note += fmt.Sprintf("获得 %s*%d，指定收货：%s，共花费：%d %s", prop.MModel.Name, num, useCode, payAll, currencyName)
	} else {
		if ps.GetDb().Model(prop).Update(gin.H{
			"uid": userId,
		}).RowsAffected == 0 {

			return false, "服务器错误！"
		}
		note += fmt.Sprintf("获得 %s*%d，道具id:%d，强化：%s,镶嵌：%s，指定收货：%s，共花费：%d %s", prop.MModel.Name, num, prop.ID, prop.PlusTmsEft, prop.FHoleInfo, useCode, payAll, currencyName)
	}
	ps.OptSvc.Commit()
	GameLog(prop.Uid, userId, note, 102)
	return true, fmt.Sprintf("购买成功！共花费：%d %s", payAll, currencyName)
}

func (ps *PropService) ShopSell(userId, id, num int) (bool, string) {
	prop := ps.GetProp(userId, id, false)
	if prop == nil {
		return false, "道具不存在！"
	}
	if prop.Sums < num {
		return false, "道具数量不足！"
	}
	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	if ps.GetDb().Model(prop).Where("sums >= ?", num).Update(gin.H{"sums": gorm.Expr("sums-?", num)}).RowsAffected == 0 {
		return false, "道具数量不足！"
	}
	user := ps.OptSvc.UserSrv.GetUserById(userId)
	prop.GetM()
	newMoney := user.Money + prop.MModel.SellJb*num
	if newMoney >= 1000000000 {
		newMoney = 1000000000
	}
	ps.GetDb().Model(user).Update(gin.H{"money": newMoney})
	ps.OptSvc.Commit()
	return true, "出售成功！"
}

func (ps *PropService) ShopPurchase(userId, pid, num int, currency string) (bool, string) {
	mprop := common.GetMProp(pid)
	if mprop == nil {
		return false, "道具不存在！"
	}
	if num > 100 {
		return false, "单次购买数量不能超过100！"
	}
	ps.OptSvc.Begin()
	defer ps.OptSvc.Rollback()
	if currency == "jb" {
		if mprop.BuyJb == 0 || mprop.BuyYb > 0 {
			return false, "道具不存在！"
		}
		payAll := mprop.BuyJb * num
		if ps.GetDb().Model(&models.User{ID: userId}).Where("money >= ?", payAll).Update(gin.H{"money": gorm.Expr("money-?", payAll)}).RowsAffected == 0 {
			return false, "您的金币不足！"
		}
	} else if currency == "sj" {
		if mprop.BuySj == 0 || mprop.BuyYb == 99999 || mprop.Stime <= 0 {
			return false, "道具不存在！"
		}
		payAll := mprop.BuySj * num
		if ps.GetDb().Model(&models.UserInfo{}).Where("uid = ? and sj >= ?", userId, payAll).Update(gin.H{"sj": gorm.Expr("sj-?", payAll)}).RowsAffected == 0 {
			return false, "您的水晶不足！"
		}
	} else if currency == "yb" {
		if mprop.BuyYb == 0 || mprop.BuyYb == 99999 || mprop.Stime <= 0 {
			return false, "道具不存在！"
		}
		payAll := mprop.BuyYb * num
		user := ps.OptSvc.UserSrv.GetUserById(userId)
		addVip := (user.Useyb + payAll) / 100
		newUseyb := (user.Useyb + payAll) % 100
		if ps.GetDb().Model(user).Where("yb >= ?", payAll).Update(gin.H{"yb": gorm.Expr("yb-?", payAll), "vip": gorm.Expr("vip+?", addVip), "score": gorm.Expr("score+?", addVip), "useyb": newUseyb}).RowsAffected == 0 {
			return false, "您的元宝不足！"
		}
		ps.GetDb().Create(&models.YbLog{Pid: pid, Account: user.Account, UseYb: payAll, Btime: ps.NowUnix(), Pnote: fmt.Sprintf("%s, 购买时剩余元宝：%d", mprop.Name, user.Yb-payAll), Num: num})
	} else if currency == "vip" {
		if mprop.Vip == 0 || mprop.Vip == 99999 || mprop.Stime <= 0 {
			return false, "道具不存在！"
		}
		payAll := mprop.Vip * num
		user := ps.OptSvc.UserSrv.GetUserById(userId)
		if ps.GetDb().Model(user).Where("vip >= ?", payAll).Update(gin.H{"vip": gorm.Expr("vip-?", payAll)}).RowsAffected == 0 {
			return false, "您的VIP不足！"
		}
	} else if currency == "ww" {
		if mprop.Prestige == 0 {
			return false, "道具不存在！"
		}
		payAll := mprop.Prestige * num
		user := ps.OptSvc.UserSrv.GetUserById(userId)
		if ps.GetDb().Model(user).Where("prestige >= ?", payAll).Update(gin.H{"prestige": gorm.Expr("prestige-?", payAll)}).RowsAffected == 0 {
			return false, "您的威望不足！"
		}
	} else if currency == "zkyb" {
		timeSetting := common.GetWelcome("timelimitbuy")
		if timeSetting == nil {
			return false, "活动未开启！"
		}
		timeStr := strings.SplitAfterN(timeSetting.Text, "|", 2)
		startTime := utils.StrParseMustTime(timeStr[0])
		endTime := utils.StrParseMustTime(timeStr[1])
		now := time.Now()
		if startTime.Sub(now).Seconds() > 0 {
			return false, "活动未开启！"
		}
		lefttime := endTime.Sub(now).Seconds()
		if lefttime <= 0 {
			return false, "活动未开启！"
		}
		if num != 1 {
			return false, "每次最多只可以购买1个！"
		}
		buyFlag := false
		for _, goodStr := range strings.Split(timeSetting.Content, ",") {
			if items := strings.Split(goodStr, ":"); len(items) == 2 {
				if pid == com.StrTo(items[0]).MustInt() {
					maxNum := com.StrTo(items[1]).MustInt()
					leftNum := maxNum
					zkNum, err := rcache.GetZKGoodNum(pid)
					if err == nil {
						leftNum = maxNum - zkNum
					}
					if leftNum <= 0 {
						return false, "该道具已被抢购一空！"
					}
					if num > leftNum {
						return false, "您所购买的数量过多！"
					}

					payAll := mprop.ZhekouYb * num
					user := ps.OptSvc.UserSrv.GetUserById(userId)

					addVip := (user.Useyb + payAll) / 100
					newUseyb := (user.Useyb + payAll) % 100
					if ps.GetDb().Model(user).Where("yb >= ?", payAll).Update(gin.H{"yb": gorm.Expr("yb-?", payAll), "vip": gorm.Expr("vip+?", addVip), "score": gorm.Expr("score+?", addVip), "useyb": newUseyb}).RowsAffected == 0 {
						return false, "您的元宝不足！"
					}

					rcache.SetZKGoodNum(pid, leftNum-num)
					ps.GetDb().Create(&models.YbLog{Pid: pid, Account: user.Account, UseYb: payAll, Btime: ps.NowUnix(), Pnote: fmt.Sprintf("%s, 购买时剩余元宝：%d", mprop.Name, user.Yb-payAll), Num: num})
					buyFlag = true
					break
				}
			}
		}
		if !buyFlag {
			return false, "没有该道具！"
		}
	} else {
		return false, "参数错误！"
	}
	if !ps.AddProp(userId, pid, num, true) {
		return false, "背包空间不足！"
	}
	ps.OptSvc.Commit()
	return true, "购买成功！"
}

func (ps *PropService) DropPetZb(uid, petId int) {
	ps.GetDb().Where("uid = ? and zbpets = ?", uid, petId).Delete(&models.UProp{})
}

func (ps *PropService) GetSmShopGood(update bool) gin.H {
	cacheData := gin.H{}
	if !update {

		if cacheData, err := rcache.GetSmShopGood(); err == nil {
			return cacheData
		}
	}

	ybGoods, sjGoods, vipGoods := make(map[string][]gin.H), make(map[string][]gin.H), make(map[string][]gin.H)
	category := []string{"remai", "jinhua", "chongwu", "zhuangbei"}
	for _, c := range category {
		ybGoods[c] = []gin.H{}
		sjGoods[c] = []gin.H{}
		vipGoods[c] = []gin.H{}
	}
	shopGoods := []models.MProp{}
	ps.GetDb().Where("(yb>0 or sj>0 or vip>0) and stime>0").Order("stime").Find(&shopGoods)
	for i := 0; i < len(shopGoods); i++ {
		good := &shopGoods[i]
		goodInfo := gin.H{
			"id":        good.ID,
			"name":      good.Name,
			"vary_id":   good.VaryName,
			"overlying": good.Vary == 1,
		}
		for j, c := range category {
			if strconv.Itoa(good.Stime)[:1] == strconv.Itoa(j+1) {
				if good.BuyYb > 0 {
					goodInfo["price"] = good.BuyYb
					ybGoods[c] = append(ybGoods[c], goodInfo)
				} else if good.BuySj > 0 {
					goodInfo["price"] = good.BuySj
					sjGoods[c] = append(sjGoods[c], goodInfo)
				} else {
					goodInfo["price"] = good.Vip
					vipGoods[c] = append(vipGoods[c], goodInfo)
				}
			}
		}
	}
	cacheData = gin.H{
		"yb_goods":  ybGoods,
		"sj_goods":  sjGoods,
		"vip_goods": vipGoods,
	}
	rcache.SetSmShopGood(cacheData)
	return cacheData

}

func (ps *PropService) GetSmShopQgList() gin.H {
	timeSetting := common.GetWelcome("timelimitbuy")
	if timeSetting == nil {
		return gin.H{
			"lefttime": 0,
		}
	}
	timeStr := strings.SplitAfterN(timeSetting.Text, "|", 2)
	startTime := utils.StrParseMustTime(timeStr[0])
	endTime := utils.StrParseMustTime(timeStr[1])
	now := time.Now()
	if startTime.Sub(now).Seconds() > 0 {
		return gin.H{
			"lefttime": 0,
		}
	}
	lefttime := endTime.Sub(now).Seconds()
	if lefttime <= 0 {
		return gin.H{
			"lefttime": 0,
		}
	}
	zkList := [][]int{}
	for _, goodStr := range strings.Split(timeSetting.Content, ",") {
		if items := strings.Split(goodStr, ":"); len(items) == 2 {
			zkList = append(zkList, []int{com.StrTo(items[0]).MustInt(), com.StrTo(items[1]).MustInt()})
		}

	}
	ybGoods := []gin.H{}
	zkNumList, err := rcache.GetZKGoodNumList()
	if err == nil {
		for pid, num := range zkNumList {
			for i := 0; i < len(zkList); i++ {
				if pid == zkList[i][0] {
					zkList[i][1] = zkList[i][1] - num
				}
			}

		}
	}
	for i := 0; i < len(zkList); i++ {
		pid := zkList[i][0]
		num := zkList[i][1]
		good := common.GetMProp(pid)
		if good != nil {
			ybGoods = append(ybGoods, gin.H{
				"id":        good.ID,
				"name":      good.Name,
				"price":     good.ZhekouYb,
				"vary_id":   good.VaryName,
				"num":       num,
				"overlying": good.Vary == 1,
			})
			fmt.Printf("id:%d, name:%s\n", good.ID, good.Name)
		}
	}
	return gin.H{
		"lefttime": lefttime,
		"qg_list":  ybGoods,
	}
}

func (ps *PropService) GetDjShopGood(update bool) gin.H {

	shopData := gin.H{}
	if !update {
		if shopData, err := rcache.GetDjShopGood(); err == nil {
			return shopData
		}
	}
	jbProp := []models.MProp{}
	ps.GetDb().Where("buy>0 and yb=0 and varyname!=9").Find(&jbProp)
	jbData := []gin.H{}
	for _, prop := range jbProp {
		jbData = append(jbData, gin.H{
			"id":        prop.ID,
			"name":      prop.Name,
			"price":     prop.BuyJb,
			"vary_id":   prop.VaryName,
			"overlying": prop.IsVary(),
		})
	}
	wwProp := []struct {
		Id       int
		Price    int
		Downtime int
		Selltype string
		Varyname int
		Name     string
		Vary     int
	}{}
	sql := `select sl.id as id,
	sl.price as price,
	sl.down_time as downtime,
	sl.sell_type as selltype,
	p.varyname as varyname,
	p.name as name,
	p.vary as vary
	FROM shop_list sl left join props p on sl.props_id=p.id where sl.sell_type='prestige' and sl.price>0 order by sl.rank desc`
	ps.GetDb().Raw(sql).Scan(&wwProp)
	wwData := []gin.H{}
	for _, prop := range wwProp {
		wwData = append(wwData, gin.H{
			"id":        prop.Id,
			"name":      prop.Name,
			"price":     prop.Price,
			"vary_id":   prop.Varyname,
			"overlying": prop.Vary == 1,
		})
	}
	shopData["jbprops"] = jbData
	shopData["wwprops"] = wwData
	rcache.SetDjShopGood(shopData)
	return shopData
}

func (ps *PropService) GetCkPropData(userId int) []gin.H {
	ckProps := []models.UProp{}
	ps.GetDb().Where("bsum>0 and uid=?", userId).Find(&ckProps)
	ckData := []gin.H{}
	for _, prop := range ckProps {
		prop.GetM()
		ckData = append(ckData, gin.H{
			"id":        prop.ID,
			"name":      prop.MModel.Name,
			"price":     prop.MModel.SellJb,
			"vary_id":   prop.MModel.VaryName,
			"num":       prop.Bsum,
			"overlying": prop.MModel.IsVary(),
		})
	}
	return ckData
}

func (ps *PropService) GetCkPropCnt(userId int) int {
	cnt := 0
	ps.GetDb().Model(&models.UProp{}).Where("bsum>0 and uid=?", userId).Count(&cnt)
	return cnt
}

func (ps *PropService) AddPropRand(user *models.User, carryPropCnt int, ratePids []*models.RatePid, IsOpenMap bool) (propNameList []string) {

	pids := []int{}
	for _, ratePid := range ratePids {
		rate := ratePid.Rate
		if !IsOpenMap {
			rate *= NoKeyMapPunishmentRate
		}
		if rand.Intn(ratePid.Rate) == 0 {

			pids = append(pids, ratePid.Pid)
		}
	}
	if len(pids) == 0 {
		return
	}
	propNameList = ps.AddPropByPids(user, carryPropCnt, pids)
	return
}

func (ps *PropService) AddPropByPids(user *models.User, carryPropCnt int, pids []int) (propNameList []string) {
	var findHas bool        // 循环中使用，减少重复声明
	var findInbb bool       // 循环中使用，减少重复声明
	var mprop *models.MProp // 循环中使用，减少重复声明
	hasProps := []*models.UProp{}
	ps.GetDb().Where("uid=? and pid in (?)", user.ID, pids).Find(&hasProps)
	hasVaryMap := make(map[int]int)
	noHasVaryMap := make(map[int]int)
	noVaryMap := make(map[int]int)
	for _, pid := range pids {
		// 添加道具
		// 先检查是否可叠加
		// 若可叠加则检查背包是否有这个道具
		// 若不可叠加或背包没有这个道具，则检查背包是否上限了
		mprop = common.GetMProp(pid)
		if mprop == nil {
			continue
		}
		if carryPropCnt < user.BagPlace {
			propNameList = append(propNameList, mprop.Name)
			if mprop.Vary == 2 {
				// 不可叠加
				carryPropCnt += 1
				noVaryMap[pid] += 1
			} else {
				findHas = false
				findInbb = false
				for _, prop := range hasProps {
					if pid == prop.Pid {
						findHas = true
						if prop.Sums > 0 {
							findInbb = true
						}
					}
				}
				if !findInbb {
					carryPropCnt += 1
				}
				if !findHas {
					noHasVaryMap[pid] += 1
				} else {
					hasVaryMap[pid] += 1
				}
			}
		} else {
			break
		}
	}

	go func() {
		db := persistence.GetOrm().Begin()
		defer db.Commit()
		num2Pids := make(map[int][]int)
		for id, num := range hasVaryMap {
			ls := num2Pids[num]
			ls = append(ls, id)
			num2Pids[num] = ls
		}

		// 可叠加、身上存在该道具的处理
		for num, ids := range num2Pids {
			db.Model(&models.UProp{}).Where("uid=? and pid in (?)", user.ID, ids).Update(gin.H{"sums": gorm.Expr("sums+?", num)})
		}

		insertItems := []string{}
		insertInterfaces := []interface{}{}

		now := utils.NowUnix()
		// 可叠加、身上不存在该道具的处理
		for id, num := range noHasVaryMap {
			insertItems = append(insertItems, "(?,?,?,?,?,?)")
			items := []interface{}{user.ID, id, 0, 1, num, now}
			insertInterfaces = append(insertInterfaces, items...)
		}

		// 不可叠加道具的处理
		for id, num := range noVaryMap {
			insertItems = append(insertItems, "(?,?,?,?,?,?)")
			items := []interface{}{user.ID, id, 0, 2, 1, now}
			for i := 0; i < num; i++ {
				insertInterfaces = append(insertInterfaces, items...)
			}
		}
		if len(insertItems) > 0 {
			db.Exec(fmt.Sprintf("insert userbag(uid,pid,sell,vary,sums,stime) values%s", strings.Join(insertItems, ",")), insertInterfaces...)
		}
	}()
	return
}
