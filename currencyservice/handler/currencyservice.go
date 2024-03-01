package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "currencyservice/proto"
)

type CurrencyService struct{}

// 获得货币
func (s *CurrencyService) GetSupportedCurrencies(ctx context.Context, in *pb.Empty) (out *pb.GetSupportedCurrenciesResponse, e error) {
	data, err := ioutil.ReadFile("data/currency_conversion.json")

	if err != nil {
		return nil, status.Errorf(codes.Internal, "加载货币数据失败: %+v", err)
	}

	currencies := make(map[string]float32)
	if err := json.Unmarshal(data, &currencies); err != nil {
		return nil, status.Errorf(codes.Internal, "解析货币数据失败: %+v", err)
	}

	fmt.Printf("货币: %v\n", currencies)

	out = new(pb.GetSupportedCurrenciesResponse)

	out.CurrencyCodes = make([]string, 0, len(currencies))

	for k := range currencies {
		out.CurrencyCodes = append(out.CurrencyCodes, k)
	}
	return out, nil
}

// 转换
func (s *CurrencyService) Convert(ctx context.Context, in *pb.CurrencyConversionRequest) (out *pb.Money, e error) {
	data, err := ioutil.ReadFile("data/currency_conversion.json")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "加载货币失败: %+v", err)
	}
	currencies := make(map[string]float64)
	if err := json.Unmarshal(data, &currencies); err != nil {
		return nil, status.Errorf(codes.Internal, "解析货币失败: %+v", err)
	}
	fromCurrency, found := currencies[in.From.CurrencyCode]
	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "不支持的币种: %s", in.From.CurrencyCode)
	}
	toCurrency, found := currencies[in.ToCode]
	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "不支持的币种: %s", in.ToCode)
	}

	out = new(pb.Money)
	out.CurrencyCode = in.ToCode
	total := int64(math.Floor(float64(in.From.Units*10^9+int64(in.From.Nanos)) / fromCurrency * toCurrency))
	out.Units = total / 1e9
	out.Nanos = int32(total % 1e9)
	return out, nil
}
