package repositories

import (
	"github.com/jinzhu/gorm"
	"pokemon/pkg/models"
)

//对应的数据库中表: users user_infos
type UserRepository struct {
	*BaseRepository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{BaseRepository: NewBaseRepository()}
}

//获取全部用户数据
func (ur *UserRepository) GetAll() ([]models.User, bool) {
	users := make([]models.User, 0)
	return users, ur.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Preload("UserInfo")

	}, &users)
}

//根据用户Id获取用户
func (ur *UserRepository) GetByUserId(id int) (*models.User, bool) {
	user := &models.User{}
	found := ur.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}, user)
	return user, found
}

//根据用户名(唯一)获取用户
func (ur *UserRepository) GetByUserName(username string) (*models.User, bool) {
	user := &models.User{}
	found := ur.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_name = ?", username)
	}, user)
	return user, found
}

//根据用户名和密码获取用户
func (ur *UserRepository) GetByUserNameAndPwd(username, password string) (*models.User, bool) {
	user := &models.User{}
	found := ur.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("account = ? and passwd_wd5 = ?", username, password)
	}, user)
	return user, found
}

//删除表中全部用户
func (ur *UserRepository) DelAll() *gorm.DB {
	err := ur.Del(func(db *gorm.DB) *gorm.DB {
		return db
	}, &models.User{})
	return err
}

//根据用户id从表中删除某个用户,与其关联的用户信息
func (ur *UserRepository) DelByUid(uid int) *gorm.DB {
	return ur.Del(func(db *gorm.DB) *gorm.DB {
		return db
	}, &models.User{})
}

func (ur *UserRepository) InitTable() error {
	if !ur.db.HasTable("user") {
		return ur.db.CreateTable(&models.User{}).Error
	} else {
		return ur.db.AutoMigrate(&models.User{}).Error
	}
}

//根据用户ID，更新用户托管空间大小
func (ur *UserRepository) SetTgPlace(UserId, Place int) *gorm.DB {
	return ur.Update(&models.User{ID: UserId},
		map[string]interface{}{"tg_place": Place})
}

//根据用户ID，更新用户剩余托管时间，单位小时
func (ur *UserRepository) SetTgTime(UserId, time int) *gorm.DB {
	return ur.Update(&models.User{ID: UserId},
		map[string]interface{}{"tg_time": time})
}

// 更新用户开启地图
func (ur *UserRepository) SetOpenMaps(UserId int, OpenMaps string) *gorm.DB {
	return ur.Update(&models.User{ID: UserId},
		map[string]interface{}{"open_map": OpenMaps})
}

// 更新用户自动战斗次数
func (ur *UserRepository) SetAutoTimes(UserId, Times int, Type string) *gorm.DB {
	var TypeName string
	if Type == "money" {
		TypeName = "auto_fight_time_m"
	} else if Type == "yb" {
		TypeName = "auto_fight_time_yb"
	} else if Type == "team" {
		TypeName = "team_auto_times"
	} else {
		return nil
	}
	return ur.Update(&models.User{ID: UserId},
		map[string]interface{}{TypeName: Times})
}

func (ur *UserRepository) SetShowTimes(UserId, ShowTimes int) *gorm.DB {
	return ur.Update(&models.User{ID: UserId},
		map[string]interface{}{"show_times": ShowTimes})
}
func (ur *UserRepository) SetPrestige(User *models.User) *gorm.DB {
	return ur.Update(User,
		map[string]interface{}{"prestige": User.Prestige})
}

func (ur *UserRepository) SetSj(UserId, AddNum int) *gorm.DB {
	return ur.Update(&models.User{ID: UserId},
		map[string]interface{}{"sj": AddNum})
}

func (ur *UserRepository) SetYb(UserId, AddNum int) *gorm.DB {
	return ur.Update(&models.User{ID: UserId},
		map[string]interface{}{"yb": AddNum})
}

func (ur *UserRepository) SetBagPlace(UserId, NewPlace, MaxPlace int) *gorm.DB {
	return ur.UpdateWhere(func(db *gorm.DB) *gorm.DB {
		return db.Where("bag_place < ?", MaxPlace)
	}, &models.User{ID: UserId},
		map[string]interface{}{"bag_place": NewPlace})
}

func (ur *UserRepository) SetCkPlace(UserId, NewPlace, MaxPlace int) *gorm.DB {
	return ur.UpdateWhere(func(db *gorm.DB) *gorm.DB {
		return db.Where("base_place < ?", MaxPlace)
	}, &models.User{ID: UserId},
		map[string]interface{}{"base_place": NewPlace})
}
