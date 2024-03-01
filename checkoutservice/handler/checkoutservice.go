package handler

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"checkoutservice/money"
	pb "checkoutservice/proto"
)

// 日志
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.LstdFlags)
)

type CheckoutService struct {
	CartService           pb.CartServiceClient
	CurrencyService       pb.CurrencyServiceClient
	EmailService          pb.EmailServiceClient
	PaymentService        pb.PaymentServiceClient
	ProductCatalogService pb.ProductCatalogServiceClient
	ShippingService       pb.ShippingServiceClient
}

// 下订单
func (s *CheckoutService) PlaceOrder(ctx context.Context, in *pb.PlaceOrderRequest) (out *pb.PlaceOrderResponse, e error) {
	logger.Printf("[PlaceOrder] user_id=%q user_currency=%q", in.UserId, in.UserCurrency)

	out = new(pb.PlaceOrderResponse)
	orderID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
		return nil, status.Errorf(codes.Internal, "生成订单id失败")
	}

	prep, err := s.prepareOrderItemsAndShippingQuoteFromCart(ctx, in.UserId, in.UserCurrency, in.Address)
	if err != nil {
		log.Fatal(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	total := &pb.Money{CurrencyCode: in.UserCurrency, Units: 0, Nanos: 0}
	total = money.Must(money.Sum(total, prep.shippingCostLocalized))
	for _, it := range prep.orderItems {
		multPrice := money.MultiplySlow(it.Cost, uint32(it.GetItem().GetQuantity()))
		total = money.Must(money.Sum(total, multPrice))
	}

	txID, err := s.chargeCard(ctx, total, in.CreditCard)
	if err != nil {
		log.Fatal(err)
		return nil, status.Errorf(codes.Internal, "更改卡失败: %+v", err)
	}
	logger.Printf("付款交易 (transaction_id: %s)", txID)

	shippingTrackingID, err := s.shipOrder(ctx, in.Address, prep.cartItems)
	if err != nil {
		log.Fatal(err)
		return nil, status.Errorf(codes.Unavailable, "配送错误: %+v", err)
	}

	if err := s.emptyUserCart(ctx, in.UserId); err != nil {
		logger.Printf("清空用户购物车失败： %s: %+v", in.UserId, err)
		log.Fatal(err)
	}

	orderResult := &pb.OrderResult{
		OrderId:            orderID.String(),
		ShippingTrackingId: shippingTrackingID,
		ShippingCost:       prep.shippingCostLocalized,
		ShippingAddress:    in.Address,
		Items:              prep.orderItems,
	}

	if err := s.sendOrderConfirmation(ctx, in.Email, orderResult); err != nil {
		logger.Printf("发送订单确认信息失败： %q: %+v", in.Email, err)
		log.Fatal(err)
	} else {
		logger.Printf("订单确认信息发送成功： %q", in.Email)
		log.Printf(in.Email)
	}
	out.Order = orderResult
	return out, nil
}

// 准备订单
type orderPrep struct {
	orderItems            []*pb.OrderItem
	cartItems             []*pb.CartItem
	shippingCostLocalized *pb.Money
}

// 准备订单和配送
func (s *CheckoutService) prepareOrderItemsAndShippingQuoteFromCart(ctx context.Context, userID, userCurrency string, address *pb.Address) (orderPrep, error) {
	var out orderPrep

	cartItems, err := s.getUserCart(ctx, userID)
	if err != nil {
		return out, fmt.Errorf("购物车错误: %+v", err)
	}
	orderItems, err := s.prepOrderItems(ctx, cartItems, userCurrency)
	if err != nil {
		return out, fmt.Errorf("准备订单失败: %+v", err)
	}
	shippingUSD, err := s.quoteShipping(ctx, address, cartItems)
	if err != nil {
		return out, fmt.Errorf("配送配额失败: %+v", err)
	}
	shippingPrice, err := s.convertCurrency(ctx, shippingUSD, userCurrency)
	if err != nil {
		return out, fmt.Errorf("货币转换失败: %+v", err)
	}

	out.shippingCostLocalized = shippingPrice
	out.cartItems = cartItems
	out.orderItems = orderItems
	return out, nil
}

// 配送配额
func (s *CheckoutService) quoteShipping(ctx context.Context, address *pb.Address, items []*pb.CartItem) (*pb.Money, error) {
	shippingQuote, err := s.ShippingService.GetQuote(ctx, &pb.GetQuoteRequest{
		Address: address,
		Items:   items,
	})
	if err != nil {
		return nil, fmt.Errorf("配送配额失败: %+v", err)
	}
	return shippingQuote.GetCostUsd(), nil
}

// 获得用户购物车
func (s *CheckoutService) getUserCart(ctx context.Context, userID string) ([]*pb.CartItem, error) {
	cart, err := s.CartService.GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("获得用户购物车失败: %+v", err)
	}
	return cart.GetItems(), nil
}

// 清空用户购物车
func (s *CheckoutService) emptyUserCart(ctx context.Context, userID string) error {
	if _, err := s.CartService.EmptyCart(ctx, &pb.EmptyCartRequest{UserId: userID}); err != nil {
		return fmt.Errorf("清空购物车失败: %+v", err)
	}
	return nil
}

// 准备订单项
func (s *CheckoutService) prepOrderItems(ctx context.Context, items []*pb.CartItem, userCurrency string) ([]*pb.OrderItem, error) {
	out := make([]*pb.OrderItem, len(items))
	for i, item := range items {
		product, err := s.ProductCatalogService.GetProduct(ctx, &pb.GetProductRequest{Id: item.GetProductId()})
		if err != nil {
			return nil, fmt.Errorf("获得商品失败 #%q", item.GetProductId())
		}
		price, err := s.convertCurrency(ctx, product.GetPriceUsd(), userCurrency)
		if err != nil {
			return nil, fmt.Errorf("价格转换失败 %q to %s", item.GetProductId(), userCurrency)
		}
		out[i] = &pb.OrderItem{Item: item, Cost: price}
	}
	return out, nil
}

// 货币转换
func (s *CheckoutService) convertCurrency(ctx context.Context, from *pb.Money, toCurrency string) (*pb.Money, error) {
	result, err := s.CurrencyService.Convert(context.TODO(), &pb.CurrencyConversionRequest{
		From:   from,
		ToCode: toCurrency,
	})
	if err != nil {
		return nil, fmt.Errorf("价格转换失败: %+v", err)
	}
	return result, err
}

// 结算卡
func (s *CheckoutService) chargeCard(ctx context.Context, amount *pb.Money, paymentInfo *pb.CreditCardInfo) (string, error) {
	paymentResp, err := s.PaymentService.Charge(ctx, &pb.ChargeRequest{
		Amount:     amount,
		CreditCard: paymentInfo,
	})
	if err != nil {
		return "", fmt.Errorf("不能更换卡: %+v", err)
	}
	return paymentResp.GetTransactionId(), nil
}

// 发送确认信息
func (s *CheckoutService) sendOrderConfirmation(ctx context.Context, email string, order *pb.OrderResult) error {
	_, err := s.EmailService.SendOrderConfirmation(ctx, &pb.SendOrderConfirmationRequest{
		Email: email,
		Order: order,
	})
	return err
}

// 配送订单
func (s *CheckoutService) shipOrder(ctx context.Context, address *pb.Address, items []*pb.CartItem) (string, error) {
	resp, err := s.ShippingService.ShipOrder(ctx, &pb.ShipOrderRequest{
		Address: address,
		Items:   items,
	})
	if err != nil {
		return "", fmt.Errorf("配送失败: %+v", err)
	}
	return resp.GetTrackingId(), nil
}
