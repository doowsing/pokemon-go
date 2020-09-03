package services

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var MsgQueue chan string
var isConn bool
var addr = flag.String("addr", "localhost:1986", "http service address")

func AddMsgQueue(msg string) {
	if isConn {
		MsgQueue <- msg
	}
}

func StartChatClient() {

	flag.Parse()
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())
	for {
		MsgQueue = make(chan string)
		connectChat(u.String())
		time.Sleep(3)
	}
}

func connectChat(wsUrl string) {
	c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Println("connect failed!", err)
	}
	defer c.Close()
	isConn = true
	for {
		select {
		case msg := <-MsgQueue:
			err := c.WriteMessage(websocket.TextMessage, []byte("announceAll"+msg))
			if err != nil {
				log.Println("connect failed!", err)
				isConn = false
				return
			}
		}
	}
}
