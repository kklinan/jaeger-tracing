package main

import (
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"

	"github.com/kklinan/jaeger-tracing/app/example/3/base"
)

func main() {
	tracer, closer := base.InitJaeger("publisher")
	defer closer.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("publish", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloTo := r.FormValue("hello_to")
		println("-->", helloTo)

		span.LogFields(
			otlog.String("event", "string-publish"),
			otlog.String("value", helloTo),
		)
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
