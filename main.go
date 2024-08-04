package main

import (
	"gin-monitoring/metrics"
	"gin-monitoring/tracing"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 미들웨어 추가
	r.Use(metrics.MetricsMiddleware())
	r.Use(tracing.TracingMiddleware())

	// /metrics 엔드포인트 추가
	r.GET("/metrics", metrics.MetricsHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // 기본 포트 8080에서 실행
}
