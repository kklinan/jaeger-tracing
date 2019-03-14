package main

import (
	"context"
	"fmt"
	"io"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

// Tracer global tracer
var Tracer opentracing.Tracer

// Closer global closer
var Closer io.Closer

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}
	helloTo := os.Args[1]

	Tracer, Closer := initJaeger("example")
	defer Closer.Close()
	opentracing.SetGlobalTracer(Tracer)

	span := Tracer.StartSpan("example-2")
	span.SetTag("example", "2")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)
	span.LogKV("exampleNo", "2")
}

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

// 传递 span 方式
// func formatString(rootSpan opentracing.Span, helloTo string) string {
// 	span := rootSpan.Tracer().StartSpan("formatString")
// 	defer span.Finish()

// 	span.SetTag("example", "2")

// 	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
// 	span.LogFields(
// 		log.String("exampleNo", "2"),
// 		log.String("value", helloStr),
// 	)

// 	return helloStr
// }

// 传递上下文方式
func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	span.SetTag("example", "2")
	defer span.Finish()

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("exampleNo", "2"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	span.SetTag("example", "2")
	defer span.Finish()

	println(helloStr)
	span.LogKV("print", helloStr)
}
