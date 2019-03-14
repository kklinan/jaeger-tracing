package main

import (
	"context"
	"io"
	"io/ioutil"

	"net/http"
	"net/url"
	"os"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

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

	span := Tracer.StartSpan("example-4")
	span.SetTag("example", "4")
	defer span.Finish()
	span.SetBaggageItem("username", helloTo)

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	formatString(ctx, helloTo)
	printHello(ctx, helloTo)
	span.LogKV("username", helloTo)
}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	span.SetTag("example", "4")
	defer span.Finish()

	v := url.Values{}
	v.Set("hello_to", helloTo)
	url := "http://localhost:8091/format?" + v.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	req = req.WithContext(ctx)

	req, ht := nethttp.TraceRequest(span.Tracer(), req)
	defer ht.Finish()

	client := &http.Client{Transport: &nethttp.Transport{}}

	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	helloStr := string(body)

	span.LogFields(
		log.String("exampleNo", "4"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	span.SetTag("example", "4")
	defer span.Finish()

	v := url.Values{}
	v.Set("hello_to", helloStr)
	url := "http://localhost:8092/publish?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	req = req.WithContext(ctx)

	req, ht := nethttp.TraceRequest(span.Tracer(), req)
	defer ht.Finish()

	client := &http.Client{Transport: &nethttp.Transport{}}
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	res.Body.Close()

	println(helloStr)
	span.LogKV("print", helloStr)
}
