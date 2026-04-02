package middleware

import (
	"managify/internal/metrics"
	"time"

	"github.com/gofiber/fiber/v2"
)

func MetricMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)

	key := c.Method() + " " + c.Path()

	m, exists := metrics.Metrics[key]
	if !exists {
		m = &metrics.EndpointMetrics{
			MinTime: duration,
			MaxTime: duration,
		}
		metrics.Metrics[key] = m
	}

	m.Count++
	m.TotalTime += duration
	if duration < m.MinTime {
		m.MinTime = duration
	}
	if duration > m.MaxTime {
		m.MaxTime = duration
	}

	return err
}
