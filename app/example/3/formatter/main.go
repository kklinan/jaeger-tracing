package main

import (
	"fmt"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"

	"github.com/kklinan/jaeger-tracing/app/example/3/base"
)

func main() {
	tracer, closer := base.InitJaeger("formatter")
	defer closer.Close()

	http.HandleFunc("/format", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("format", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloTo := r.FormValue("hello_to")
		helloStr := fmt.Sprintf("Hello %s", helloTo)
		w.Write([]byte(helloStr))

		span.LogFields(
			otlog.String("event", "string-format"),
			otlog.String("value", helloStr),
		)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
