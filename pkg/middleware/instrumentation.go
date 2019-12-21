package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type InstrumentationMiddleware struct {
	duration metrics.Histogram
}

func NewInstrumentationMiddleware(operationName string) *InstrumentationMiddleware {
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "lens_users_svc",
			Subsystem: "lens_users_svc",
			Name:      operationName,
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}

	return &InstrumentationMiddleware{
		duration.With("method", operationName),
	}
}

func (m *InstrumentationMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer func(begin time.Time) {
			m.duration.With("success", fmt.Sprint(r.Context().Err() == nil)).Observe(time.Since(begin).Seconds())
		}(time.Now())

		next.ServeHTTP(w, r)
	})
}