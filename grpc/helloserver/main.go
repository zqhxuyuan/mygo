package main

import (
	"net"

	pb "mygo/proto" // 引入编译生成的包

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

// 定义helloService并实现约定的接口
type helloService struct{}

// HelloService ...
var HelloService = helloService{}

// 实现helloService，对应了proto的Hello service接口
// Hello service只有一个参数HelloRequest，这里还注入了Context
func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// 创建一个响应对象，这个响应对象会返回给客户端
	resp := new(pb.HelloReply)
	// HelloRequest定义了变量Name，所以这里可以获取到Name变量的值
	// HelloReply定义了变量Message，这里设置Message的值
	resp.Message = "Hello " + in.Name + "."

	return resp, nil
}

func main() {
	// 监听TCP端口
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册HelloService。注意：这里并不是小写的helloService（类型），而是注入一个变量（对象）
	pb.RegisterHelloServer(s, HelloService)

	grpclog.Println("Listen on " + Address)

	s.Serve(listen)
}
