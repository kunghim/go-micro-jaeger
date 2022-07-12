package main

import (
	"context"
	"errors"
	"github.com/asim/go-micro/v3"
	wrapperTrace "github.com/go-micro/plugins/v3/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	cons "go-micro-jaeger/constant"
	"go-micro-jaeger/jaeger"
	"go-micro-jaeger/proto/hello"
	"log"
	"time"
)

func main() {
	// 创建 tracer
	closer, err := jaeger.NewTracer(cons.ClientTracer)
	if err != nil {
		log.Fatal("创建 client tracer 失败 -> ", err)
		return
	}
	// 服务结束时关闭 tracer
	defer closer.Close()
	tracer := opentracing.GlobalTracer()

	// 创建 micro 服务
	service := microService()

	// client parent span 父Span
	ctx, parentSpan, err := wrapperTrace.StartSpanFromContext(context.Background(), tracer, "simple.client.parent")
	defer parentSpan.Finish()
	parentSpan.LogFields(traceLog.String("info", "this is simple client parent span"))
	parentSpan.SetTag("name", "simple client parent")

	// client child span 子Span A
	childOneSpan := tracer.StartSpan("simple.client.childOne", opentracing.ChildOf(parentSpan.Context()))
	defer childOneSpan.Finish()
	childOne(childOneSpan)

	// client child span 子Span B
	childTwoSpan := tracer.StartSpan("simple.client.childTwo", opentracing.ChildOf(parentSpan.Context()))
	defer childTwoSpan.Finish()
	childTwo(childTwoSpan)

	// 调用 hello 服务 call 接口
	helloWorldService := hello.NewHelloWorldService(cons.HelloMicroServer, service.Client())
	callResponse, err := helloWorldService.Call(ctx, &hello.CallRequest{Name: "张三"})
	if err != nil {
		log.Fatal("调用 hello 服务的 call 接口失败 -> ", err)
		return
	}
	log.Println("执行成功 -> ", callResponse.Msg)

	parentSpan.LogFields(traceLog.String("info", "finish business..."))
}

func microService() micro.Service {
	// 创建 micro 服务
	service := micro.NewService(
		// 设置 micro 服务名称
		micro.Name(cons.ClientMicroServer),
		// 加入 opentracing 的中间件
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	// 初始化 micro 服务
	service.Init()
	return service
}

// 客户端业务处理A
func childOne(span opentracing.Span) {
	time.Sleep(100 * time.Millisecond)
	span.LogFields(traceLog.String("info", "this is simple client child one span"))
	span.SetTag("name", "simple.client.child.one")
	// TODO BUSINESS A
}

// 客户端业务处理B
func childTwo(span opentracing.Span) {
	time.Sleep(150 * time.Millisecond)
	span.LogFields(traceLog.String("info", "this is simple client child two span"))
	err := errors.New("test fail addition")
	ext.LogError(span, err)
	span.LogFields(traceLog.Error(errors.New("test fail addition")))
	span.SetTag("name", "simple.client.child.two")
	// TODO BUSINESS B

	//  client child span 子Span B 的 子Span C
	time.Sleep(200 * time.Millisecond)
	childSpan := opentracing.GlobalTracer().StartSpan("simple.client.child.two.child", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	childSpan.LogFields(traceLog.String("info", "this is simple client child two child span"))
	childSpan.SetTag("name", "simple.client.child.two.child")
	// TODO BUSINESS C
}
