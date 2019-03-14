package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kklinan/jaeger-tracing/app/gin/controller/health"
	"github.com/kklinan/jaeger-tracing/app/gin/controller/user"
)

// InitRouter Init Router
func InitRouter(r *gin.Engine) {
	r.GET("/ping", health.Ping)
	r.GET("/user/:id", user.Get)
	r.POST("/user", user.Post)
}
