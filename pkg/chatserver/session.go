package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func decodeJson(d []byte) (map[interface{}]interface{}, error) {
	m := map[interface{}]interface{}{}
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type session struct {
	value map[interface{}]interface{}
}

func (s *session) IntGet(name interface{}) int {
	return s.value[name.(string)].(int)
}

func (s *session) StrGet(name interface{}) string {
	return s.value[name.(string)].(string)
}

func GetSession(sessionId string) *session {
	conn := GetR()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("get", "session_"+sessionId))
	if err != nil {
		fmt.Printf("redis server err:%s\n", err)
		return nil
	}
	dataMap, err := decodeJson(data)
	if err != nil {
		fmt.Printf("redis data:%s\n", data)
		fmt.Printf("redis data decode err:%s\n", err)
		return nil
	}
	return &session{value: dataMap}
}
