package controllers

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/psampaz/slice"
	"log"
	"net/http"
	"pokemon/common/rcache"
	ginapp "pokemon/game/ginapp"
	"pokemon/game/models"
	"pokemon/game/services"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"pokemon/game/utils/captcha"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var UserCtl = NewUserController()

type UserController struct {
	service *services.UserService
}

func NewUserController() *UserController {
	return &UserController{service: services.NewUserService(nil)}
}

func Login(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	account := c.PostForm("account")
	passwd := c.PostForm("password")

	valid := validation.Validation{}
	valid.Required(account, "account").Message("账号不能为空")
	valid.Required(passwd, "password").Message("密码不能为空")
	if valid.HasErrors() {
		gapp.JSONDATAOK(valid.Errors[0].Message, nil)
		return
	}
	user := gapp.OptSvc.UserSrv.GetByUserNameAndPwd(account, passwd)
	if user == nil {
		gapp.JSONDATAOK("账号或密码错误！", nil)
		return
	} else if user.Password != "" && user.Password != "0" {
		gapp.JSONDATAOK("您已被禁止登陆！", nil)
		return
	}

	ip := utils.GetClientIp(c)
	users := rcache.GetIPUsers(ip)
	if users == nil {
		users = make(map[int]int)
	}
	find := false
	for i, _ := range users {
		if i == user.ID {
			find = true
			break
		}
	}

	if !find && len(users) >= rcache.IpLimitCount {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "当前IP活跃账号过多，不可进入登录！", "data": nil})
		return
	}
	now := int(time.Now().Unix())
	users[user.ID] = now
	rcache.SetIPUsers(ip, users)

	bs, err := common.GenerateVerifyToken(user.ID, user.Account)
	if err != nil {
		log.Printf("生成PTOKEN错误！err:%s\n", err.Error())
		gapp.JSONDATAOK("登录失败！请咨询管理员！", nil)
		return
	}
	tokenStr := string(bs)
	err = rcache.SetIdToken(user.ID, tokenStr)
	if err != nil {
		log.Printf("存入PTOKEN错误！err:%s\n", err.Error())
	}

	gapp.JSONDATAOK("", gin.H{"ptoken": tokenStr})
}

func Captcha(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	// 需要返回图像的base64字符串，并将验证码存入redis
	id, answer, base64Str, err := captcha.NewCaptcha()
	if err != nil {
		gapp.JSONDATAOK("服务器生成验证码出错！", nil)
		return
	}
	rcache.SetCaptchaAnswer(id, answer)
	gapp.JSONDATAOK("", gin.H{"captcha": base64Str, "cap-uuid": id})
}

func PhoneCaptcha(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	captchaAnswer := c.Query("captcha")
	capUUID := c.Query("cap-uuid")
	if capUUID == "" {
		gapp.JSONDATAOK("请重新获取验证码！", gin.H{"result": false})
		return
	}
	trueCaptcha := rcache.GetCaptchaAnswer(capUUID)
	if trueCaptcha == "" {
		gapp.JSONDATAOK("验证码已过期，请重新获取验证码！", gin.H{"result": false})
		return
	}
	if trueCaptcha != captchaAnswer {
		gapp.JSONDATAOK("验证码错误，请重新输入验证码！", gin.H{"result": false})
		return
	}
	phoneNumber := c.Query("phone")
	if phoneNumber == "" {
		gapp.JSONDATAOK("请检查输入的手机号码是否正确！", gin.H{"result": false})
		return
	}
	// 需要返回图像的base64字符串，将手机验证码存入redis
	gapp.JSONDATAOK("", gin.H{"result": true})
}

func CheckUsername(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	account := c.Query("account")
	nickname := c.Query("nickname")

	var ok bool
	var msg string
	defer func() {
		_ = msg
		_ = ok
		gapp.JSONDATAOK(msg, gin.H{"result": ok})
	}()
	valid := validation.Validation{}
	valid.Required(account, "account").Message("账号不能为空")
	if l := utf8.RuneCountInString(account); l < 4 || l > 16 {
		msg = "账号长度需在4~16"
		return
	}
	valid.Required(nickname, "nickname").Message("昵称不能为空")
	if l := utf8.RuneCountInString(nickname); l < 4 || l > 16 {
		msg = "昵称长度需在4~16"
		return
	}
	user := &models.User{}
	gapp.OptSvc.GetDb().Where("name=? or nickname=?", account, nickname).First(user)
	if user.ID > 0 {
		if user.Account == account {
			msg = "已存在相同的用户名！"
		} else {
			msg = "已存在相同的昵称！"
		}
		return
	}
	ok = true
	return
}

func Register(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	account := c.PostForm("account")
	nickname := c.PostForm("nickname")
	passwd := c.PostForm("password")
	img := c.PostForm("img")
	pet := c.PostForm("pet")
	phone := c.PostForm("phone")
	//phoneCaptcha := c.PostForm("phonecaptcha")

	var ok bool
	var msg string
	defer func() {
		_ = msg
		_ = ok
		gapp.JSONDATAOK(msg, gin.H{"result": ok})
	}()
	valid := validation.Validation{}
	valid.Required(account, "account").Message("账号不能为空")
	if l := utf8.RuneCountInString(account); l < 4 || l > 16 {
		msg = "账号长度需在4~16"
		return
	}
	valid.Required(nickname, "nickname").Message("昵称不能为空")
	if l := utf8.RuneCountInString(nickname); l < 4 || l > 16 {
		msg = "昵称长度需在4~16"
		return
	}
	valid.Required(passwd, "password").Message("密码不能为空")
	if l := utf8.RuneCountInString(passwd); l < 8 || l > 16 {
		msg = "密码长度需在8~16"
		return
	}
	// ^[A-Za-z\d$@$!%*?&.]{6, 16}
	ok, err := regexp.MatchString(`^[A-Za-z\d$@!%*?&.]{6, 16}`, passwd)
	if err != nil {
		msg = err.Error()
		return
	}

	if phone != "13800001111" {
		valid.Phone(phone, "phone").Message("手机号码格式不对！")
	}
	if valid.HasErrors() {
		msg = valid.Errors[0].Message
		return
	}

	if img == "" || !strings.Contains("123456", img) {
		img = "1"
	}
	sex := "帅哥"
	if slice.ContainsString([]string{"2", "4", "6"}, img) {
		sex = "美女"
	}
	petId := 1
	Id2PetId := map[string]int{
		"1": 1,
		"2": 13,
		"3": 23,
		"4": 32,
		"5": 42,
	}
	if petId, ok = Id2PetId[pet]; !ok {
		petId = 1
	}
	ok, msg = gapp.OptSvc.UserSrv.CreateUser(account, nickname, passwd, petId, img, sex, phone)
}

func SetMuchangPwd(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	old := c.Query("oldpasswd")
	new := c.Query("newpasswd")
	ok, msg := gapp.OptSvc.UserSrv.SetMuchangPwd(gapp.Id(), old, new)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func SetCangKuPwd(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	old := c.Query("oldpasswd")
	new := c.Query("newpasswd")
	ok, msg := gapp.OptSvc.UserSrv.SetCangKuPwd(gapp.Id(), old, new)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func GetPmsMoneys(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.OptSvc.UserSrv.GetPmsMoneys(gapp.Id())
	gapp.JSONDATAOK("取回交易所资金成功！", gin.H{"result": true})
}

func CheckLogin(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	// 需要返回图像的base64字符串，并将验证码存入redis
	gapp.JSONDATAOK("", gin.H{"result": true})
}

func Logout(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	rcache.SetIdToken(gapp.Id(), "")
	gapp.JSONDATAOK("", nil)
}

func EmailMsg(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	msgs := []*models.EmailMsg{}
	gapp.OptSvc.GetDb().Where("uid=?", gapp.Id()).Order("id desc").Limit(10).Find(&msgs)
	msgData := []gin.H{}
	for _, m := range msgs {
		msgData = append(msgData, gin.H{"time": m.Time.Format("2006-01-02 15:04:05"), "content": m.Content})
	}
	gapp.JSONDATAOK("", gin.H{"msg": msgData})
}
