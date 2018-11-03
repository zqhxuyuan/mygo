package main

import (
	"fmt"
	pb "mygo/proto" // 引入proto包

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		grpclog.Fatalln(err)
	}

	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	fmt.Println("connect to rpc server.")

	// 调用方法，客户端RPC调用服务端，需要传递请求，即HelloRequest
	reqBody := new(pb.HelloRequest)
	reqBody.Name = "gRPC"

	// 执行RPC调用，返回响应结果
	r, err := c.SayHello(context.Background(), reqBody)
	if err != nil {
		fmt.Println("error")
		grpclog.Fatalln(err)
	}
	fmt.Println(r.Message)
	grpclog.Println(r.Message)
}
