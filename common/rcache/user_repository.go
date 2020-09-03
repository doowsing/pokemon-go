package rcache

import (
	"pokemon/game/common"
	"strconv"
)

const (
	TOKEN   = "token_"
	CAPTCHA = "captcha_"

	IpLimitCount = 15
)

func GetIdToken(id int) (string, error) {
	data, ok := GetGCache().Get(TOKEN + strconv.Itoa(id))
	if ok {
		return data.(string), nil
	} else {
		return RdbOperator.Get(TOKEN + strconv.Itoa(id)).String()
	}
}

func DelIdToken(id int) (bool, error) {
	GetGCache().Del(TOKEN + strconv.Itoa(id))
	return RdbOperator.Delete(TOKEN + strconv.Itoa(id)).Bool()
}

func SetIdToken(id int, token string) error {
	GetGCache().SetWithTTL(TOKEN+strconv.Itoa(id), token, common.LoginExpireTime)
	return RdbOperator.SetEx(TOKEN+strconv.Itoa(id), token, common.LoginExpireTime).Error()
}

func UpdateIdToken(id int) (bool, error) {
	if token, err := GetIdToken(id); err == nil {
		GetGCache().SetWithTTL(TOKEN+strconv.Itoa(id), token, common.LoginExpireTime)
	}
	return RdbOperator.Expire(TOKEN+strconv.Itoa(id), common.LoginExpireTime).Bool()
}

func GetIPUsers(ip string) map[int]int {
	list, ok := GetGCache().Get("ip_users" + ip)
	if !ok {
		return nil
	}
	users, ok := list.(map[int]int)
	if !ok {
		return nil
	}
	return users
}

func SetIPUsers(ip string, users map[int]int) {
	GetGCache().Set("ip_users"+ip, users)
}

func SetCaptchaAnswer(id, answer string) {
	RdbOperator.SetEx(CAPTCHA+id, answer, 60*5)
}

func GetCaptchaAnswer(id string) string {
	a, _ := RdbOperator.Get(CAPTCHA + id).String()
	return a
}
