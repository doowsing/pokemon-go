package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
	"pokemon/common/rcache"
	"pokemon/game/models"
	"pokemon/game/repositories"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"time"
)

type UserService struct {
	BaseService
	repo *repositories.UserRepository
}

func NewUserService(osrc *OptService) *UserService {
	us := &UserService{repo: repositories.NewUserRepository()}
	us.SetOptSrc(osrc)
	return us
}

func (us *UserService) Logout(id int) bool {
	rcache.DelIdToken(id)
	return true
}

func (us *UserService) InitCreate() bool {
	user := models.User{
		Account:   "GGNNN",
		PasswdMd5: utils.Sha1("GGNNN321"),
		Nickname:  "GGNNN",
		Regtime:   int(time.Now().Unix()),
		BagPlace:  30,
		BasePlace: 40,
		McPlace:   20,
	}
	err := us.repo.Insert(&user)
	if err != nil {
		return false
	}
	return true
}

func (us *UserService) InitTable() error {
	return us.repo.InitTable()
}

// 2019 11 08 新写入
func (us *UserService) Login(account, passwd string) int {
	// 返回结果：
	// 密码错误：-1
	// 正确：0
	// 被封号：1

	if user, ok := us.repo.GetByUserNameAndPwd(account, utils.Md5(passwd)); !ok {
		return -1
	} else if user.Password != "" && user.Password != "0" {
		return 1
	} else {
		return 0
	}
}

func (us *UserService) GetByUserNameAndPwd(account, passwd string) *models.User {
	user := &models.User{}
	if result := us.GetDb().Where("Name = ? and secret = ?", account, utils.Md5(passwd)).First(&user); result.RowsAffected == 0 {
		return nil
	} else {
		return user
	}
}

func (us *UserService) GetUserById(id int) *models.User {
	user := &models.User{}
	us.GetDb()
	if result := us.GetDb().Where("Id=?", id).First(user); result.RowsAffected > 0 {
		return user
	} else {
		return nil
	}
}

func (us *UserService) GetUserInfoById(id int) *models.UserInfo {
	userInfo := &models.UserInfo{}
	if result := us.GetDb().Where("uid=?", id).First(userInfo); result.RowsAffected > 0 {
		return userInfo
	} else {
		return nil
	}
}

func (us *UserService) UpdateUserInfo(userId int, where ...interface{}) func(attrs ...interface{}) *gorm.DB {
	opdb := us.GetDb().Model(&models.UserInfo{ID: userId})
	if len(where) > 0 {
		opdb = opdb.Where(where[0], where[1:]...)
	}
	return opdb.Update
}

func (us *UserService) UpdateUser(userId int, where ...interface{}) func(attrs ...interface{}) *gorm.DB {
	opdb := us.GetDb().Model(&models.User{ID: userId})
	if len(where) > 0 {
		opdb = opdb.Where(where[0], where[1:]...)
	}
	return opdb.Update
}

func (us *UserService) SetMBid(userId, petid int) bool {
	return us.UpdateUser(userId)(UpMap{"mbid": petid, "fightbb": petid}).RowsAffected > 0
}

func (us *UserService) AddTgPlace(userId int) bool {
	return us.UpdateUser(userId)(UpMap{"tgmax": gorm.Expr("tgmax+1")}).RowsAffected > 0
}

func (us *UserService) AddShowTimes(userId, num int) bool {
	return us.UpdateUserInfo(userId)(UpMap{"tgmax": gorm.Expr("tgmax+?", num)}).RowsAffected > 0
}

func (us *UserService) AddSj(userId, num int) bool {
	return us.UpdateUserInfo(userId)(UpMap{"sj": gorm.Expr("sj+?", num)}).RowsAffected > 0
}

// 扣取金币
func (us *UserService) DecreaseJb(userId, num int) bool {
	return us.GetDb().Model(&models.User{ID: userId}).Where("money >= ?", num).Update(UpMap{"money": gorm.Expr("money - ?", num)}).RowsAffected > 0
}

// 增加金币
func (us *UserService) IncreaseJb(userId, num int) bool {
	return us.GetDb().Model(&models.User{ID: userId}).Update(UpMap{"money": gorm.Expr("money + ?", num)}).RowsAffected > 0
}

// 扣取水晶
func (us *UserService) DecreaseSj(userId, num int) bool {
	return us.GetDb().Model(&models.UserInfo{ID: userId}).Where("sj >= ?", num).Update(UpMap{"sj": gorm.Expr("sj - ?", num)}).RowsAffected > 0
}

func (us *UserService) AddYb(userId, num int) bool {
	return us.UpdateUser(userId)(UpMap{"yb": gorm.Expr("yb+?", num)}).RowsAffected > 0
}

func (us *UserService) AddBagPlace(userId, num, maxnum int) bool {
	return us.UpdateUser(userId, "maxnum<?", maxnum)(UpMap{"maxnum": gorm.Expr("maxnum+?", num)}).RowsAffected > 0
}

func (us *UserService) AddCkPlace(userId, num, maxnum int) bool {
	return us.UpdateUser(userId, "maxbase<?", maxnum)(UpMap{"maxbase": gorm.Expr("maxbase+?", num)}).RowsAffected > 0
}

func (us *UserService) AddMcPlace(userId, num, maxnum int) bool {
	return us.UpdateUser(userId, "maxmc<?", maxnum)(UpMap{"maxmc": gorm.Expr("maxmc+?", num)}).RowsAffected > 0
}

func (us *UserService) AddTgTime(userId, num int) bool {
	return us.UpdateUser(userId)(UpMap{"tgtime": gorm.Expr("tgtime+?", num)}).RowsAffected > 0
}

func (us *UserService) SetMaps(userId int, openMaps string) bool {
	return us.UpdateUser(userId)(UpMap{"openmap": openMaps}).RowsAffected > 0
}

func (us *UserService) AddJbAuto(userId, num int) bool {
	return us.UpdateUser(userId)(UpMap{"sysautosum": gorm.Expr("sysautosum+?", num)}).RowsAffected > 0
}

func (us *UserService) AddYbAuto(userId, num int) bool {
	return us.UpdateUser(userId)(UpMap{"maxautofitsum": gorm.Expr("maxautofitsum+?", num)}).RowsAffected > 0
}

func (us *UserService) AddTeamAuto(userId, num int) bool {
	return us.UpdateUserInfo(userId)(UpMap{"team_auto_times": gorm.Expr("team_auto_times+?", num)}).RowsAffected > 0
}

func (us *UserService) AddSsCzl(userId, num int) bool {
	return us.UpdateUserInfo(userId)(UpMap{"czl_ss": gorm.Expr("czl_ss+?", num)}).RowsAffected > 0
}

func (us *UserService) AddZcScore(userId int, buffStatus string) bool {
	// 使女神要塞双倍积分
	return us.UpdateUserInfo(userId)(UpMap{"buff_status": fmt.Sprintf("add_zc_jifen:%s,%s;", time.Now().Format("Ymd"), buffStatus)}).RowsAffected > 0
}

func (us *UserService) GetCardTile(codeName string) *models.CardTitle {
	return common.GetCardTitle(codeName)
}

func (us *UserService) GetDatiPlayer(userId int) *models.AoyunPlayer {
	player := &models.AoyunPlayer{}
	us.GetDb().Where("uid = ?", userId).First(player)
	if player.Id > 0 {
		return player
	}
	return nil
}

// 是否已扫过雷
func (us *UserService) GetSaoleiHasRecord(userId int) bool {
	flag, err := rcache.GetSaoleiTodayUser(userId)
	if err == nil && flag == 1 {
		return true
	}
	return false
}

//  设置为已扫过雷
func (us *UserService) SetSaoleiStatus(userId int) {
	rcache.SetSaoleiTodayUser(userId, 1)
}

// 是否可以扫雷
func (us *UserService) GetSaoleiStatus(userId int) (int, bool) {
	userInfo := us.GetUserInfoById(userId)
	useTicket := us.GetSaoleiTicketRecord(userId)
	if us.GetSaoleiHasRecord(userId) {
		if useTicket {
			return userInfo.FSaoleiPoints, true
		} else {
			return userInfo.FSaoleiPoints, false
		}
	}
	user := us.GetUserById(userId)
	mainPet := us.OptSvc.PetSrv.GetPetById(user.Mbid)
	if mainPet != nil && com.StrTo(mainPet.Czl).MustFloat64() > 65 {
		return userInfo.FSaoleiPoints, true
	}
	if useTicket {
		return userInfo.FSaoleiPoints, true
	}
	return userInfo.FSaoleiPoints, false
}

// 使用闯关卡进入的记录
func (us *UserService) GetSaoleiTicketRecord(userId int) bool {
	flag, err := rcache.GetSaoleiTicketUser(userId)
	if err == nil && flag == 1 {
		return true
	}
	return false
}

// 清除闯关卡使用记录
func (us *UserService) ClearSaoleiTicketRecord(userId int) {
	rcache.SetSaoleiTicketUser(userId, 0)
}

// 设置牧场新密码
func (us *UserService) SetMuchangPwd(userId int, old, new string) (bool, string) {
	user := us.GetUserById(userId)
	if user.McPwd != "" {
		if old != "" && utils.Md5(old) != user.McPwd {
			return false, "旧密码错误！"
		}
	}
	if us.CheckNeedPwd(old, user.McPwd) {
		return false, "旧密码错误！"
	}
	if new != "" {
		new = utils.Md5(new)
	}
	us.GetDb().Model(user).Update(gin.H{"fieldpwd": new})
	return true, "设置新密码成功！"
}

// 设置仓库新密码
func (us *UserService) SetCangKuPwd(userId int, old, new string) (bool, string) {
	user := us.GetUserById(userId)
	if us.CheckNeedPwd(old, user.CkPwd) {
		return false, "旧密码错误！"
	}
	if new != "" {
		new = utils.Md5(new)
	}
	us.GetDb().Model(user).Update(gin.H{"ckpwd": new})
	return true, "设置新密码成功！"
}

func (us *UserService) CheckNeedPwd(inputPwd, truePwdMd5 string) bool {
	need_pass := false
	if truePwdMd5 != "" {
		if inputPwd == "" || utils.Md5(inputPwd) != truePwdMd5 {
			need_pass = true
		}
	}
	return need_pass
}

// 取交易所金钱
func (us *UserService) GetPmsMoneys(userId int) {
	user := us.GetUserById(userId)
	userInfo := us.GetUserInfoById(userId)
	us.GetDb().Model(userInfo).Update(gin.H{"sj": userInfo.Sj + userInfo.Paisj, "paisj": 0, "paiyb": 0})
	us.GetDb().Model(user).Update(gin.H{"money": user.Money + user.PaiMoney, "yb": user.Yb + userInfo.Paiyb, "paimoney": 0})
}

// 将威望转换为贵族威望
func (us *UserService) GivePrestige(userId, num int) bool {
	return us.GetDb().Model(&models.User{ID: userId}).Where("prestige>=?", num).
		Update(gin.H{"jprestige": gorm.Expr("jprestige+?", num), "prestige": gorm.Expr("prestige-?", num)}).RowsAffected > 0
}

// 消耗金币自动战斗次数
func (us *UserService) DecreaseAutoJb(userId int) bool {
	return us.GetDb().Model(&models.User{ID: userId}).Where("sysautosum>0").Update(gin.H{"sysautosum": gorm.Expr("sysautosum-1")}).RowsAffected > 0
}

// 消耗元宝自动战斗次数
func (us *UserService) DecreaseAutoYb(userId int) bool {
	return us.GetDb().Model(&models.User{ID: userId}).Where("maxautofitsum>0").Update(gin.H{"maxautofitsum": gorm.Expr("maxautofitsum-1")}).RowsAffected > 0
}

// 新增玩家
func (us *UserService) CreateUser(account, nickname, password string, pet int, img, sex string, phone string) (bool, string) {

	us.OptSvc.Begin()
	defer us.OptSvc.Rollback()
	user := &models.User{}
	us.GetDb().Where("name=? or nickname=?", account, nickname).First(user)
	if user.ID > 0 {
		if user.Account == account {
			return false, "已存在相同的用户名！"
		} else {
			return false, "已存在相同的昵称！"
		}
	}
	phone = "13800001111"
	if phone != "13800001111" {
		userCnt := struct {
			Cnt int
		}{}
		us.GetDb().Raw("select count(1) as cnt from player where phone=?", phone).Scan(&userCnt)
		if userCnt.Cnt >= 3 {
			return false, "单个手机号最多只能注册3个账号！"
		}
	}
	now := utils.NowUnix()
	user = &models.User{
		Account:   account,
		PasswdMd5: utils.Md5(password),
		Nickname:  nickname,
		Sex:       sex,
		Regtime:   now,
		Headimg:   img,
		BagPlace:  30,
		BasePlace: 40,
		McPlace:   10,
		TgPlace:   1,
		Money:     0,
		Yb:        0,
		OpenMap:   "1",
		Phone:     phone,
	}
	if us.GetDb().Create(user).Error != nil {
		return false, "服务器错误！u"
	}
	userInfo := &models.UserInfo{
		Uid:        user.ID,
		Bbshow:     5,
		GpcGroupId: 0,
	}
	if us.GetDb().Create(userInfo).Error != nil {
		return false, "服务器错误！ui"
	}
	ok, newPet := us.OptSvc.PetSrv.CreatPetById(user, pet)
	if !ok {
		return false, "服务器错误！pe"
	}
	us.GetDb().Model(user).Update(gin.H{"mbid": newPet.ID})
	us.OptSvc.Commit()
	return true, "注册成功！恭喜您进入游戏！"
}

func (us *UserService) GetIdToken(id int) (string, error) {
	return rcache.GetIdToken(id)
}
