package rcache

import (
	"pokemon/common/persistence"
)

var RdbOperator = NewRdbHandler(persistence.GetRedisConn)

func RPush(key, value string) (int, error) {
	// 往队列右处插入值
	return RdbOperator.RPush(key, value).Int()
}

func LLen(key string) (int, error) {
	// 取队列长度
	return RdbOperator.LLen(key).Int()
}

func LRange(key string, start, stop int) ([][]byte, error) {
	// 返回列表中指定区间内的元素，取[start, stop]的区间
	return RdbOperator.LRange(key, start, stop).ByteSlices()
}

func LTrim(key string, start, stop int) (bool, error) {
	// 对一个列表进行修剪,只剩下[start, stop]的区间
	return RdbOperator.LTrim(key, start, stop).Bool()
}
