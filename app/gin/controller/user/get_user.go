package user

import (
	"github.com/gin-gonic/gin"

	"github.com/kklinan/jaeger-tracing/app/gin/jaeger"
)

// Get Get user by userID
func Get(c *gin.Context) {
	span := jaeger.Tracer.StartSpan("get_user")
	defer span.Finish()

	id := c.Param("id")
	c.JSON(200, gin.H{
		"user_id":   id,
		"user_name": "张三",
		"sex":       1,
		"addr":      "北京东城区",
	})
}
