package chat

import (
	"encoding/json"
)

var chatUserHub = make(map[int]*ChatUser)

func GetUserFromHub(id int) (*ChatUser, bool) {
	u, ok := chatUserHub[id]
	return u, ok
}

func AddUser(u *ChatUser) {
	oldUser, ok := chatUserHub[u.info.Id]
	if ok {
		CloseUser(oldUser)
	}
	chatUserHub[u.info.Id] = u
}

func DelUser(u *ChatUser) {
	u1, ok := GetUserFromHub(u.info.Id)
	if ok && u == u1 {
		delete(chatUserHub, u.info.Id)
	}
}

func Send2ALL(v interface{}) {
	data, err := json.Marshal(v)
	if err == nil {
		for _, u := range chatUserHub {
			if u.info != nil {
				u.SendJson(data)
			}
		}
	}

}

func Send2One(t string, sender map[string]interface{}, receiver int, msg string) bool {
	for _, u := range chatUserHub {
		if u.info != nil || u.info.Id == receiver {
			u.SendStruct(t, sender, msg)
			return true
		}
	}
	return false
}

func Send2Group(t string, u *ChatUser, msg interface{}) {

}

func Send2Family(t string, u *ChatUser, msg interface{}) {

}
