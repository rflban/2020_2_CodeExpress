package monitoring

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Monitoring struct {
	Hits     *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

func NewMonitoring(e *echo.Echo) *Monitoring {
	hits := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
		Help: "help",
	}, []string{"status", "path", "method"})

	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "duration",
		Help: "help",
	}, []string{"status", "path", "method"})

	var monitoring = &Monitoring{
		Hits:     hits,
		Duration: duration,
	}

	prometheus.MustRegister(monitoring.Hits, monitoring.Duration)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	return monitoring
}
