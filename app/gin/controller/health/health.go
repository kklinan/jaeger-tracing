package health

import (
	"github.com/gin-gonic/gin"

	"github.com/opentracing/opentracing-go/log"

	"github.com/kklinan/jaeger-tracing/app/gin/jaeger"
)

// Ping Ping
func Ping(c *gin.Context) {
	// time.Sleep(time.Second * 5)

	span := jaeger.Tracer.StartSpan("ping")
	span.SetTag("health", "ping")
	defer span.Finish()

	c.JSON(200, gin.H{
		"message": "pong",
	})

	span.LogFields(
		log.String("ip", "0.0.0.0"),
		log.String("server_name", "vipmember"),
	)
	span.LogKV("ip", "0.0.0.0")
	span.LogKV("server_name", "vipmember")
}
