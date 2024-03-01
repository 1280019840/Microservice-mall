package handler

import (
	"fmt"
	"math"
)

// 货币描述
type Quote struct {
	Dollars uint32
	Cents   uint32
}

// 报价描述
func (q Quote) String() string {
	return fmt.Sprintf("$%d.%d", q.Dollars, q.Cents)
}

// 根据商品数量创建报价
func CreateQuoteFromCount(count int) Quote {
	return CreateQuoteFromFloat(8.99)
}

// 创建报价
func CreateQuoteFromFloat(value float64) Quote {
	units, fraction := math.Modf(value)
	return Quote{
		uint32(units),
		uint32(math.Trunc(fraction * 100)),
	}
}
