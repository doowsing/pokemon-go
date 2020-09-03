package repositories

import (
	"github.com/jinzhu/gorm"
	"pokemon/pkg/persistence"
)

//动态连接查询
type Where func(*gorm.DB) *gorm.DB

//通用输入输出
type Out interface{}

//DAO 接口
type IRepository interface {
	Insert(out Out) *gorm.DB
	Exec(where Where, out Out, offset int, limit int) bool
	Select(where Where, out Out) bool
	SelectMany(where Where, out Out) bool
	Del(where Where, out Out) *gorm.DB
	Update(out Out, params map[string]interface{}) *gorm.DB
}

type BaseRepository struct {
	db     *gorm.DB
	TsOpen bool
}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{TsOpen: false, db: persistence.GetOrm()}
}

//接口实现自检
var _ IRepository = &BaseRepository{}

//获取数据库实例
func (br *BaseRepository) _db() *gorm.DB {
	return persistence.GetOrm()
}

func (br *BaseRepository) SetTsDb(tsDb *gorm.DB) {
	br.TsOpen = true
	br.db = tsDb
}

func (br *BaseRepository) SetNoTsDb() {
	br.TsOpen = false
	br.db = br._db()
}

func (br *BaseRepository) BeginTs() {
	br.TsOpen = true
	br.db = br._db().Begin()
}

func (br *BaseRepository) EndTs() {
	br.TsOpen = false
	br.db = br._db()
}

func (br *BaseRepository) CommitTs() {
	if br.TsOpen {
		br.db.Commit()
		br.EndTs()
	}
}

func (br *BaseRepository) CallbackTs() {
	if br.TsOpen {
		br.db.Callback()
		br.EndTs()
	}
}

func (br *BaseRepository) GetThisDb() *gorm.DB {
	if br.TsOpen {
		//fmt.Println("tsopen is opening!")
	} else {
		//fmt.Println("tsopen is closed!")
	}

	return br.db
}

//通用查询
func (br *BaseRepository) Exec(where Where, out Out, offset int, limit int) bool {
	return !br.GetThisDb().Scopes(where).Offset(offset).Limit(limit).Find(out).RecordNotFound()
}

//插入数据
func (br *BaseRepository) Insert(out Out) *gorm.DB {
	return br.GetThisDb().Create(out)
}

//查单个
func (br *BaseRepository) Select(where Where, out Out) bool {
	return br.Exec(where, out, 0, 1)
}

//查全部
func (br *BaseRepository) SelectMany(where Where, out Out) bool {
	return br.Exec(where, out, 0, -1)
}

//更新数据 ,单个 或者 多个
func (br *BaseRepository) Update(out Out, params map[string]interface{}) *gorm.DB {
	return br.GetThisDb().Model(out).Update(params)
}

//更新数据 ,单个 或者 多个
func (br *BaseRepository) UpdateWhere(where Where, out Out, params map[string]interface{}) *gorm.DB {
	return br.GetThisDb().Scopes(where).Model(out).Update(params)
}

//完全更新
func (br *BaseRepository) Save(out Out) *gorm.DB {
	return br.GetThisDb().Save(out)
}

//Count
func (br *BaseRepository) Count(out Out, where Where) (err error, count int) {
	return br.GetThisDb().Model(out).Scopes(where).Count(&count).Error, count
}

//删除数据, 单个或者多个
func (br *BaseRepository) Del(where Where, out Out) *gorm.DB {
	return br.GetThisDb().Scopes(where).Delete(out)
}

// 创建表
func (br *BaseRepository) CreateTable(out Out) error {
	return br.GetThisDb().CreateTable(out).Error
}

// 迁移表
func (br *BaseRepository) AutoMigrate(out Out) error {
	return br.GetThisDb().AutoMigrate(out).Error
}

// 查询表是否存在
func (br *BaseRepository) HasTable(tableName string) bool {
	return br.GetThisDb().HasTable(tableName)
}

// 查询表是否存在
func (br *BaseRepository) ExecSql(sql string, values ...interface{}) *gorm.DB {
	return br.GetThisDb().Exec(sql, values...)
}
