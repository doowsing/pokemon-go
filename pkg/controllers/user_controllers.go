package controllers

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/devfeel/dotweb/framework/crypto/uuid"
	"github.com/gin-gonic/gin"
	ginapp "pokemon/pkg/ginapp"
	"pokemon/pkg/rcache"
	"pokemon/pkg/services"
	"pokemon/pkg/utils"
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
	user := gapp.OptSrv.UserSrv.GetByUserNameAndPwd(account, passwd)
	if user == nil {
		gapp.JSONDATAOK("账号或密码错误！", nil)
		return
	} else if user.Password != "" && user.Password != "0" {
		gapp.JSONDATAOK("您已被禁止登陆！", nil)
		return
	}
	session := gapp.Session()
	session.Set("username", user.Account)
	session.Set("nickname", user.Nickname)
	session.Set("name", user.Account)
	session.Set("id", user.ID)
	fmt.Printf("userid:%s\n", user.ID)
	//session.Set("LoginApiState", 1) // 用不着了
	_ = session.Save()
	err := gapp.OptSrv.UserSrv.SetIdToken(user.ID, session.SessionId())
	if err != nil {
		fmt.Printf("记录错误！err:%s\n", err.Error())
	}

	gapp.JSONDATAOK("", gin.H{"ptoken": session.GetToken()})
}

func Captcha(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	// 需要返回图像的base64字符串，并将验证码存入redis
	gapp.JSONDATAOK("", gin.H{"captcha": "base64_STR", "cap-uuid": uuid.NewV4().String()})
}

func PhoneCaptcha(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	captcha := c.PostForm("captcha")
	capUUID := c.PostForm("cap-uuid")
	if capUUID == "" {
		gapp.JSONDATAOK("请重新获取验证码！", gin.H{"result": false})
		return
	}
	trueCaptcha, _ := rcache.Get("cap_uuid_" + capUUID)
	if trueCaptcha == nil {
		gapp.JSONDATAOK("验证码已过期，请重新获取验证码！", gin.H{"result": false})
		return
	}
	if string(trueCaptcha) != captcha {
		gapp.JSONDATAOK("验证码错误，请重新输入验证码！", gin.H{"result": false})
		return
	}
	phoneNumber := c.PostForm("phone")
	if phoneNumber == "" {
		gapp.JSONDATAOK("请检查输入的手机号码是否正确！", gin.H{"result": false})
		return
	}
	// 需要返回图像的base64字符串，将手机验证码存入redis
	gapp.JSONDATAOK("", gin.H{"result": true})
}

func SetMuchangPwd(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	old := c.Query("oldpasswd")
	new := c.Query("newpasswd")
	ok, msg := gapp.OptSrv.UserSrv.SetMuchangPwd(gapp.Id(), old, new)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func SetCangKuPwd(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	old := c.Query("oldpasswd")
	new := c.Query("newpasswd")
	ok, msg := gapp.OptSrv.UserSrv.SetCangKuPwd(gapp.Id(), old, new)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func GetPmsMoneys(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.OptSrv.UserSrv.GetPmsMoneys(gapp.Id())
	gapp.JSONDATAOK("取回交易所资金成功！", gin.H{"result": true})
}

func CheckLogin(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	// 需要返回图像的base64字符串，并将验证码存入redis
	gapp.JSONDATAOK("", gin.H{"result": true})
}

func (uc *UserController) ActiveLogin(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	session := gapp.Session()
	nowtime := utils.NowUnix()
	if lastLoginTime, ok := session.Get("lastlogintime"); ok {
		if nowtime-lastLoginTime.(int) <= 5 {
			gapp.String(`<script>alert('两次登录间隔不能小于5秒')</script>`)
			gapp.Redirect("/passport/login.php")
			return
		}
	}
	session.Set("lastlogintime", nowtime)
	defer session.Save()
	gapp.OptSrv.SysSrv.LoginLog(session.MustGet("username").(string), utils.GetClientIp(c), nowtime)

	gapp.Redirect("/")
}

func (uc *UserController) Logout(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	if gapp.Id() > 0 {
		session := gapp.Session()
		session.Clear()
		session.Session().Options.MaxAge = 0
		session.Save()
	}
	gapp.Redirect("/passport/login.php")
}
