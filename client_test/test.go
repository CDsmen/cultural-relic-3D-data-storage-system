package main

import (
	pb "artifact_svr/grpc_gen/artifact_svr"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {
	id := testCreate()
	if id != "" {
		testQuery(id)
		testUpdate(id, "11", "22", "33")
		testQuery(id)
	}
}

func testCreate() string {
	f := "./assets/stl.zip"

	// 读取目录下的压缩文件
	file, _ := os.ReadFile(f)
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(256 * 1024 * 1024)), // 设置接收消息的最大大小
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
		return ""
	}
	defer conn.Close()

	client := pb.NewArtifactServiceClient(conn)
	resp, err := client.StoreArtifact(context.Background(), &pb.StoreArtifactReq{
		File: file,
		Type: pb.FileType_STL,
		Info: &pb.ArtifactInfo{
			Name:        "testName",
			Description: "testDesc",
			Location:    "testLoca",
		},
	})
	if err != nil {
		log.Fatalf("client.StoreArtifact err: %v", err)
		return ""
	}

	fmt.Printf("Create resp: %v\n", resp)
	return resp.ArtifactId
}

func testQuery(id string) {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(256 * 1024 * 1024)), // 设置接收消息的最大大小
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewArtifactServiceClient(conn)
	resp, err := client.QueryArtifact(context.Background(), &pb.QueryArtifactReq{
		ArtifactId:    id,
		NeedBasicInfo: true,
		// NeedSource_3DFile:     true,
		// NeedCompressed_3DFile: true,
	})
	if err != nil {
		log.Fatalf("client.QueryArtifact err: %v", err)
		return
	}

	fmt.Printf("Query resp: %v\n", resp)
}

func testUpdate(id, name, desc, loca string) {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(256 * 1024 * 1024)), // 设置接收消息的最大大小
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewArtifactServiceClient(conn)
	resp, err := client.UpdateArtifact(context.Background(), &pb.UpdateArtifactReq{
		ArtifactId: id,
		Info: &pb.ArtifactInfo{
			Name:        name,
			Description: desc,
			Location:    loca,
		},
	})
	if err != nil {
		log.Fatalf("client.UpdateArtifact err: %v", err)
		return
	}

	fmt.Printf("Update resp: %v\n", resp)
}
