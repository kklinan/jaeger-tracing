package user

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"

	"github.com/kklinan/jaeger-tracing/app/gin/jaeger"
)

// Post Create user
func Post(c *gin.Context) {
	span := jaeger.Tracer.StartSpan("create_user")
	defer span.Finish()

	user := gin.H{
		"user_id":   1,
		"user_name": "张三",
		"sex":       1,
		"addr":      "北京东城区",
	}
	c.JSON(200, user)

	userStr, _ := json.Marshal(user)
	span.LogFields(
		log.String("user", string(userStr)),
	)
}
