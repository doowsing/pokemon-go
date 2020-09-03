package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:1986", "http service address")
var blankRegexp = regexp.MustCompile(`\s+`)
var upgrader = websocket.Upgrader{} // use default options

var loginChatUser = make(map[string]*ChatUser)

type ChatUser struct {
	id       int
	nickname string
	conn     *websocket.Conn
	msgType  int
	msgChan  chan []byte
	stopChan chan int
	isVip    bool
	isMerge  bool
}

func (u *ChatUser) IsValid() bool {
	if u.id <= 0 {
		//fmt.Printf("error!no valid user!\n")
		return false
	}
	return true
}

func (u *ChatUser) Vip() string {
	vipStr := ""
	if u.isVip {
		vipStr = `<font color="#ff0000">(Vip)</font>`
	}
	return vipStr
}
func (u *ChatUser) Merge() string {
	mergeStr := ""
	if u.isVip {
		mergeStr = `<img src="../images/merge.gif" />`
	}
	return mergeStr
}

func (u *ChatUser) Close() {
	if u.IsValid() {
		u.id = 0
		u.stopChan <- 1
		u.conn.Close()
		delete(loginChatUser, u.nickname)
	}
}

func (u *ChatUser) StartSend() {
	fmt.Printf("start send...\n")
	defer close(u.msgChan)
	defer close(u.stopChan)
	for {
		var msg []byte
		select {
		case msg = <-u.msgChan:
			u.conn.WriteMessage(u.msgType, msg)
		case <-u.stopChan:
			fmt.Printf("conn is stoped...\n")
			return
		}
	}
}

func (u *ChatUser) AddMsg(msg string) {
	if u.IsValid() {
		// 加入到u用户发送写队列中
		//u.conn.WriteMessage(u.msgType, []byte(msg))
		u.msgChan <- append([]byte(msg), 0)
	}
}
func (u *ChatUser) Send2One(nickname, msg string) {
	// 私聊
	if nickname == u.nickname {
		u.AddMsg("SYSM|不能对自己说话！")
		return
	}
	if one, ok := loginChatUser[nickname]; ok {
		one.AddMsg(msg)
	} else {
		u.AddMsg(fmt.Sprintf("SYSM|%s 不在线！", nickname))
	}
}
func (u *ChatUser) Send2All(msg string) {
	// 全聊
	Send2ALL(msg)
}

func (u *ChatUser) ProcessMsg(mt int, msg []byte) {
	u.msgType = mt
	msgStr := string(msg)
	msgItems := blankRegexp.Split(msgStr, -1)
	switch msgItems[0] {
	case "announceAll":
		u.Send2All(msgItems[1])
	case "login":
		session := GetSession(msgItems[1])
		if session == nil {
			u.Close()
			return
		}
		//fmt.Printf("login:%s\n", session)
		if u.id = session.IntGet("id"); u.id != 0 {
			u.nickname = session.StrGet("nickname")
			u.stopChan = make(chan int, 1)
			u.msgChan = make(chan []byte)
			go u.StartSend()
			u.AddMsg("L|mask")
			u.AddMsg("SYSI|欢迎, " + u.nickname)
			if olduser, ok := loginChatUser[u.nickname]; ok {
				fmt.Printf("closed by myself\n")
				olduser.Close()
			}
			loginChatUser[u.nickname] = u
		} else {
			u.Close()
		}
		break
	case "CHAT":
		msgItems = blankRegexp.Split(msgStr, 2)
		u.Send2All(strings.TrimSpace(fmt.Sprintf("C|$%s`%s%s说：%s", u.nickname, u.Vip(), u.Merge(), msgItems[1])))
		break
	case "WP":
		msgItems = blankRegexp.Split(msgStr, -1)
		fmt.Printf("私聊：%v, len:%d\n", msgItems, len(msgItems))
		u.Send2One(msgItems[1], fmt.Sprintf("WP|$%s`%s对你说：%s", u.nickname, u.Vip(), msgItems[2]))
		u.AddMsg(fmt.Sprintf("WP|你对$%s`说：%s", msgItems[1], msgItems[2]))
		break
	case "W":
		u.AddMsg("W")
		break
	case "SGCHAT":
		// 组队聊天
		u.AddMsg("W")
		break
	case "GCHAT":
		// 工会聊天
		u.AddMsg("W")
		break
	default:
		u.AddMsg("SYSM|Error command")
		fmt.Printf("never prepared data:%s", msgStr)
	}
}

func Send2ALL(msg string) {
	for _, v := range loginChatUser {
		v.AddMsg(msg)
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	user := &ChatUser{conn: c}
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			user.Close()
			break
		}
		log.Printf("recv: %s", message)
		user.ProcessMsg(mt, message)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
