package rcache

const (
	tjpshopgoods      = "tjpshopgoods"
	tjpshopgoods_time = 3600
	zbfj_info         = "zbfj_info"
)

func GetTJPGoods() (map[string]interface{}, error) {
	goods := make(map[string]interface{})
	err := RdbOperator.Get(tjpshopgoods).Struct(&goods)
	return goods, err
}

func SetTJPGoods(goods map[string]interface{}) {
	RdbOperator.SetEx(tjpshopgoods, goods, tjpshopgoods_time)
}

func GetZbfjTimes(userId int) (int, error) {
	times, err := RdbOperator.Hget(zbfj_info, userId).Int()
	return times, err
}

// 设置装备分解次数
func SetZbfjTimes(userId, times int) {
	RdbOperator.Hset(zbfj_info, userId, times)
}

// 清空所有的装备分解次数
func ClearZbfjTimes() {
	RdbOperator.Delete(zbfj_info)
}
