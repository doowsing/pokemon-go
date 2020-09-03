package scheduled

import (
	"github.com/parnurzeal/gorequest"
	"log"
)

var baseServer = "http://127.0.0.1/scheduled"
var chatServer = "http://127.0.0.1:2020/check"

var name2func = map[string]func(){
	"CheckUnExpireProp": CheckUnExpireProp,
	"DelZeroProp":       DelZeroProp,
	"EndSSBattle":       EndSSBattle,
	"ClearSaoLei":       ClearSaoLei,
	"CheckChatServer":   CheckChatServer,
}

func CheckUnExpireProp() {
	_, _, errs := gorequest.New().Get(baseServer + "/check-unexpired-prop").End()
	if len(errs) > 0 {
		log.Printf("CheckUnExpireProp err:%s\n", errs[0])
		return
	}
}

func DelZeroProp() {
	_, _, errs := gorequest.New().Get(baseServer + "/del-zero-prop").End()
	if len(errs) > 0 {
		log.Printf("DelZeroProp err:%s\n", errs[0])
		return
	}
}

func EndSSBattle() {
	_, _, errs := gorequest.New().Get(baseServer + "/end-ss-battle").End()
	if len(errs) > 0 {
		log.Printf("EndSSBattle err:%s\n", errs[0])
		return
	}
}

func ClearSaoLei() {
	_, _, errs := gorequest.New().Get(baseServer + "/clear-saolei").End()
	if len(errs) > 0 {
		log.Printf("ClearSaoLei err:%s\n", errs[0])
		return
	}
}

func CheckChatServer() {
	_, _, errs := gorequest.New().Get(chatServer).End()
	if len(errs) > 0 {
		log.Printf("chat server check err:%s\n", errs[0])
		return
	}
}
