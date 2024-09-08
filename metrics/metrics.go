package metrics

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	serviceName string
)

var (
	requestCount       = make(map[string]int)
	requestDuration    = make(map[string]time.Duration)
	requestStatusCodes = make(map[string]map[int]int)
	mu                 sync.Mutex
)

type MetricInit struct {
	ServiceName string
}

// recordMetrics는 메트릭을 기록합니다.
func recordMetrics(method, endpoint string, duration time.Duration, status int, serviceName string) {
	mu.Lock()
	defer mu.Unlock()

	// 요청 수 증가
	key := method + endpoint
	requestCount[key]++

	// 요청 지연 시간 추가
	requestDuration[key] += duration

	// 상태 코드 수 증가
	if _, exists := requestStatusCodes[key]; !exists {
		requestStatusCodes[key] = make(map[int]int)
	}
	requestStatusCodes[key][status]++
}

// MetricsHandler는 메트릭을 제공하고 초기화합니다.
func MetricsHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	metricsData := make(map[string]interface{})
	for endpoint, count := range requestCount {
		avgDuration := float64(requestDuration[endpoint].Milliseconds()) / float64(count)
		metricsData[endpoint] = map[string]interface{}{
			"count":        count,
			"avg_duration": avgDuration,
			"status_codes": requestStatusCodes[endpoint],
		}
	}

	response := map[string]interface{}{
		"service_name": serviceName,
		"metrics":      metricsData,
	}

	c.JSON(200, response)
}

// MetricsMiddleware는 서비스 이름을 인자로 받는 미들웨어를 생성합니다.
func MetricsMiddleware(init MetricInit) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName = init.ServiceName

		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "/unknown"
		}

		recordMetrics(c.Request.Method, endpoint, duration, status, serviceName)
	}
}
