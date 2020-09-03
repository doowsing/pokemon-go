package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
	"pokemon/pkg/common"
	"pokemon/pkg/models"
	"pokemon/pkg/rcache"
	"pokemon/pkg/repositories"
	"pokemon/pkg/utils"
	"strconv"
	"time"
)

type UserService struct {
	baseService
	repo   *repositories.UserRepository
	rdrepo *rcache.UserRedisRepository
}

func NewUserService(osrc *OptService) *UserService {
	us := &UserService{repo: repositories.NewUserRepository(), rdrepo: rcache.NewUserRedisRepository()}
	us.SetOptSrc(osrc)
	return us
}

func (us *UserService) Logout(id int) bool {
	us.DelIdToken(id)
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
	return GetCardTitle(codeName)
}

func (us *UserService) GetIdToken(id int) (string, error) {
	return redis.String(rcache.Get("token_" + strconv.Itoa(id)))
}

func (us *UserService) DelIdToken(id int) (bool, error) {
	return rcache.Delete("token_" + strconv.Itoa(id))
}

func (us *UserService) SetIdToken(id int, token string) error {
	return rcache.SetEx("token_"+strconv.Itoa(id), token, common.LoginExpireTime)
}

func (us *UserService) UpdateIdToken(id int) (bool, error) {
	return rcache.Expire("token_"+strconv.Itoa(id), common.LoginExpireTime)
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
	rkey := "today_sl_user"
	str, err := rcache.Hget(rkey, strconv.Itoa(userId))
	if err == nil && len(str) > 0 {
		if string(str) == "1" {
			return true
		}
	}
	return false
}

//  设置为已扫过雷
func (us *UserService) SetSaoleiStatus(userId int) {
	rkey := "today_sl_user"
	rcache.Hset(rkey, strconv.Itoa(userId), "1")
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
	mainPet := us.OptSrc.PetSrv.GetPetById(user.Mbid)
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
	rkey := "today_is_use_ticket"
	str, err := rcache.Hget(rkey, strconv.Itoa(userId))
	if err == nil && len(str) > 0 {
		if string(str) == "1" {
			return true
		}
	}
	return false
}

// 清除闯关卡使用记录
func (us *UserService) ClearSaoleiTicketRecord(userId int) {
	rkey := "today_is_use_ticket"
	rcache.Hset(rkey, strconv.Itoa(userId), 0)
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
	if userInfo.Paisj > 0 {
		us.GetDb().Model(userInfo).Update(gin.H{"sj": userInfo.Sj + userInfo.Paisj})
	}
	if user.PaiMoney > 0 || userInfo.Paiyb > 0 {
		us.GetDb().Model(user).Update(gin.H{"money": user.Money + user.PaiMoney, "yb": user.Yb + userInfo.Paiyb})
	}
}

// 将威望转换为贵族威望
func (us *UserService) GivePrestige(userId, num int) bool {
	return us.GetDb().Model(&models.User{ID: userId}).Where("prestige>=?", num).
		Update(gin.H{"jprestige": gorm.Expr("jprestige+?", num), "prestige": gorm.Expr("prestige-?", num)}).RowsAffected > 0
}
