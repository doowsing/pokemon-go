package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	"log"
)

const (
	chat_server_host_url = "http://127.0.0.1:2020"
)

type ChatMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type GroupCommandMsg struct {
	Command string      `json:"command"`
	UUid    string      `json:"uuid"`
	Receive int         `json:"receive"` //0则为全部
	Data    interface{} `json:"data"`
}

func send2ChatServe(dt interface{}) {
	//data, err := json.Marshal(&map[string]interface{}{"data": dt})
	//if err != nil {
	//	log.Printf("send chat err:%s\n", err)
	//}
	//log.Printf("send data:%s\n", data)
	data, err := json.Marshal(dt)
	if err != nil {
		log.Printf("send chat err:%s\n", err)
	}
	log.Printf("send data:%s\n", data)
	_, _, errs := gorequest.New().Post(chat_server_host_url + "/sysmsg").Type("multipart").Send(map[string]interface{}{"data": string(data)}).End()
	if len(errs) > 0 {
		log.Printf("send chat err:%s\n", errs[0])
	}
}

func AnnounceChat(msg string) {
	v := &ChatMsg{
		Type: "announce",
		Data: msg,
	}
	send2ChatServe(v)
}

func ShowPropChat(userId, propId int, propName string, propColor int) {
	v := &ChatMsg{
		Type: "prop-show",
		Data: map[string]interface{}{
			"userId":    userId,
			"propId":    propId,
			"propName":  propName,
			"propColor": propColor,
		},
	}
	send2ChatServe(v)
}

func ShowPetChat(userId, petId int, petName string) {
	v := &ChatMsg{
		Type: "pet-show",
		Data: map[string]interface{}{
			"userId":  userId,
			"petId":   petId,
			"petName": petName,
		},
	}
	send2ChatServe(v)
}

func UpdateUserInfo2Chat(userInfo gin.H) {
	v := &ChatMsg{
		Type: "userinfo",
		Data: userInfo,
	}
	send2ChatServe(v)
}

func NoticeTips(userId int, tips string) {
	v := &ChatMsg{
		Type: "tips",
		Data: gin.H{"id": userId, "tips": tips},
	}
	send2ChatServe(v)
}

func GroupNotice(data interface{}) {
	v := &ChatMsg{
		Type: "group",
		Data: data,
	}
	send2ChatServe(v)
}

// 通知聊天服务器更新用户的队伍Id
func GroupUpdateId(uuid string, receiveId int) {
	GroupNotice(&GroupCommandMsg{
		Command: "update-uuid",
		UUid:    uuid,
		Receive: receiveId,
		Data:    nil,
	})
}

// 通知用户刷新队伍
func GroupUpdate(uuid string, receiveId int) {
	GroupNotice(&GroupCommandMsg{
		Command: "update",
		UUid:    uuid,
		Receive: receiveId,
		Data:    nil,
	})
}

// 通知聊天服务器队伍解散
func GroupDissolve(uuid string) {
	GroupNotice(&GroupCommandMsg{
		Command: "dissolve",
		UUid:    uuid,
		Receive: 0,
		Data:    nil,
	})
}

// 通知聊天服务器发送聊天信息，邀请加入队伍
func GroupInvite(uuid string, receiveId int, leaderId int) {
	GroupNotice(&GroupCommandMsg{
		Command: "invite",
		UUid:    uuid,
		Receive: receiveId,
		Data:    leaderId,
	})
}

// 通知玩家进入战斗
func GroupStartFight(uuid string, data interface{}) {
	GroupNotice(&GroupCommandMsg{
		Command: "fight",
		UUid:    uuid,
		Receive: 0,
		Data:    data,
	})
}

// 通知玩家响应攻击结果
func GroupAttack(uuid string, data interface{}) {
	GroupNotice(&GroupCommandMsg{
		Command: "attack",
		UUid:    uuid,
		Receive: 0,
		Data:    data,
	})
}

// 通知玩家进入关卡翻牌
func GroupEnterCard(uuid string) {
	GroupNotice(&GroupCommandMsg{
		Command: "enter-card",
		UUid:    uuid,
		Receive: 0,
		Data:    nil,
	})
}

// 通知玩家更新翻牌结果
func GroupUpdateCard(uuid string, data interface{}) {
	GroupNotice(&GroupCommandMsg{
		Command: "do-card",
		UUid:    uuid,
		Receive: 0,
		Data:    data,
	})
}

// 通知玩家进入Boss翻牌
func GroupEnterBossCard(uuid string) {
	GroupNotice(&GroupCommandMsg{
		Command: "enter-boss-card",
		UUid:    uuid,
		Receive: 0,
		Data:    nil,
	})
}

// 通知玩家更新翻牌结果
func GroupUpdateBossCard(uuid string, data interface{}) {
	GroupNotice(&GroupCommandMsg{
		Command: "do-boss-card",
		UUid:    uuid,
		Receive: 0,
		Data:    data,
	})
}
