package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moolya-ai/go-trace-sdk/pkg/trace"
)

func main() {
	trace.InitLogger("http://localhost:3000/logs", "<YOUR_API_KEY>")

	router := gin.Default();
	router.Use(trace.GinTraceMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	router.PUT("/test", func(c *gin.Context) {
		c.JSON(401, gin.H{"message": "Hello, World!"})
	})

	router.DELETE("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	router.PATCH("/test", func(c *gin.Context) {
		c.JSON(400, gin.H{"message": "Hello, World!"})
	})

	router.Run(":4000")
}
