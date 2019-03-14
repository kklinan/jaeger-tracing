package jaeger

import (
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	// "github.com/uber/jaeger-lib/metrics/prometheus"
)

// 初始化变量
var (
	err    error
	Tracer opentracing.Tracer
	Closer io.Closer
)

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(service string) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	// metricsFactory := prometheus.New()
	// Tracer, Closer, err = cfg.NewTracer(
	// 	config.Metrics(metricsFactory),
	// )

	Tracer, Closer, err = cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(Tracer)
	// return Tracer, Closer
}
