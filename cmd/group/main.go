package main

import (
	"pokemon/common/persistence"
	"pokemon/group/rdc-server"
)

func main() {
	persistence.InitRedisCluster()
	rdc_server.InitServer()
}
