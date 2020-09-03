package app

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"pokemon/pkg/config"
	"pokemon/pkg/controllers"
	"pokemon/pkg/controllers/login_middleware"
	"pokemon/pkg/middleware"
	"pokemon/pkg/models"
	"pokemon/pkg/persistence"
	"pokemon/pkg/routers"
	"pokemon/pkg/services"
	"time"
)

//应用服务器
type App struct {
	DB     *gorm.DB
	Conf   *config.BlogConfig
	Server *gin.Engine
}

func NewApp() *App {
	return &App{}
}

//启动服务器
func (app *App) Launch() error {
	app.Conf = config.Config()
	app.initDB()
	app.initRedis()
	app.initServer()
	app.initUtils()
	app.initRouter()
	app.initTask()
	return app.Server.Run(app.Conf.ServerPort)
}

//关闭操作
func (app *App) Destory() {
	if app.DB != nil {
		app.DB.Close()
	}
	if app.Server != nil {
		app.Server = nil
	}
}

//根据配置文件初始化数据库
func (app *App) initDB() {
	app.DB = persistence.GetOrm()
}

//根据配置文件初始化Redis
func (app *App) initRedis() {
	//app.RPool = persistence.GetRedisPool()
	persistence.GetRedisConn()
}

//根据配置初始化服务器
func (app *App) initServer() {
	app.Server = gin.New()

	//配置自定义error
	app.initError()

	//配置环境模式
	if app.Conf.EnvProd {
		//app.Server.SetProductionMode()
	} else {
		//app.Server.SetDevelopmentMode()
	}
	//开启日志
	app.Server.Use(gin.Logger())
	//异常返回500
	//app.Server.Use(gin.Recovery())
	//开启Gzip压缩
	app.Server.Use(gzip.Gzip(gzip.DefaultCompression))

	//开启跨域
	app.Server.Use(middleware.Cross())
	//开启ip计数
	app.Server.Use(middleware.NewCROSMiddleware().Handle())
	//开启静态文件查找
	app.Server.Use(middleware.StaticMdl("static"))
	//开启session
	app.Server.Use(loginmiddleware.SessionMdl())
}

//初始化路由配置
func (app *App) initRouter() {
	r := routers.NewApiRouter(app.Server)
	// api/v1/xx api
	r.Init()
}

//配置自定义error
func (app *App) initError() {
	ec := controllers.NewErrorController()
	app.Server.NoRoute(ec.PageNotFound)
	app.Server.NoMethod(ec.MethodNotFound)
}

// 初始化常用内存数据
func (app *App) initUtils() {
	// 初始化数据缓存
	services.InitDataStore()

	// 设置自动访问函数
	models.SetMPetFunc(services.GetMpet, services.GetMpetByName)
	models.SetMPropFunc(services.GetMProp)
	models.SetMSkillFunc(services.GetMskill)
	models.SetTaskFunc(services.GetTask)

	return
	// 定时任务部分
	optSrv := services.NewOptService()
	go func() {
		go func() {
			time.Sleep(3)
			services.StartChatClient() // 开启聊天服务器
		}()

		time.Sleep(time.Duration(5) * time.Second)
		go optSrv.PropSrv.CheckPropExpire(60) // 清理过期物品

		time.Sleep(time.Duration(5) * time.Second)
		go optSrv.PropSrv.CheckPropValid(60) // 清理无效物品
	}()
	go optSrv.SysSrv.AutoInsertGameLog()

}

// 初始化定时任务

func (app *App) initTask() {
	sys := &services.SysService{}
	go func() {
		//sys.PrepareTestDate()
		go sys.AutoInsertGameLog()
		go sys.SaveGameLog()
		go sys.SaveLoginLog()
	}()
	optSrv := services.NewOptService()
	go func() {
		//go func() {
		//	time.Sleep(3)
		//	services.StartChatClient() // 开启聊天服务器
		//}()

		time.Sleep(time.Duration(5) * time.Second)
		go optSrv.PropSrv.CheckPropExpire(60) // 清理过期物品

		time.Sleep(time.Duration(5) * time.Second)
		go optSrv.PropSrv.CheckPropValid(60) // 清理无效物品
	}()
}
