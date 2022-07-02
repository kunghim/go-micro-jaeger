package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go-micro-jaeger/proto/hello"
	"go-micro-jaeger/proto/notice"
	"log"
	"time"
)

type HelloService struct {
	NoticeServer notice.NoticeService
}

func (h HelloService) Call(ctx context.Context, request *hello.CallRequest, response *hello.CallResponse) error {
	log.Println("this is HelloService.Call")
	response.Msg = "Hello, " + request.Name
	// 获取上游传下来的链路，StartSpanFromContext 会将当前 context 注入到 xx 中，
	// 用于将当前的 span 传递到下一层服务，这样调用 notice 服务才能衔接上链路。
	span, ctx := opentracing.StartSpanFromContext(ctx, "server.Call")
	// 函数执行结束后关闭 span
	defer span.Finish()
	// 设置一个 tag
	span.SetTag("CallRequest.Name", request.Name)
	// 调用 notice 服务的 send 接口
	sendResponse, err := h.NoticeServer.Send(ctx, &notice.SendRequest{Name: request.Name})
	if err != nil {
		return err
	}
	log.Println("执行成功 -> ", sendResponse.Msg)
	// 模拟执行成功（可忽略）
	time.Sleep(1 * time.Second)
	ext.Error.Set(span, false)
	return nil
}
