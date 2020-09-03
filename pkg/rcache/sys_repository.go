package rcache

import (
	"fmt"
	"pokemon/pkg/common"
	"time"
)

//Redis缓存管理, 文章之类的 不细分
type SysRedisRepository struct {
	BaseRedisRepository
}

func NewRdRepository() *SysRedisRepository {
	return &SysRedisRepository{BaseRedisRepository: BaseRedisRepository{}}
}

//获取总独立ip个数
func (rr *SysRedisRepository) CountIps() (int, error) {
	return SCARD(common.IPKey)
}

//插入日活ip,独立ip
func (rr *SysRedisRepository) InsertIp(ip string) {
	if len(ip) <= 0 {
		return
	}
	now := time.Now()
	date := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	key := fmt.Sprintf("%s::%s", common.IPKey, date)
	//独立
	SADD(common.IPKey, ip)
	//日活
	SADD(key, ip)
}

//获取日活ip
func (rr *SysRedisRepository) CountUV() (int, error) {
	now := time.Now()
	date := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	key := fmt.Sprintf("%s::%s", common.IPKey, date)
	return SCARD(key)
}
