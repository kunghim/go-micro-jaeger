package main

import (
	"github.com/opentracing/opentracing-go"
	traceLog "github.com/opentracing/opentracing-go/log"
	"go-micro-jaeger/jaeger"
	"log"
	"time"
)

func main() {
	// 创建 tracer
	closer, err := jaeger.NewTracer("simple")
	if err != nil {
		log.Fatal("创建 server tracer 失败 -> ", err)
		return
	}
	// 服务结束时关闭 tracer
	defer closer.Close()
	tracer := opentracing.GlobalTracer()

	spanA := tracer.StartSpan("Span A")
	defer spanA.Finish()

	spanB := opentracing.GlobalTracer().StartSpan("span B", opentracing.ChildOf(spanA.Context()))
	defer spanB.Finish()
	childSpanB(spanB)

	spanC := opentracing.GlobalTracer().StartSpan("span C", opentracing.ChildOf(spanA.Context()))
	defer spanC.Finish()
	childSpanC(spanC)

	spanA.LogFields(traceLog.String("info", "finish business"))
}

func childSpanB(spanB opentracing.Span) {
	time.Sleep(500 * time.Millisecond)
	spanD := opentracing.GlobalTracer().StartSpan("span D", opentracing.ChildOf(spanB.Context()))
	defer spanD.Finish()
	spanD.SetTag("name", "span D")
}

func childSpanC(spanC opentracing.Span) {
	spanE := opentracing.GlobalTracer().StartSpan("span E", opentracing.ChildOf(spanC.Context()))
	defer spanE.Finish()
	spanE.SetTag("name", "span E")
	time.Sleep(100 * time.Millisecond)

	spanF := opentracing.GlobalTracer().StartSpan("span F", opentracing.ChildOf(spanC.Context()))
	defer spanF.Finish()
	spanF.SetTag("name", "span F")
	time.Sleep(100 * time.Millisecond)

	spanG := opentracing.GlobalTracer().StartSpan("span G", opentracing.FollowsFrom(spanF.Context()))
	defer spanG.Finish()
	spanG.SetTag("name", "span G")
	time.Sleep(100 * time.Millisecond)

	spanH := opentracing.GlobalTracer().StartSpan("span H", opentracing.FollowsFrom(spanG.Context()))
	defer spanH.Finish()
	spanH.SetTag("name", "span H")
	time.Sleep(200 * time.Millisecond)
}
