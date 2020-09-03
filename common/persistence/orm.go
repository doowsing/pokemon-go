package persistence

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"pokemon/common/config"
)

// 持久化 之 gorm  mysql
var orm *gorm.DB
var logOrm *gorm.DB

//根据配置初始化gorm 打开数据库连接
func InitMysql() {
	conf := config.Config().DBCfg
	var err error
	orm, err = gorm.Open(conf.Dtype,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			conf.User, conf.Password, conf.Addr, conf.Port, conf.Name, "utf8"))
	orm.DB().SetMaxIdleConns(500)
	//orm.Begin()
	if err != nil {
		panic(err)
	}
	logOrm, err = gorm.Open(conf.Dtype,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			"root", conf.Password, conf.Addr, conf.Port, "poke_log", "utf8"))

	if err != nil {
		panic(err)
	}

	orm.LogMode(conf.Debug)
	//logOrm.LogMode(conf.Debug)
}

// 获取gorm全局实例
func GetOrm() *gorm.DB {
	return orm
}

func GetLogOrm() *gorm.DB {
	return logOrm
}
