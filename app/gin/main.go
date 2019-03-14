package main

import (
	"github.com/gin-gonic/gin"

	"github.com/kklinan/jaeger-tracing/app/gin/jaeger"
	"github.com/kklinan/jaeger-tracing/app/gin/router"
)

func main() {
	jaeger.Init("gin-server")
	defer jaeger.Closer.Close()

	r := gin.Default()
	router.InitRouter(r)

	r.Run(":8089") // listen and serve on 0.0.0.0:8080
}
