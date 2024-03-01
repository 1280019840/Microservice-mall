package handler

import (
	"bytes"
	"context"
	"log"
	"strconv"

	creditcard "github.com/durango/go-credit-card"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "paymentservice/proto"
)

// 日志
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

type PaymentService struct{}

// 结算
func (s *PaymentService) Charge(ctx context.Context, in *pb.ChargeRequest) (out *pb.ChargeResponse, e error) {
	card := creditcard.Card{
		Number: in.CreditCard.CreditCardNumber,
		Cvv:    strconv.FormatInt(int64(in.CreditCard.CreditCardCvv), 10),
		Year:   strconv.FormatInt(int64(in.CreditCard.CreditCardExpirationYear), 10),
		Month:  strconv.FormatInt(int64(in.CreditCard.CreditCardExpirationMonth), 10),
	}
	out = new(pb.ChargeResponse)
	if err := card.Validate(); err != nil {
		return out, status.Errorf(codes.InvalidArgument, err.Error())
	}

	logger.Printf(`事务处理: %s, Amount: %s%d.%d`, in.CreditCard.CreditCardNumber, in.Amount.CurrencyCode, in.Amount.Units, in.Amount.Nanos)

	out.TransactionId = uuid.NewString()
	return out, nil
}
