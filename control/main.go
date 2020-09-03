package control

import (
	rpc_group "pokemon/common/rpc-client/rpc-group"
	"sync"
	"time"
)

var smap = make(map[string]bool)
var mutex = &sync.RWMutex{}

func setFlag(serverName string, value bool) {
	mutex.Lock()
	defer mutex.Unlock()
	smap[serverName] = value
}

func SetFlag(serverName string, value bool) {
	flag := GetFlag(serverName)
	if flag == value {
		return
	}
	setFlag(serverName, value)
}

func GetFlag(serverName string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return smap[serverName]
}

func CheckGroupRpc() {
	serverName := "Group serve"
	for {
		err := rpc_group.CheckConnectError()
		if err != nil {
			SetFlag(serverName, false)
			// 重启
		} else {
			SetFlag(serverName, true)
		}
		time.Sleep(time.Second)
	}
}

func main() {
	CheckGroupRpc()
}
