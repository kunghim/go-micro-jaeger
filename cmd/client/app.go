package main

import (
	"context"
	"github.com/asim/go-micro/v3"
	wrapperTrace "github.com/go-micro/plugins/v3/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	cons "go-micro-jaeger/constant"
	"go-micro-jaeger/jaeger"
	"go-micro-jaeger/proto/hello"
	"log"
)

func main() {
	// 创建 tracer
	trace, err := jaeger.Tracer(cons.ClientTracer, cons.TracerAgent).Create()
	if err != nil {
		log.Fatal("创建 client tracer 失败 -> ", err)
		return
	}
	// 服务结束时关闭 tracer
	defer trace.Closer.Close()

	// 创建 micro 服务
	service := micro.NewService(
		// 设置 micro 服务名称
		micro.Name(cons.ClientMicroServer),
		// 加入 opentracing 的中间件
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	// 初始化 micro 服务
	service.Init()
	ctx, span, err := wrapperTrace.StartSpanFromContext(context.Background(), opentracing.GlobalTracer(), "client")
	name := "张三"
	defer span.Finish()
	span.SetTag("name", name)
	helloWorldService := hello.NewHelloWorldService(cons.ServerMicroServer, service.Client())
	callResponse, err := helloWorldService.Call(ctx, &hello.CallRequest{Name: name})
	if err != nil {
		log.Fatal("调用 notice 服务的 send 接口失败 -> ", err)
		return
	}
	log.Println("执行成功 -> ", callResponse.Msg)
}
