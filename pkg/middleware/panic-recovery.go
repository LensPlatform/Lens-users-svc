package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type PanicRecovery struct {
	logger zap.Logger
}

func NewPanicRecovery(logger zap.Logger) *PanicRecovery {
	return &PanicRecovery{logger:logger}
}

func (panic *PanicRecovery) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				panic.logger.Error("Panic Occured", zap.Any("Panic", err))
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}



