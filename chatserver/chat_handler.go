package chat

import (
	"encoding/json"
	"log"
)

const (
	TIP      = "tips"
	LOGIN    = "login"
	ANNOUNCE = "announce"
)

type ReceiveMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type SendMsg struct {
	Type   string      `json:"type"`
	Sender interface{} `json:"sender"`
	Data   interface{} `json:"data"`
}

type GroupCommandMsg struct {
	Command string      `json:"command"`
	UUid    string      `json:"uuid"`
	Receive int         `json:"receive"` //0则为全部
	Data    interface{} `json:"data"`
}

func NewSendMsg(t string, sr interface{}, data interface{}) *SendMsg {
	return &SendMsg{
		Type:   t,
		Sender: sr,
		Data:   data,
	}
}

func ProcessMsg(u *ChatUser, msg []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("process client user:%s, msg:%s, err:%s\n", u.info.NickName, msg, err)
		}
	}()
	rMsg := &ReceiveMsg{}
	err := json.Unmarshal(msg, rMsg)
	if err != nil {
		log.Printf("receive msg err:%s\n", err)
		return
	}
	log.Printf("receive msg :%s\n", msg)
	switch rMsg.Type {
	case "login":
		token, ok := rMsg.Data.(string)
		//log.Printf("receive token:%s\n", token)
		if ok {
			user, err := GetUserInfo(token)
			if err == nil {
				u.info = user
				u.SendStruct(LOGIN, nil, user)
				AddUser(u)
			} else {
				//log.Printf("receive msg err:%s\n", err)
				u.SendStruct(TIP, nil, err.Error())
				CloseUser(u)
			}
		} else {
			u.SendStruct(TIP, nil, "登录凭证出错！")
			CloseUser(u)
		}

		break
	case "public-chat":
		msg, ok := rMsg.Data.(string)
		if ok {
			Send2ALL(NewSendMsg("public-chat", u.ToMap(), msg))
		}
		break
	case "secret-chat":

		log.Printf("receive data:%s\n", msg)
		data, ok := rMsg.Data.(map[string]interface{})
		if !ok {
			break
		}
		msg, ok := data["msg"].(string)
		if !ok {
			break
		}
		id, ok := data["receive"].(float64)
		if !ok {
			break
		}
		if u.info.Id == int(id) {
			u.SendStruct(TIP, nil, "不可以给自己发私聊")
			break
		}
		receiver, ok := GetUserFromHub(int(id))
		if ok {
			receiver.SendStruct("secret-chat", u.ToMap(), msg)
			u.SendStruct("secret-chat-result", nil, map[string]interface{}{
				"receiver": receiver.ToMap(),
				"msg":      msg,
			})
		} else {
			u.SendStruct(TIP, nil, "用户未上线！")
		}
		break
	case "group-chat":
		msg, ok := rMsg.Data.(string)
		if ok {
			Send2Group(rMsg.Type, u, msg)
		}
		break
	case "family-chat":
		msg, ok := rMsg.Data.(string)
		if ok {
			Send2Family(rMsg.Type, u, msg)
		}
		break
	case "check":
		u.SendStruct("check", nil, "ok")
		break
	default:
		log.Printf("never prepared data:%s\n", msg)
	}
}

func ProcessSysMsg(msg string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("process sys:, msg:%s, err:%s\n", msg, err)
		}
	}()
	rMsg := &ReceiveMsg{}
	err := json.Unmarshal([]byte(msg), rMsg)
	log.Printf("receive msg:%s\n", msg)
	if err != nil {
		log.Printf("receive msg err:%s\n", err)
		return
	}
	switch rMsg.Type {
	case "announce":
		Send2ALL(NewSendMsg(rMsg.Type, nil, rMsg.Data))
		break
	case "prop-show":
		data := &struct {
			UserId    int    `json:"userId"`
			PropId    int    `json:"propId"`
			PropName  string `json:"propName"`
			PropColor int    `json:"propColor"`
		}{}
		rMsg.Data = data

		err := json.Unmarshal([]byte(msg), rMsg)
		log.Printf("receive msg:%s\n", msg)
		if err != nil {
			log.Printf("receive msg err:%s\n", err)
			return
		}
		user, ok := GetUserFromHub(data.UserId)
		if ok {
			Send2ALL(NewSendMsg(rMsg.Type, user.ToMap(), rMsg.Data))
		}
		break
	case "pet-show":
		data := &struct {
			UserId  int    `json:"userId"`
			PetId   int    `json:"petId"`
			PetName string `json:"petName"`
		}{}
		rMsg.Data = data

		err := json.Unmarshal([]byte(msg), rMsg)
		log.Printf("receive msg:%s\n", msg)
		if err != nil {
			log.Printf("receive msg err:%s\n", err)
			return
		}
		user, ok := GetUserFromHub(data.UserId)
		if ok {
			Send2ALL(NewSendMsg(rMsg.Type, user.ToMap(), rMsg.Data))
		}
		break

	case "group-new", "group-change", "group-dissolve":
		break
	case "logout":
		idF, ok := rMsg.Data.(float64)
		if ok {
			id := int(idF)
			u, ok := GetUserFromHub(id)
			if ok {
				u.SendStruct(TIP, nil, "您已被系统下号！")
				CloseUser(u)
				delete(chatUserHub, id)
			}
		}
	case "userinfo":
		data := &UserInfo{}
		rMsg.Data = data

		err := json.Unmarshal([]byte(msg), rMsg)
		log.Printf("receive msg:%s\n", msg)
		if err != nil {
			log.Printf("receive msg err:%s\n", err)
			return
		}
		u, ok := GetUserFromHub(data.Id)
		if ok {
			u.info = data
		}
	case TIP:
		data := rMsg.Data.(map[string]interface{})
		id := int(data["id"].(float64))
		tips := data["tips"].(string)
		user, ok := GetUserFromHub(id)
		if ok {
			user.SendStruct("tips", nil, tips)
		}

	case "group":
		groupMsg := &GroupCommandMsg{}
		rMsg.Data = groupMsg
		err := json.Unmarshal([]byte(msg), rMsg)
		log.Printf("receive msg:%s\n", msg)
		if err != nil {
			log.Printf("receive msg err:%s\n", err)
			return
		}
		HandleGroupCommand(groupMsg)
	}
}

func HandleGroupCommand(groupMsg *GroupCommandMsg) {
	switch groupMsg.Command {
	case "update-uuid":
		user, ok := GetUserFromHub(groupMsg.Receive)
		if ok {
			user.info.GroupUUid = groupMsg.UUid
			// 通知客户端刷新界面
			user.SendStruct("group-update", nil, nil)
		}
		break
	case "clear-group-uuid":
		for _, u := range chatUserHub {
			if u.info != nil {
				u.info.GroupUUid = ""
				u.SendStruct("group-update", nil, nil)
			}
		}
		break
	case "update":
		if groupMsg.Receive > 0 {
			user, ok := GetUserFromHub(groupMsg.Receive)
			if ok && user.info.GroupUUid == groupMsg.UUid {
				user.SendStruct("group-update", nil, nil)
			}
		} else {
			for _, u := range chatUserHub {
				if u.info.GroupUUid == groupMsg.UUid {
					u.SendStruct("group-update", nil, nil)
				}
			}
		}

		break
	case "dissolve":
		for _, u := range chatUserHub {
			if u.info.GroupUUid == groupMsg.UUid {
				u.SendStruct("tips", nil, "队伍已解散！")
				u.info.GroupUUid = ""
			}
		}
		break
	case "invite":
		user, ok := GetUserFromHub(groupMsg.Receive)
		if !ok {
			break
		}
		senderId, ok := groupMsg.Data.(int)
		if !ok {
			break
		}
		sender, ok := GetUserFromHub(senderId)
		if !ok {
			break
		}
		user.SendStruct("group-invite", sender.ToMap(), groupMsg.UUid)
		break
	case "fight":
		for _, u := range chatUserHub {
			if u.info.GroupUUid == groupMsg.UUid {
				u.SendStruct("group-fight", nil, groupMsg.Data)
			}
		}
		break
	case "attack":
		for _, u := range chatUserHub {
			if u.info.GroupUUid == groupMsg.UUid {
				u.SendStruct("group-attack", nil, groupMsg.Data)
			}
		}
		break
	case "enter-card":
		for _, u := range chatUserHub {
			if u.info.GroupUUid == groupMsg.UUid {
				u.SendStruct("enter-card", nil, nil)
			}
		}
		break
	case "do-card":
		for _, u := range chatUserHub {
			if u.info.GroupUUid == groupMsg.UUid {
				u.SendStruct("do-card", nil, groupMsg.Data)
			}
		}
		break
	case "enter-boss-card":
		for _, u := range chatUserHub {
			if u.info.GroupUUid == groupMsg.UUid {
				u.SendStruct("enter-boss-card", nil, nil)
			}
		}
		break
	case "do-boss-card":
		for _, u := range chatUserHub {
			if u.info.GroupUUid == groupMsg.UUid {
				u.SendStruct("do-boss-card", nil, groupMsg.Data)
			}
		}
		break
	}
}
