package rcache

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"pokemon/pkg/persistence"
	"reflect"
)

type BaseRedisRepository struct {
}

var GetR = persistence.GetRedisConn

type RDb struct {
	Data interface{}
	Err  error
}

func NewRDb(data interface{}, err error) *RDb {
	return &RDb{Data: data, Err: err}
}

func (this *RDb) ToInt() (int, error) {
	if this.Err != nil {
		return 0, this.Err
	}
	return redis.Int(this.Data, this.Err)
}

func (this *RDb) ToBool() (int, error) {
	if this.Err != nil || reflect.ValueOf(this.Data).IsNil() {
		return 0, this.Err
	}
	return redis.Int(this.Data, this.Err)
}

//获取Redis实例
func (rr *BaseRedisRepository) redis() redis.Conn {
	return GetR()
}

func Set(key string, data interface{}) (bool, error) {
	conn := GetR()

	defer conn.Close()

	value, err := Marshal(data)
	if err != nil {
		return false, err
	}

	reply, err := redis.Bool(conn.Do("SET", key, value))
	return reply, err
}

func SetEx(key string, data interface{}, second int) error {
	conn := GetR()
	defer conn.Close()

	value, err := Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, value, "EX", second)
	return err
}

func Expire(key string, second int) (bool, error) {
	conn := GetR()
	defer conn.Close()
	return redis.Bool(conn.Do("EXPIRE", key, second))
}

func Exists(key string) bool {
	conn := GetR()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := GetR()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := GetR()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func Hget(key, field string) ([]byte, error) {
	conn := GetR()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("HGET", key, field))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func Hset(key, field string, data interface{}) (bool, error) {
	conn := GetR()
	defer conn.Close()

	value, err := Marshal(data)
	if err != nil {
		return false, err
	}
	return redis.Bool(conn.Do("HSET", key, field, value))
}

func Hmget(key string) (map[string]interface{}, error) {
	conn := GetR()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("HMGET", key))
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(reply, &mapData)
	if err != nil {
		return nil, err
	}
	return mapData, nil
}

func Hmset(key string, filed2data *map[string]interface{}) (bool, error) {

	conn := GetR()
	defer conn.Close()
	var args = []interface{}{key}
	for i, v := range *filed2data {
		vJson, err := Marshal(v)
		if err != nil {
			fmt.Println("批量序列化失败，宠物模型ID：", i)
		} else {
			args = append(args, i, vJson)
		}
	}
	return redis.Bool(conn.Do("HMSET", args...))
}

func Hdel(key, field string) (bool, error) {
	conn := GetR()
	defer conn.Close()

	return redis.Bool(conn.Do("HDEL", key, field))
}

func HExist(key, field string) (bool, error) {
	conn := GetR()
	defer conn.Close()

	return redis.Bool(conn.Do("HEXIST", key, field))
}

func RPush(key, value string) (int, error) {
	// 往队列右处插入值
	conn := GetR()
	defer conn.Close()
	return redis.Int(conn.Do("RPUSH", key, value))
}

func LLen(key string) (int, error) {
	// 取队列长度
	conn := GetR()

	defer conn.Close()
	return redis.Int(conn.Do("LLEN", key))
}

func LRange(key string, start, stop int) ([][]byte, error) {
	// 返回列表中指定区间内的元素，取[start, stop]的区间
	conn := GetR()
	defer conn.Close()
	return redis.ByteSlices(conn.Do("LRANGE", key, start, stop))
}

func LTrim(key string, start, stop int) (bool, error) {
	// 对一个列表进行修剪,只剩下[start, stop]的区间
	conn := GetR()
	defer conn.Close()
	return redis.Bool(conn.Do("LTRIM", key, start, stop))
}

func SCARD(key string) (int, error) {
	conn := GetR()
	return redis.Int(conn.Do("SCARD", key))
}

func SADD(key, field string) (interface{}, error) {
	conn := GetR()
	return conn.Do("SADD", key, field)
}
func Marshal(v interface{}) ([]byte, error) {
	var bs []byte
	var err error
	switch t := v.(type) {
	case string:
		bs = []byte(t)
	case []byte:
		bs = t
	default:
		bs, err = json.Marshal(t)
	}
	return bs, err
}

func Unmarshal(response interface{}, err error) ([]byte, bool, error) {
	if response == nil && err == nil { // value does not exist
		return nil, false, nil
	}
	bs, err := redis.Bytes(response, err)
	if err != nil {
		return nil, false, err
	}

	return bs, true, nil
}
