package main

import (
	"errors"
	"log"
	addHttp "mygo/rpc/http"
	"net"
	"net/http"
	"net/rpc"
)

// 定义一个服务对象
type Arith int

// 实现服务类型的方法
func (t *Arith) Multiply(args *addHttp.Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *addHttp.Args, quo *addHttp.Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

// 实现RPC服务端
func main() {
	arith := new(Arith)
	// 注册服务
	rpc.Register(arith)
	// 通过HTTP暴露出来
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	// 启动HTTP服务
	http.Serve(l, nil)
}
