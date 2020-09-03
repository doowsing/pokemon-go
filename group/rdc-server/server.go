package rdc_server

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func InitServer() {
	arith := new(GroupHandle)
	rpc.Register(arith)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	//if e != nil {
	//	log.Fatal("listen error:", e)
	//}
	//http.Serve(l, nil)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}
