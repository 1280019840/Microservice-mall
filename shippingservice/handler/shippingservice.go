package handler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	pb "shippingservice/proto"
)

type ShippingService struct{}

// 日志
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func (s *ShippingService) GetQuote(ctx context.Context, in *pb.GetQuoteRequest) (out *pb.GetQuoteResponse, e error) {
	logger.Print("[GetQuote] 收到请求")
	defer logger.Print("[GetQuote] 完成请求")

	// 1. 根据商品数量生成报价
	out = new(pb.GetQuoteResponse)
	quote := CreateQuoteFromCount(0)

	// 2. 生成响应
	out.CostUsd = &pb.Money{
		CurrencyCode: "USD",
		Units:        int64(quote.Dollars),
		Nanos:        int32(quote.Cents * 10000000),
	}
	return out, nil
}

func (s *ShippingService) ShipOrder(ctx context.Context, in *pb.ShipOrderRequest) (out *pb.ShipOrderResponse, e error) {
	logger.Print("[ShipOrder] 收到请求")
	defer logger.Print("[ShipOrder] 请求完成")
	// 1. 创建跟踪id
	out = new(pb.ShipOrderResponse)
	baseAddress := fmt.Sprintf("%s, %s, %s", in.Address.StreetAddress, in.Address.City, in.Address.State)
	id := CreateTrackingId(baseAddress)

	// 2. 生成响应
	out.TrackingId = id
	return out, nil
}
