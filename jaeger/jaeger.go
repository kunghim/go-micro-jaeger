package jaeger

import (
	"errors"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

type OpenTrace struct {
	traceName string
	agentAddr string
	Tracer    opentracing.Tracer
	Closer    io.Closer
}

func Tracer(traceName, agentAddr string) *OpenTrace {
	if len(traceName) == 0 || len(agentAddr) == 0 {
		log.Fatal("traceName 或 agentAddr 不能为空")
		return nil
	}
	return &OpenTrace{traceName: traceName, agentAddr: agentAddr}
}

func (o *OpenTrace) Create() (*OpenTrace, error) {
	if o == nil {
		return nil, errors.New("无法连接 jaeger agent")
	}
	tracer, closer, err := o.newTracer(o.traceName, o.agentAddr)
	if err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	o.Closer = closer
	o.Tracer = tracer
	return o, nil
}

func (o *OpenTrace) newTracer(traceName string, agentAddr string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: traceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(agentAddr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		config.Reporter(reporter),
	)

	return tracer, closer, err
}
