package metrics

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	requestCount       = make(map[string]int)
	requestDuration    = make(map[string]time.Duration)
	requestStatusCodes = make(map[string]map[int]int)
	mu                 sync.Mutex
)

// 메트릭을 기록하는 함수
func recordMetrics(method, endpoint string, duration time.Duration, status int) {
	mu.Lock()
	defer mu.Unlock()

	// 요청 수 증가
	requestCount[method+endpoint]++

	// 요청 지연 시간 추가
	requestDuration[method+endpoint] += duration

	// 상태 코드 수 증가
	if _, exists := requestStatusCodes[method+endpoint]; !exists {
		requestStatusCodes[method+endpoint] = make(map[int]int)
	}
	requestStatusCodes[method+endpoint][status]++
}

// 메트릭 핸들러
func MetricsHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	metricsData := make(map[string]interface{})
	for endpoint, count := range requestCount {
		avgDuration := requestDuration[endpoint].Seconds() / float64(count)
		metricsData[endpoint] = map[string]interface{}{
			"count":        count,
			"avg_duration": avgDuration,
			"status_codes": requestStatusCodes[endpoint],
		}
	}

	c.JSON(200, metricsData)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "unknown"
		}

		recordMetrics(c.Request.Method, endpoint, duration, status)
	}
}
