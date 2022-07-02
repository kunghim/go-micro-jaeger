package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go-micro-jaeger/proto/notice"
	"log"
	"time"
)

type NoticeService struct {
}

func (n NoticeService) Send(ctx context.Context, request *notice.SendRequest, response *notice.SendResponse) error {
	log.Println("this is NoticeService.Send")
	response.Msg = "NoticeService 接收到请求啦"
	// 获取上游传下来的链路，'_' 为 context，StartSpanFromContext 会将当前 context 注入到 xx 中，
	// 用于将当前的 span 传递到下一层服务，因为 notice 是本示例中最底层，所以不需要用到了。
	span, _ := opentracing.StartSpanFromContext(ctx, "notice.Send")
	// 函数执行结束后关闭 span
	defer span.Finish()
	// 设置一个 tag
	span.SetTag("SendRequest.Name", request.Name)
	// 模拟执行 error
	ext.Error.Set(span, true)
	time.Sleep(3 * time.Second)
	return nil
}
