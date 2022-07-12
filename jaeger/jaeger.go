package jaeger

import (
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	cons "go-micro-jaeger/constant"
	"io"
	"time"
)

func NewTracer(traceName string) (io.Closer, error) {
	if len(traceName) == 0 {
		return nil, errors.New("trace name can not be null")
	}
	tracer, closer, err := newTracer(traceName, cons.JaegerAgent)
	opentracing.SetGlobalTracer(tracer)
	return closer, err
}

func newTracer(traceName string, agentAddr string) (opentracing.Tracer, io.Closer, error) {
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
