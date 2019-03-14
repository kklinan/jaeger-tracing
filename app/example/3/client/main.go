package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/http"

	"github.com/kklinan/jaeger-tracing/app/example/3/base"
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

	Tracer, Closer := base.InitJaeger("example")
	defer Closer.Close()
	opentracing.SetGlobalTracer(Tracer)

	span := Tracer.StartSpan("example-3")
	span.SetTag("example", "3")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)
	span.LogKV("exampleNo", "3")
}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	span.SetTag("example", "3")
	defer span.Finish()

	v := url.Values{}
	v.Set("hello_to", helloTo)
	url := "http://localhost:8081/format?" + v.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := xhttp.Do(req)
	if err != nil {
		panic(err.Error())
	}

	helloStr := string(resp)

	span.LogFields(
		log.String("exampleNo", "3"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	span.SetTag("example", "3")
	defer span.Finish()

	v := url.Values{}
	v.Set("hello_to", helloStr)
	url := "http://localhost:8082/publish?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	if _, err := xhttp.Do(req); err != nil {
		panic(err.Error())
	}

	println(helloStr)
	span.LogKV("print", helloStr)
}
