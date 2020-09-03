package rcache

import (
	"fmt"
	"pokemon/game/common"
	"time"
)

const (
	rankCache      = "rankCache"
	rankCache_time = 60
)

//获取总独立ip个数
func CountIps() (int, error) {
	return RdbOperator.SCARD(common.IPKey).Int()
}

//插入日活ip,独立ip
func InsertIp(ip string) {
	if len(ip) <= 0 {
		return
	}
	now := time.Now()
	date := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	key := fmt.Sprintf("%s::%s", common.IPKey, date)
	//独立
	RdbOperator.SADD(common.IPKey, ip)
	//日活
	RdbOperator.SADD(key, ip)
}

//获取日活ip
func CountUV() (int, error) {
	now := time.Now()
	date := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	key := fmt.Sprintf("%s::%s", common.IPKey, date)
	return RdbOperator.SCARD(key).Int()
}

func GetRankList() (map[string][]map[string]interface{}, error) {
	data := make(map[string][]map[string]interface{})
	err := RdbOperator.Get(rankCache).Struct(&data)
	return data, err
}

func SetRankList(data interface{}) {
	RdbOperator.SetEx(rankCache, data, rankCache_time)
}
