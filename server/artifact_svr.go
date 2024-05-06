package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"

	"artifact_svr/db"
	pb "artifact_svr/grpc_gen/artifact_svr"
	"artifact_svr/grpc_gen/converter"
	"artifact_svr/rpc"

	"google.golang.org/grpc"
)

// server 是 ArtifactService 的服务器类型
type server struct {
	pb.UnimplementedArtifactServiceServer
}

// 耗时计算函数
// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	fmt.Printf("%s 执行耗时: %s\n", name, elapsed)
// }

// 存储文物数据的服务
func (s *server) StoreArtifact(ctx context.Context, req *pb.StoreArtifactReq) (*pb.StoreArtifactResp, error) {
	id := generateID()
	resp := &pb.StoreArtifactResp{
		ArtifactId: id,
	}

	rpc_resp, err := rpc.Client.ConvertToGltf(context.Background(), &converter.ConvertReq{
		IsBin:     true,
		NeedDraco: false,
		NoZip:     true,
		File:      req.File,
		Type:      fileType2Str(req.Type),
	})
	if err != nil {
		log.Fatalf("rpc err: %v", err)
	}
	filepath := fmt.Sprintf("save_3d_file/%s.glb", id)
	err = os.WriteFile(filepath, rpc_resp.File, 0644)
	if err != nil {
		log.Fatalf("could not write file: %v", err)
		return nil, err
	}

	c_filepath := fmt.Sprintf("save_3d_file/C_%s.glb", id)
	// 命令和参数
	cmd := exec.Command("gltfpack", "-i", filepath, "-o", c_filepath, "-tc", "-cc")

	// 执行命令
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("could not exec cmd:", err)
		return nil, err
	}

	// 信息插入数据库
	if req.Info != nil {
		err = db.InsertArtifact(id, req.Info.Name, req.Info.Description, req.Info.Location, filepath, c_filepath)
	} else {
		err = db.InsertArtifact(id, "", "", "", filepath, c_filepath)
	}
	if err != nil {
		fmt.Println("Failed to insert artifact, err:", err)
		return nil, err
	}

	return resp, nil
}

// 修改文物信息的服务
func (s *server) UpdateArtifact(ctx context.Context, req *pb.UpdateArtifactReq) (*pb.UpdateArtifactResp, error) {
	err := db.UpdateArtifact(req.ArtifactId, req.Info.Name, req.Info.Description, req.Info.Location)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateArtifactResp{}, nil
}

// 查询文物数据的服务
func (s *server) QueryArtifact(ctx context.Context, req *pb.QueryArtifactReq) (*pb.QueryArtifactResp, error) {
	modelList, err := db.QueryArtifacts(req.ArtifactId)
	if err != nil || len(modelList) == 0 {
		return nil, err
	}
	model := modelList[0]
	resp := &pb.QueryArtifactResp{}
	if req.NeedBasicInfo {
		resp.BasicInfo = &pb.ArtifactInfo{
			Name:        model.Name,
			Description: model.Description,
			Location:    model.Location,
		}
	}
	if req.NeedSource_3DFile {
		resp.Source_3DFile, err = readFileFromUri(model.Source3DFileURI)
		if err != nil {
			return nil, err
		}
	}
	if req.NeedCompressed_3DFile {
		resp.Compressed_3DFile, err = readFileFromUri(model.Compressed3DFileURI)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func fileType2Str(fileType pb.FileType) string {
	switch fileType {
	case pb.FileType_FBX:
		return "fbx"
	case pb.FileType_STL:
		return "stl"
	case pb.FileType_OBJ:
		return "obj"
	case pb.FileType_STP:
		return "stp"
	case pb.FileType_IGES:
		return "iges"
	default:
		return "err"
	}
}

func generateID() string {
	// 使用当前时间戳作为ID的一部分
	timestamp := time.Now().UnixNano()

	// 使用随机数生成剩余部分
	rand.Seed(time.Now().UnixNano())
	randomPart := fmt.Sprintf("%05d", rand.Intn(100000))

	// 组合时间戳和随机数部分生成ID
	id := fmt.Sprintf("%d%s", timestamp, randomPart)
	return id
}

func readFileFromUri(uri string) ([]byte, error) {
	file, err := os.ReadFile(uri)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Run() error {
	lis, err := net.Listen("tcp", ":8997")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	var opts = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(256 * 1024 * 1024),
		grpc.MaxSendMsgSize(256 * 1024 * 1024),
	}
	// 创建 gRPC 服务器
	s := grpc.NewServer(opts...)
	// 注册 ArtifactService 服务
	pb.RegisterArtifactServiceServer(s, &server{})
	log.Println("Server started listening on port :8997...")
	// 启动 gRPC 服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
