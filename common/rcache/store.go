package rcache

const (
	smshopgoods      = "smshopgoods"
	smshopgoods_time = 3600
	zhekou_buyed_num = "zhekou_buyed_num"
	djshopgoods      = "djshopgoods"
	djshopgoods_time = 3600
)

func GetDjShopGood() (map[string]interface{}, error) {
	goods := make(map[string]interface{})
	err := RdbOperator.Get(djshopgoods).Struct(&goods)
	return goods, err
}

func SetDjShopGood(goods map[string]interface{}) {
	RdbOperator.SetEx(djshopgoods, goods, djshopgoods_time)
}

func GetSmShopGood() (map[string]interface{}, error) {
	goods := make(map[string]interface{})
	err := RdbOperator.Get(smshopgoods).Struct(&goods)
	return goods, err
}

func SetSmShopGood(goods map[string]interface{}) {
	RdbOperator.SetEx(smshopgoods, goods, smshopgoods_time)
}

func GetZKGoodNumList() (map[int]int, error) {
	id2num := make(map[int]int)
	err := RdbOperator.Hmget(zhekou_buyed_num).Struct(&id2num)
	return id2num, err
}

// 清理购买记录
func ClearZKGoodNumList() {
	RdbOperator.Delete(zhekou_buyed_num)
}

func GetZKGoodNum(pid int) (int, error) {
	num, err := RdbOperator.Hget(zhekou_buyed_num, pid).Int()
	return num, err
}

func SetZKGoodNum(pid, newNum int) {
	RdbOperator.Hset(zhekou_buyed_num, pid, newNum)
}
