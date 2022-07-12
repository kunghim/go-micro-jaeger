package handler

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
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

	parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "simple.hello.parent")
	defer parentSpan.Finish()
	parentSpan.LogFields(traceLog.String("info", "this is simple hello parent span"))
	parentSpan.SetTag("name", "simple hello parent")

	childOneSpan := opentracing.StartSpan("simple.hello.child.one", opentracing.ChildOf(parentSpan.Context()))
	defer childOneSpan.Finish()
	h.childOne(childOneSpan)

	childTwoSpan := opentracing.StartSpan("simple.hello.child.two", opentracing.ChildOf(parentSpan.Context()))
	defer childTwoSpan.Finish()
	return h.childTwo(ctx, request, childTwoSpan)
}

func (h *HelloService) childOne(span opentracing.Span) {
	span.LogFields(traceLog.String("info", "this is simple hello child one span"))
	span.SetTag("name", "simple.hello.child.one")
	time.Sleep(100 * time.Millisecond)

	childSpan := opentracing.GlobalTracer().StartSpan("simple.hello.child.one.next", opentracing.FollowsFrom(span.Context()))
	defer childSpan.Finish()
	childSpan.LogFields(traceLog.String("info", "this is simple hello child one next span"))
	childSpan.SetTag("name", "simple.hello.child.one.next")
	time.Sleep(160 * time.Millisecond)
}

func (h *HelloService) childTwo(ctx context.Context, request *hello.CallRequest, span opentracing.Span) error {
	span.LogFields(traceLog.String("info", "this is simple hello child two span"))
	err := errors.New("test fail addition")
	ext.LogError(span, err)
	span.LogFields(traceLog.Error(errors.New("test fail addition")))
	span.SetTag("name", "simple.hello.child.two")
	time.Sleep(120 * time.Millisecond)

	childSpan, ctx := opentracing.StartSpanFromContext(ctx, "simple.hello.child.tow.next", opentracing.FollowsFrom(span.Context()))
	defer childSpan.Finish()
	childSpan.LogFields(traceLog.String("info", "request to notice server"))
	childSpan.SetTag("name", "simple.hello.child.two.next")
	// 调用 notice 服务的 send 接口
	sendResponse, err := h.NoticeServer.Send(ctx, &notice.SendRequest{Name: request.Name})
	if err != nil {
		return err
	}
	log.Println("执行成功 -> ", sendResponse.Msg)
	return nil
}
