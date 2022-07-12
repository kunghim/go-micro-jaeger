package main

import (
	"github.com/asim/go-micro/v3"
	wrapperTrace "github.com/go-micro/plugins/v3/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	cons "go-micro-jaeger/constant"
	"go-micro-jaeger/handler"
	"go-micro-jaeger/jaeger"
	"go-micro-jaeger/proto/hello"
	"go-micro-jaeger/proto/notice"
	"log"
)

func main() {
	// 创建 tracer
	closer, err := jaeger.NewTracer(cons.HelloTracer)
	if err != nil {
		log.Fatal("创建 hello tracer 失败 -> ", err)
		return
	}
	// 服务结束时关闭 tracer
	defer closer.Close()

	// 创建 micro 服务
	service := microService()
	// 获取 micro-notice 服务的 noticeService，才能在 Call 中调用 notice send
	noticeService := notice.NewNoticeService(cons.NoticeMicroServer, service.Client())
	err = hello.RegisterHelloWorldHandler(service.Server(), handler.HelloService{NoticeServer: noticeService})
	if err != nil {
		log.Fatal("注册 notice service 失败 -> ", err)
		return
	}
	// 启动服务
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}

func microService() micro.Service {
	// 创建 micro 服务
	service := micro.NewService(
		// 设置 micro 服务名称
		micro.Name(cons.HelloMicroServer),
		// 加入 opentracing 的中间件
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	// 初始化 micro 服务
	service.Init()
	return service
}
