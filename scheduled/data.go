package scheduled

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"pokemon/common/utils/path"
)

var Tasks = &TaskList{}

func init() {

	setBytes, err := ioutil.ReadFile(path.GetRootDir() + "/config/task_set.json")
	if err != nil {
		log.Printf("load tasks err:%s\n", err)
		return
	}
	data := &struct {
		GameServer string  `json:"game_server"`
		ChatServer string  `json:"chat_server"`
		TaskLists  []*Task `json:"task_list"`
	}{}
	err = json.Unmarshal(setBytes, data)
	if err != nil {
		log.Printf("load tasks err:%s\n", err)
		return
	}
	Tasks.list = data.TaskLists
	baseServer = data.GameServer
	chatServer = data.ChatServer

}
