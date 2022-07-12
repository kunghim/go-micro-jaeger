package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go-micro-jaeger/proto/notice"
	"time"
)

type NoticeService struct {
}

func (n NoticeService) Send(ctx context.Context, request *notice.SendRequest, response *notice.SendResponse) error {
	response.Msg = "NoticeService 接收到请求啦"
	span, _ := opentracing.StartSpanFromContext(ctx, "notice.parent")
	defer span.Finish()
	span.SetTag("name", "simple.notice.parent")
	time.Sleep(100 * time.Millisecond)
	return nil
}
