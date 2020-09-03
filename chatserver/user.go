package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
} // use default options

type UserInfo struct {
	Id        int    `json:"id"`
	NickName  string `json:"nickname"`
	IsVip     bool   `json:"is_vip"`
	IsMerge   bool   `json:"is_merge"`
	Img       string `json:"img"`
	GroupUUid string `json:"group_uuid"`
}

type ChatUser struct {
	info     *UserInfo
	conn     *websocket.Conn
	msgChan  chan []byte
	stopChan chan int
}

func (u *ChatUser) IsValid() bool {
	if u.info != nil {
		//fmt.Printf("error!no valid user!\n")
		return false
	}
	return true
}

func CloseUser(u *ChatUser) {
	if u != nil {
		if u.info != nil {

			DelUser(u)
			if u.info.GroupUUid != "" {
				go SetUserUnReady(u.info.Id, u.info.GroupUUid)
			}
		}
		u.conn.Close()
	}
}

func (u *ChatUser) SendStruct(t string, sr, data interface{}) {
	u.SendJson(NewSendMsg(t, sr, data))
}

func (u *ChatUser) SendJson(v interface{}) {
	var data []byte
	var err error
	switch v.(type) {
	case []byte:
		data = v.([]byte)
		break
	default:
		data, err = json.Marshal(v)
	}

	if err == nil {
		u.msgChan <- data
	}
}

func (u *ChatUser) SendString(msg string) {
	u.msgChan <- []byte(msg)
}

func (u *ChatUser) ToMap() map[string]interface{} {
	if u.info == nil {
		return nil
	}
	data := make(map[string]interface{})
	data["id"] = u.info.Id
	data["nickname"] = u.info.NickName
	if u.info.IsVip {
		data["is_vip"] = u.info.IsVip
	}
	if u.info.IsMerge {
		data["is_merge"] = u.info.IsMerge
	}
	if u.info.Img != "" {
		data["img"] = u.info.Img
	}
	return data
}

func (u *ChatUser) SecretChat(id int, msg string) {
	// 私聊
	if id == u.info.Id {
		u.SendString("SYSM|不能对自己说话！")
		return
	}
	if one, ok := chatUserHub[id]; ok {
		one.SendString(msg)
	} else {
		u.SendString(fmt.Sprintf("SYSM|%s 不在线！", one.info.NickName))
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (u *ChatUser) readPump() {
	defer func() {
		close(u.msgChan)
		CloseUser(u)
	}()
	u.conn.SetReadLimit(maxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := u.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		ProcessMsg(u, message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (u *ChatUser) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		CloseUser(u)
	}()
	for {
		select {
		case message, ok := <-u.msgChan:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := u.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(u.msgChan)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-u.msgChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &ChatUser{conn: conn, msgChan: make(chan []byte, 256)}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
