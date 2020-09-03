package relations

import (
	"github.com/jinzhu/gorm"
	"pokemon/common/config"
	"pokemon/common/persistence"
	"pokemon/game/models"
	"pokemon/game/utils"
	"time"
)

/*
	generator.go 根据结构体(Models) 统一创建数据库关系

        1. 初始化分类
		2. 初始化默认用户

*/
func InitRelations() {
	db := persistence.GetOrm()
	// 判断存不存在表, 不存在就新建, 否则就是自动迁移(其他修改)
	if !db.HasTable("player") {
		db.CreateTable(&models.User{})
	} else {
		db.AutoMigrate(&models.User{})
	}
	initDefaultUser(db)
}

// 插入默认用户 数据库记录
func initDefaultUser(db *gorm.DB) {
	defaultUser := config.Config().DefaultClientUser
	if found := db.Find(&models.User{}, "account = ?", defaultUser).RecordNotFound(); found {
		// 使用UserService插入
		// 使用UserRepo插入
		// 直接db插入
		user := &models.User{
			Account:   defaultUser,
			Nickname:  defaultUser,
			PasswdMd5: utils.Md5(defaultUser + "ccc"),
			Regtime:   time.Now(),
			Sex:       1,
		}
		db.Create(&user)
	}
}
