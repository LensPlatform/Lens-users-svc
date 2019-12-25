package middleware

import (
	"net/http"
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type CircuitBreaker struct {
	CircuitBreaker *gobreaker.CircuitBreaker
	Logger zap.Logger
}

func NewCircuitBreaker(name string, maxRequest uint32, timeout time.Duration, interval time.Duration, readyToTrip func(counts gobreaker.Counts) bool, logger zap.Logger ) *CircuitBreaker{
	var st gobreaker.Settings
	st.Name = name
	st.MaxRequests = maxRequest
	st.Timeout = timeout
	st.Interval = interval

	if readyToTrip == nil {
		st.ReadyToTrip = defaultReadyToTrip
	} else {
		st.ReadyToTrip = readyToTrip
	}
	return &CircuitBreaker{gobreaker.NewCircuitBreaker(st), logger}
}

func defaultReadyToTrip(counts gobreaker.Counts) bool{
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.ConsecutiveFailures > 4 || failureRatio > 0.5
}


func (breaker CircuitBreaker) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		_, err := breaker.CircuitBreaker.Execute(func() (interface{}, error) {
			next.ServeHTTP(w, r)
			return nil, nil
		})

		if err != nil {
			breaker.Logger.Error(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	})
}

