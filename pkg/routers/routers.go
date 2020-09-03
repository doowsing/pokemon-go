package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"pokemon/pkg/config"
	"pokemon/pkg/controllers"
	loginmiddleware "pokemon/pkg/controllers/login_middleware"
)

var DBConn *gorm.DB

type Router struct {
	server *gin.Engine
}

//路由配置
func NewApiRouter(server *gin.Engine) *Router {
	router := &Router{server: server}
	if config.Config().EnvProd {
		//router.server.Use(middleware.NewCROSMiddleware())
	}
	return router
}

func (rt *Router) Init() {
	r := rt.server
	//r.Use(middleware.CharSet("gbk"))
	serverInit(r)
	//pokeServerInit(r)

}

func serverInit(rt *gin.Engine) {
	rt.GET("/getlast", controllers.FileServerCtl.GetLastAccess)
	rt.GET("/getnow", controllers.FileServerCtl.GetNowAccess)
	rt.GET("/api/listlen/gamelog", controllers.GetGameLogListLen)
	rt.GET("/api/listlen/loginlog", controllers.GetLoginLogListLen)
	rt.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "你好啊菜鸡")
	})
	NoLoginUrlInit(rt)
	LoginUrlInit(rt)
	AdminShowServerInit(rt)
}

func NoLoginUrlInit(r *gin.Engine) {
	r.POST("/user/login", controllers.Login)
	r.GET("/user/captcha", controllers.Captcha)
	r.GET("/user/phonecaptcha", controllers.PhoneCaptcha)
}

func LoginUrlInit(r *gin.Engine) {
	rg := r.Group("/", loginmiddleware.LoginPage())

	userGroup := rg.Group("/user")
	userGroup.GET("/checklogin", controllers.CheckLogin)
	userGroup.GET("/mc/setpwd", controllers.SetMuchangPwd)
	userGroup.GET("/ck/setpwd", controllers.SetCangKuPwd)
	userGroup.GET("/getpmsmoney", controllers.GetPmsMoneys)

	petGroup := rg.Group("/pet")
	petGroup.GET("/info/:id", controllers.PetInfo)
	petGroup.GET("/putin", controllers.PutInPet)
	petGroup.GET("/putout", controllers.PutOutPet)
	petGroup.GET("/throw", controllers.ThrowPet)
	petGroup.GET("/setmbid", controllers.SetMbid)
	petGroup.GET("/carrypets", controllers.CarryPetList)
	petGroup.GET("/allpets", controllers.AllPetList)
	petGroup.GET("/learnskill", controllers.LearnSkill)
	petGroup.GET("/updateskill", controllers.UpdateSkill)
	petGroup.GET("/show/:id", controllers.ShowPet2All)

	propGroup := rg.Group("/prop")
	propGroup.GET("/info/:id", controllers.PropInfo)
	propGroup.GET("/show-info/:pid", controllers.MPropInfo)
	propGroup.GET("/carryprops", controllers.CarryPropList)
	propGroup.GET("/use/:id", controllers.UseProp)
	propGroup.GET("/putin", controllers.PutInProp)
	propGroup.GET("/putout", controllers.PutOutProp)
	propGroup.GET("/throw", controllers.ThrowProp)
	propGroup.POST("/auction", controllers.Auction)
	propGroup.GET("/reauction", controllers.ReAuction)
	propGroup.GET("/rollauction", controllers.RollAuction)
	propGroup.GET("/purchase", controllers.Purchase)
	propGroup.GET("/shop-sell", controllers.ShopSell)
	propGroup.GET("/shop-purchase", controllers.ShopPurchase)

	npcGroup := rg.Group("/npc")
	npcGroup.GET("/petinfo", controllers.PetsPage)
	npcGroup.GET("/petinfo/offzb", controllers.PetOffZb)
	npcGroup.GET("/user", controllers.UserPage)
	npcGroup.GET("/public", controllers.PublicPage)
	npcGroup.GET("/smshop", controllers.SmShopPage)
	npcGroup.GET("/smshop/qg", controllers.SmShopQgList)
	npcGroup.GET("/djshop", controllers.DjShopPage)
	npcGroup.GET("/mc", controllers.MuchangPage)
	npcGroup.GET("/ck", controllers.CangkuPage)
	npcGroup.GET("/pms", controllers.PaiMaiPage)
	npcGroup.GET("/tjp", controllers.TieJiangPuPage)
	npcGroup.GET("/tjp/fjzb", controllers.FenJieEquip)
	npcGroup.GET("/tjp/qhzb", controllers.QiangHuaEquip)
	npcGroup.GET("/tjp/qhzb-info", controllers.QiangHuaEquipInfo)
	npcGroup.GET("/tjp/merge", controllers.MergeProps)
	npcGroup.GET("/king", controllers.KingPage)
	npcGroup.GET("/king/getprize", controllers.GetDayPrize)
	npcGroup.GET("/king/giveprestige", controllers.GivePrestige)
	npcGroup.GET("/king/openegg", controllers.Zadan)
	npcGroup.GET("/sl", controllers.SaoLeiInfo)
	npcGroup.GET("/sl/update", controllers.UpdateSaoLeiAwards)
	npcGroup.GET("/sl/start", controllers.StartSaoLei)
	npcGroup.GET("/sl/into", controllers.IntoSaoLei)
	npcGroup.GET("/sl/easter", controllers.EasterSaoLei)
	npcGroup.GET("/cwsd", controllers.PetSdPage)
	npcGroup.GET("/cwsd/evolution", controllers.PetEvolution)
	npcGroup.GET("/cwsd/merge", controllers.PetMerge)
	npcGroup.GET("/cwsd/zs", controllers.PetZhuansheng)
	npcGroup.GET("/cwsd/cqcz", controllers.PetCqCzl)
	npcGroup.GET("/cwsd/zhcz", controllers.PetZhCzl)
	npcGroup.GET("/cwsd/ss-evolution", controllers.PetSSEvolution)
	npcGroup.GET("/cwsd/sszs-info", controllers.PetSSZhuanShengInfo)
	npcGroup.GET("/cwsd/sszs", controllers.PetSSZhuanSheng)
}

func AdminShowServerInit(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	adminGroup.GET("/showegginfo", controllers.ShowEggSetting)
}

func pokeServerInit(r *gin.Engine) {
	r.GET("/passport/login.php", controllers.LoginPage)
	r.POST("/passport/dealPc.php", controllers.Login)
	r.GET("/login/login.php", loginmiddleware.LoginPage(), controllers.UserCtl.ActiveLogin)
	r.GET("/login/logout.php", loginmiddleware.LoginPage(), controllers.UserCtl.Logout)
	r.GET("/", loginmiddleware.LoginPage(), controllers.IndexPage)
	r.GET("/game.php", loginmiddleware.LoginPage(), controllers.GamePage)
	r.POST("/api/time.php", controllers.SysCtl.GetTime)
	r.GET("/iframe.php", controllers.IframePage)

	fGroup := r.Group("function", loginmiddleware.LoginRequest())
	fGroup.GET("/Welcome_Mod.php", controllers.WelcomePage)
	fGroup.GET("/toHome.php", controllers.HomePage)
	fGroup.GET("/User_Mod.php", controllers.UserPage)
	fGroup.GET("/Pets_Mod.php", controllers.PetsPage)
	fGroup.GET("/Expore_Mod.php", controllers.ExporePage)
	fGroup.GET("/Exporenew_Mod.php", controllers.ExporeNewPage)
	fGroup.GET("/Muchang_Mod.php", controllers.MuchangPage)
	fGroup.GET("/Team_Mod.php", controllers.TeamModPage)
	fGroup.GET("/fb_Mod.php", controllers.FbModPage)
	fGroup.GET("/Pai_Mod.php", controllers.PaiMaiPage)

	fGroup.GET("/ext_Online.php", controllers.ExtOnline)
	fGroup.GET("/mcGate.php", controllers.McGate)
	fGroup.GET("/getBag.php", controllers.GetBag)
	fGroup.GET("/mcbbshow.php", controllers.ShowPet)
	fGroup.GET("/usedProps.php", controllers.UsePropString)
	fGroup.GET("/getPropsInfo.php", controllers.GetPropInfo)
	fGroup.GET("/offprops.php", controllers.OffZb)
	fGroup.GET("/mapGate.php", controllers.CheckOpenMap)
	fGroup.GET("/openMap.php", controllers.OpenMap)
	fGroup.GET("/team.php", controllers.TeamInfo)
}

//func (rt *Router) InitPetRouter(pr *gin.RouterGroup) {
//	petController := controllers.NewPetController()
//	pr.GET("/getmpet", petController.ShowMPet)
//	pr.GET("/getmpets", petController.ShowALLMPet)
//	pr.GET("/getmskill", petController.ShowMSkill)
//	pr.GET("/initpettable", petController.InitTable)
//	pr.GET("/", middleware.JWT(), petController.ShowPets)
//	pr.GET("/createpet", middleware.JWT(), petController.CreatePet)
//	pr.GET("/getexp", middleware.JWT(), petController.ShowExpByLv)
//	pr.GET("/getwx", middleware.JWT(), petController.ShowWx)
//	pr.GET("/addExp", middleware.JWT(), petController.IncreatExp)
//}
//func (rt *Router) InitPropRouter(pr *gin.RouterGroup) {
//	pc := controllers.NewPropController()
//	pr.GET("/", middleware.JWT(), pc.GetBpProp)
//	pr.GET("/add", middleware.JWT(), pc.AddProp)
//	pr.GET("/use/:propid", middleware.JWT(), pc.UsePropString)
//	pr.GET("/showm/:propid", middleware.JWT(), pc.ShowMProp)
//}
//func (rt *Router) InitUserRouter(pr *gin.RouterGroup) {
//
//	userController := controllers.NewUserController()
//	pr.GET("/", middleware.JWT(), userController.GetUser)
//}
//
//func (rt *Router) InitServerRouter(pr *gin.RouterGroup) {
//	sysController := controllers.NewSysController()
//	pr.GET("/init/redis", sysController.InitRdModels)
//	pr.GET("/init/table/", sysController.InitAllTables)
//	pr.GET("/init/table/prop", sysController.InitMPropTable)
//	pr.GET("/init/table/user", sysController.InitUserTable)
//	pr.GET("/init/data/user", sysController.CreateUser)
//}
//
//func (rt *Router) InitFightRouter(fr *gin.RouterGroup) {
//	fightController := controllers.NewFightController()
//	fr.GET("/test", fightController.TestMFightInfo)
//	fr.GET("/enter/:mapid", middleware.JWT(), fightController.EnterMap)
//	fr.GET("/start/:mapid", middleware.JWT(), fightController.StartFight)
//}
