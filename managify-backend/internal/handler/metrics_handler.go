package handler

import (
	"fmt"
	"managify/internal/metrics"
	"time"

	"github.com/gofiber/fiber/v2"
)

func MetricsHandler(c *fiber.Ctx) error {
	c.Type("html", "utf-8")

	html := `
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
			table { border-collapse: collapse; width: 100%; background-color: #fff; }
			th, td { border: 1px solid #ccc; padding: 8px; text-align: center; }
			th { background-color: #007bff; color: white; }
			tr:nth-child(even) { background-color: #f2f2f2; }
			tr:hover { background-color: #e0e0e0; }
		</style>
	</head>
	<body>
		<h2>Endpoint Metrics</h2>
		<table>
			<tr>
				<th>Endpoint</th>
				<th>Count</th>
				<th>Total</th>
				<th>Min</th>
				<th>Max</th>
				<th>Avg</th>
			</tr>
	`

	for key, m := range metrics.Metrics {
		avgTime := time.Duration(0)
		if m.Count > 0 {
			avgTime = m.TotalTime / time.Duration(m.Count)
		}

		formatDuration := func(d time.Duration) string {
			if d < time.Millisecond {
				return fmt.Sprintf("%.2f µs", float64(d)/1000)
			}
			return fmt.Sprintf("%.2f ms", float64(d.Microseconds())/1000)
		}

		html += fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			key,
			m.Count,
			formatDuration(m.TotalTime),
			formatDuration(m.MinTime),
			formatDuration(m.MaxTime),
			formatDuration(avgTime),
		)
	}

	html += "</table></body></html>"

	return c.SendString(html)
}
