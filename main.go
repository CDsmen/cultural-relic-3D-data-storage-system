package main

import (
	"fmt"

	"artifact_svr/db"
	"artifact_svr/server"

	"artifact_svr/rpc"
)

func Init() error {
	rpc.Init()

	err := db.InitMySQL()
	if err != nil {
		fmt.Println("Failed to init MySQL:", err)
		return err
	}
	return nil
}

func main() {
	err := Init()
	if err != nil {
		fmt.Println("Failed Init server:", err)
		return
	}
	err = server.Run()
	if err != nil {
		fmt.Println("Failed to run server:", err)
		return
	}
}
