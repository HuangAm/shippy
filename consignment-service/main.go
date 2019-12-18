// consignment_service/main.go 货运服务
package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	pb "shippy/consignment-service/proto/consignment" // 导入 protoc 自动生成的包
)

const (
	PORT = ":50051" // 常量一般全大写
)

// 仓库接口
type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error) // 存放新货物
	GetAll() []*pb.Consignment                                   // 获取仓库中所有的货物
}

// 我们存放了多批货物的仓库，实现了 IRepository 接口
type Repository struct {
	consignments []*pb.Consignment
}

// 创建一个新的货运，并更新 仓库 属性 consignments
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// 定义微服务
type service struct {
	repo Repository
}

// 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端
//
// 托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// 接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

// 获取目前所有托运的货物
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	allConsignments := s.repo.GetAll()
	return &pb.Response{Consignments: allConsignments}, nil
}

func main() {

	// 设置 gRPC 服务器
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listen on: %s\n", PORT)
	s := grpc.NewServer()
	repo := Repository{}
	// 在我们的 gRPC 服务器上注册微服务，这会将我们的代码和 *.pb.go 中的各种 interface 对应起来
	pb.RegisterShippingServiceServer(s, &service{repo})

	// 在 gRPC 服务器上注册 reflection
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
