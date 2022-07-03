package main

import (
	"flag"
	"github.com/asim/go-micro/v3"
	wrapperTrace "github.com/go-micro/plugins/v3/wrapper/trace/opentracing"
	cons "go-micro-jaeger/constant"
	"go-micro-jaeger/handler"
	"go-micro-jaeger/jaeger"
	"go-micro-jaeger/proto/hello"
	"go-micro-jaeger/proto/notice"
	"log"
)

var jaegerAgentAddr string

func init() {
	// 获取启动参数中的 jaeger-agent 地址, 默认为：127.0.0.1:5775
	flag.StringVar(&jaegerAgentAddr, "a", cons.JaegerAgent, "set your jaeger-agent address")
	flag.Parse()
}

func main() {
	// 创建 tracer
	trace, err := jaeger.Tracer(cons.ServerTracer, jaegerAgentAddr).Create()
	if err != nil {
		log.Fatal("创建 server tracer 失败 -> ", err)
		return
	}
	// 服务结束时关闭 tracer
	defer trace.Closer.Close()

	// 创建 micro 服务
	service := micro.NewService(
		// 设置 micro 服务名称
		micro.Name(cons.ServerMicroServer),
		// 加入 opentracing 的中间件
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper(trace.Tracer)),
	)
	// 初始化 micro 服务
	service.Init()

	// 获取 micro-notice 服务的 noticeService，才能在 Call 中调用 notice send
	noticeService := notice.NewNoticeService(cons.NoticeMicroServer, service.Client())
	err = hello.RegisterHelloWorldHandler(service.Server(), handler.HelloService{NoticeServer: noticeService})
	if err != nil {
		log.Fatal("注册 server service 失败 -> ", err)
		return
	}

	// 启动服务
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
