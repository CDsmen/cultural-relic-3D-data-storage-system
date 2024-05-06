package main

import (
	"fmt"

	"artifact_svr/db"
	"artifact_svr/server"

	"artifact_svr/rpc"
)

func Init() {
	rpc.Init()

	err := db.InitMySQL()
	if err != nil {
		fmt.Println("Failed to init MySQL:", err)
		return
	}
}

func main() {
	Init()
	err := server.Run()
	if err != nil {
		fmt.Println("Failed to run server:", err)
		return
	}
}
