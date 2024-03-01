package main

import (
	handler "adservice/handler"
	pb "adservice/proto"
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

// grpc监听端口和地址
const PORT = 50010
const ADDRESS = "127.0.0.1"

func main() {
	// ip和端口
	ipport := ADDRESS + ":" + strconv.Itoa(PORT)

	// -----------注册到consul上-------------------

	//初始化consul配置，默认端口和地址
	/*
		config := &Config{
			Address:   "127.0.0.1:8500",
			Scheme:    "http",
			Transport: transportFn(),
		}
	*/
	consulConfig := api.DefaultConfig()
	// 创建consul对象
	consulClient, err_consul := api.NewClient(consulConfig)
	if err_consul != nil {
		fmt.Println("consul创建对象报错：", err_consul)
		return
	}

	// 告诉consul即将注册到服务到信息
	reg := api.AgentServiceRegistration{
		// 标签
		Tags: []string{"adservice"},
		// 访问名称
		Name: "adservice",
		// 地址
		Address: ADDRESS,
		// 端口
		Port: PORT,
	}

	// 注册grpc服务到consul上
	err_agent := consulClient.Agent().ServiceRegister(&reg)
	if err_agent != nil {
		fmt.Println("consul注册grpc失败：", err_agent)
		return
	}

	//-----------------------grpc代码----------------------------------
	// 初始化grpc对象
	grpcServer := grpc.NewServer()

	// 注册服务
	pb.RegisterAdServiceServer(grpcServer, new(handler.AdService))

	// 设置监听
	listien, err := net.Listen("tcp", ipport)
	if err != nil {
		fmt.Println("监听报错:", err)
		return
	}
	defer listien.Close()

	// 启动服务
	fmt.Println("服务启动成功...")
	// 监听
	err_grpc := grpcServer.Serve(listien)
	if err_grpc != nil {
		fmt.Println("grpc服务启动报错:", err)
		return
	}
}
