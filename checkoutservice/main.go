package main

import (
	handler "checkoutservice/handler"
	pb "checkoutservice/proto"
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

func GetGrpcConn(consulClient *api.Client, serviceName string, serviceTag string) *grpc.ClientConn {
	service, _, err_service := consulClient.Health().Service(serviceName, serviceTag, true, nil)
	if err_service != nil {
		fmt.Println("获取健康服务报错：", err_service)
		return nil
	}
	s := service[0].Service
	address := s.Address + ":" + strconv.Itoa(s.Port)
	fmt.Printf("serviceName: %v\n", serviceName)

	fmt.Printf("address:%s\n", address)

	//链接grpc服务
	grpcConn, _ := grpc.Dial(address, grpc.WithInsecure())

	return grpcConn
}

const PORT = 50020
const ADDRESS = "127.0.0.1"

func main() {
	ipport := ADDRESS + ":" + strconv.Itoa(PORT)
	// ----------------注册到consul上---------------------
	// 初始化consul配置
	consulConfig := api.DefaultConfig()

	// 创建consul对象
	consulClient, err_consul := api.NewClient(consulConfig)
	if err_consul != nil {
		fmt.Println("consul创建对象报错：", err_consul)
		return
	}

	// 告诉consul即将注册到服务到信息
	reg := api.AgentServiceRegistration{
		Tags:    []string{"checkoutservice"},
		Name:    "checkoutservice",
		Address: ADDRESS,
		Port:    PORT,
	}

	// 注册grpc服务到consul上
	err_agent := consulClient.Agent().ServiceRegister(&reg)
	if err_agent != nil {
		fmt.Println("consul注册grpc失败：", err_agent)
		return
	}

	// -----------------------grpc代码-----------------------
	//  初始化grpc对象
	grpcServer := grpc.NewServer()

	// 调用其他服务
	checkoutService := &handler.CheckoutService{
		CartService:           pb.NewCartServiceClient(GetGrpcConn(consulClient, "cartservice", "cartservice")),
		CurrencyService:       pb.NewCurrencyServiceClient(GetGrpcConn(consulClient, "currencyservice", "currencyservice")),
		EmailService:          pb.NewEmailServiceClient(GetGrpcConn(consulClient, "emailservice", "emailservice")),
		ProductCatalogService: pb.NewProductCatalogServiceClient(GetGrpcConn(consulClient, "productcatalogservice", "productcatalogservice")),
		PaymentService:        pb.NewPaymentServiceClient(GetGrpcConn(consulClient, "paymentservice", "paymentservice")),
		ShippingService:       pb.NewShippingServiceClient(GetGrpcConn(consulClient, "shippingservice", "shippingservice")),
	}

	// 注册服务
	pb.RegisterCheckoutServiceServer(grpcServer, checkoutService)

	// 设置监听
	listien, err := net.Listen("tcp", ipport)
	if err != nil {
		fmt.Println("监听报错:", err)
		return
	}
	defer listien.Close()

	// 启动服务
	fmt.Println("服务启动成功。。。")

	err_grpc := grpcServer.Serve(listien)
	if err_grpc != nil {
		fmt.Println("grpc服务启动报错:", err)
		return
	}
}
