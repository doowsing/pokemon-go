package chat

import (
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
)

const (
	base_url = "http://127.0.0.1:8080"
)

type HttpData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func GetUserInfo(token string) (*UserInfo, error) {
	_, body, errs := gorequest.New().Get(base_url+"/chat/user").
		Set("ptoken", token).End()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	httpData := &HttpData{Data: &UserInfo{}}
	//log.Printf("user url:%s json :%s\n", resp.Request.URL, body)
	err := json.Unmarshal([]byte(body), httpData)
	if err != nil {
		return nil, err
	}
	if httpData.Code != 200 {
		return nil, errors.New("登录失败！")
	}
	//userInfo := &UserInfo{}
	//data, ok := httpData.Data.(map[string]interface{})
	//if !ok {
	//	return nil, errors.New("login failed!")
	//}
	//userInfo.Id=int(data["id"])
	//log.Printf("user json :%s\n", httpData.Data)
	//err = json.Unmarshal(httpData.Data, userInfo)
	//if err != nil {
	//	return nil, err
	//}
	return httpData.Data.(*UserInfo), err
}

func SetUserUnReady(userId int, groupId string) {
	gorequest.New().Get(base_url + "/chat/SetGroupUnReady").Send(map[string]interface{}{"id": userId, "uuid": groupId}).End()
}
