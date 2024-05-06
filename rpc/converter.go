package rpc

import (
	"log"

	pb "artifact_svr/grpc_gen/converter"

	"google.golang.org/grpc"
)

var Client pb.ConverterClient

func Init() {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(256 * 1024 * 1024)), // 设置接收消息的最大大小
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:8999", opts...)
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}

	Client = pb.NewConverterClient(conn)
}
