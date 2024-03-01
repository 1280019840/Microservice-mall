package handler

import (
	"bytes"
	"context"
	"log"

	pb "emailservice/proto"
)

// 发送邮件
type DummyEmailService struct{}

// 日志
var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func (s *DummyEmailService) SendOrderConfirmation(ctx context.Context, in *pb.SendOrderConfirmationRequest) (out *pb.Empty, e error) {
	logger.Printf("邮件已经发送到： %s .", in.Email)
	out = new(pb.Empty)
	return out, nil
}
