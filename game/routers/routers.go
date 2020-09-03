package routers

import (
	"github.com/gin-gonic/gin"
	"pokemon/common/config"
	"pokemon/game/controllers"
	"pokemon/game/controllers/npc"
	"pokemon/game/ginapp"
	"pokemon/game/middleware"
)

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

	NoLoginUrlInit(rt)
	LoginUrlInit(rt)
	AdminServerInit(rt)
	ScheduledServerInit(rt)
}

func NoLoginUrlInit(r *gin.Engine) {
	r.POST("/user/login", controllers.Login)
	r.GET("/user/check", controllers.CheckUsername)
	r.POST("/user/register", controllers.Register)
	r.GET("/user/captcha", controllers.Captcha)
	r.GET("/user/phonecaptcha", controllers.PhoneCaptcha)
	r.GET("/group/SetGroupUnReady", controllers.SetGroupUnReady)
}

func LoginUrlInit(r *gin.Engine) {
	rg := r.Group("/", middleware.TokenMiddleWare())

	userGroup := rg.Group("/user")
	userGroup.GET("/checklogin", controllers.CheckLogin)
	userGroup.GET("/mc/setpwd", controllers.SetMuchangPwd)
	userGroup.GET("/ck/setpwd", controllers.SetCangKuPwd)
	userGroup.GET("/getpmsmoney", controllers.GetPmsMoneys)
	userGroup.GET("/email-sys", controllers.EmailMsg)
	userGroup.GET("/logout", controllers.Logout)

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
	propGroup.GET("/show-all", controllers.ShowProp2All)
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
	npcGroup.GET("/game-info", npc.GameInfo)
	npcGroup.GET("/city", npc.CityPage)
	npcGroup.GET("/petinfo", npc.PetsPage)
	npcGroup.GET("/petinfo/offzb", npc.PetOffZb)
	npcGroup.GET("/user", npc.UserPage)
	npcGroup.GET("/public", npc.PublicPage)
	npcGroup.GET("/smshop", npc.SmShopPage)
	npcGroup.GET("/smshop/qg", npc.SmShopQgList)
	npcGroup.GET("/djshop", npc.DjShopPage)
	npcGroup.GET("/mc", npc.MuchangPage)
	npcGroup.GET("/ck", npc.CangkuPage)
	npcGroup.GET("/pms", npc.PaiMaiPage)

	npcGroup.GET("/tjp", npc.TieJiangPuPage)
	npcGroup.GET("/tjp/fjzb", npc.FenJieEquip)
	npcGroup.GET("/tjp/qhzb", npc.QiangHuaEquip)
	npcGroup.GET("/tjp/qhzb-info", npc.QiangHuaEquipInfo)
	npcGroup.GET("/tjp/merge", npc.MergeProps)

	npcGroup.GET("/king", npc.KingPage)
	npcGroup.GET("/king/getprize", npc.GetDayPrize)
	npcGroup.GET("/king/giveprestige", npc.GivePrestige)
	npcGroup.GET("/king/openegg", npc.Zadan)

	npcGroup.GET("/sl", npc.SaoLeiInfo)
	npcGroup.GET("/sl/update", npc.UpdateSaoLeiAwards)
	npcGroup.GET("/sl/start", npc.StartSaoLei)
	npcGroup.GET("/sl/into", npc.IntoSaoLei)
	npcGroup.GET("/sl/easter", npc.EasterSaoLei)

	npcGroup.GET("/cwsd", npc.PetSdPage)
	npcGroup.GET("/cwsd/evolution", npc.PetEvolution)
	npcGroup.GET("/cwsd/merge", npc.PetMerge)
	npcGroup.GET("/cwsd/zs", npc.PetZhuansheng)
	npcGroup.GET("/cwsd/cqcz", npc.PetCqCzl)
	npcGroup.GET("/cwsd/zhcz", npc.PetZhCzl)
	npcGroup.GET("/cwsd/ss-evolution", npc.PetSSEvolution)
	npcGroup.GET("/cwsd/sszs-info", npc.PetSSZhuanShengInfo)
	npcGroup.GET("/cwsd/sszs", npc.PetSSZhuanSheng)
	npcGroup.GET("/welcome", npc.WelcomeContent)

	npcGroup.GET("/card/series", npc.CardSeries)
	npcGroup.GET("/card/user", npc.UserCardSeries)
	npcGroup.GET("/card/prizes", npc.CardPrizes)
	npcGroup.GET("/card/getprize", npc.GetCardPrize)
	npcGroup.GET("/card/title", npc.CardTitles)
	npcGroup.GET("/card/title/cancel", npc.CancelCardTitle)
	npcGroup.GET("/card/title/use", npc.UseCardTitle)

	npcGroup.GET("/family", npc.FamilyPageInfo)
	npcGroup.GET("/family/info", npc.FamilyShowInfo)
	npcGroup.POST("/family/create", npc.CreateFamily)
	npcGroup.GET("/family/apply", npc.ApplyFamily)
	npcGroup.GET("/family/disband", npc.DisbandFamily)
	npcGroup.GET("/family/upgrade-info", npc.FamilyUpgradeInfo)
	npcGroup.GET("/family/upgrade", npc.UpgradeFamily)
	npcGroup.GET("/family/donate", npc.DonateFamily)
	npcGroup.GET("/family/receive-welfare", npc.ReceiveFamilyWelFare)
	npcGroup.GET("/family/handle-apply", npc.ReplyApplyFamily)
	npcGroup.GET("/family/fire", npc.FireFamilyMember)
	npcGroup.GET("/family/change", npc.ManageFamilyAuthority)
	npcGroup.GET("/family/exit", npc.ExitFamily)
	npcGroup.GET("/family/store-upgrade", npc.UpgradeFamilyStore)
	npcGroup.GET("/family/store-purchase", npc.PurchaseFamilyStore)

	npcGroup.GET("/zhanbuwu", npc.ZhanBuWuList)
	npcGroup.GET("/zhanbuwu/use", npc.ZhanBuWuUse)

	taskGroup := rg.Group("/task")
	taskGroup.GET("/accept-list", controllers.AcceptTaskList)
	taskGroup.GET("/enable-accept-list", controllers.EnableAcceptTaskList)
	taskGroup.GET("/show/:id", controllers.TaskInfo)
	taskGroup.GET("/accept-show/:id", controllers.AcceptTaskInfo)
	taskGroup.GET("/accept", controllers.AcceptTask)
	taskGroup.GET("/finish", controllers.FinishTask)
	taskGroup.GET("/throw", controllers.ThrowTask)

	fightGroup := rg.Group("/fight")
	fightGroup.GET("/openmaps", controllers.OpenMaps)
	fightGroup.GET("/openmap", controllers.OpenMap)
	fightGroup.GET("/into-map", controllers.GoInToMap)
	fightGroup.GET("/into-fbmap", controllers.GoInToFbMap)
	fightGroup.GET("/start", middleware.LimitIp(), controllers.StartFight)
	fightGroup.GET("/attack", controllers.Attack)
	fightGroup.GET("/auto-start", controllers.AutoStartFight)
	fightGroup.GET("/auto-start/skill", controllers.AutoFightSkill)
	fightGroup.GET("/auto-start/cancel", controllers.CancelAutoStartFight)
	fightGroup.GET("/tt-usesj", controllers.TTUseSj)
	fightGroup.GET("/tt-rank", controllers.TTUserRank)
	fightGroup.GET("/users", controllers.MapUsers)
	fightGroup.GET("/catch", controllers.CatchPet)
	fightGroup.GET("/start-fb", controllers.StartFb)

	ssFightGroup := rg.Group("/ss-battle")
	ssFightGroup.GET("/user-list", controllers.SSBattleUserList)
	ssFightGroup.GET("/enter", controllers.SSBattleEnter)
	ssFightGroup.GET("/start-fight", controllers.SSBattleStartFight)
	ssFightGroup.GET("/attack", controllers.SSBattleAttack)
	ssFightGroup.GET("/use", controllers.SSBattleUseProp)
	ssFightGroup.GET("/user-info", controllers.SSBattleUserInfo)
	ssFightGroup.GET("/store", controllers.SSBattleStore)
	ssFightGroup.GET("/get-award", controllers.SSBattleGetAward)
	ssFightGroup.GET("/convert-exp", controllers.SSBattleConvertExp)
	ssFightGroup.GET("/convert-prop", controllers.SSBattleConvertProp)

	familyFightGroup := rg.Group("/family-battle")
	familyFightGroup.GET("/info", controllers.FamilyBattleInfo)
	familyFightGroup.GET("/invite-battle", controllers.FamilyBattleInvite)
	familyFightGroup.GET("/accept-battle", controllers.FamilyBattleAccept)
	familyFightGroup.GET("/start-fight", controllers.FamilyBattleStartFight)
	familyFightGroup.GET("/attack", controllers.FamilyBattleAttack)

	groupGroup := rg.Group("/group")
	groupGroup.GET("/info", controllers.GroupInfo)
	groupGroup.GET("/create", controllers.CreateGroup)
	groupGroup.GET("/dissolve", controllers.DissolveGroup)
	groupGroup.GET("/request", controllers.RequestGroup)
	groupGroup.GET("/receive", controllers.ReceiveGroup)
	groupGroup.GET("/refuse", controllers.RefuseGroup)
	groupGroup.GET("/kick-out", controllers.KickOutGroup)
	groupGroup.GET("/exit", controllers.ExitGroup)
	groupGroup.GET("/invite", controllers.InviteGroup)
	groupGroup.GET("/start-fight", controllers.GroupStartFight)
	groupGroup.GET("/attack", controllers.GroupAttack)
	groupGroup.GET("/card", controllers.GroupEnterCard)
	groupGroup.GET("/card/do", controllers.GroupDoCard)
	groupGroup.GET("/card/show", controllers.GroupAllCard)
	groupGroup.GET("/boss-card", controllers.GroupEnterBossCard)
	groupGroup.GET("/boss-card/do", controllers.GroupDoBossCard)
	groupGroup.GET("/set-status", controllers.GroupSetStatus)
	groupGroup.GET("/in-map", controllers.GroupInmap)

	chatGroup := rg.Group("/chat")
	chatGroup.GET("/user", controllers.ChatLogin)
}

func ScheduledServerInit(r *gin.Engine) {
	group := r.Group("/scheduled")
	group.GET("/check-unexpired-prop", controllers.CheckUnExpireProp)
	group.GET("/del-zero-prop", controllers.DelZeroProp)
	group.GET("/end-ss-battle", controllers.EndSSBattle)
	group.GET("/clear-saolei", controllers.ClearSaoLei)
}

func AdminServerInit(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	adminGroup.GET("/showegginfo", controllers.ShowEggSetting)
	adminGroup.GET("/redisclusterstatus", controllers.ShowRedisClusterStatus)
}

func FileServerInit(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		gapp := ginapp.NewGapp(c)
		defer ginapp.DropGapp(gapp)
		gapp.String("你好啊菜鸡")
		return
		gapp.HTML("index.html", nil)
	})
}
