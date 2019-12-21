package middleware

import (
	"net/http"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
)

type ZipkinTracer struct {
	OperationName string
	Tracer *zipkin.Tracer
}

func NewZipKinTracerMiddleware(operationName string, tracer *zipkin.Tracer) *ZipkinTracer{
	return &ZipkinTracer{OperationName:operationName, Tracer:tracer}
}

func (tracer *ZipkinTracer) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sc model.SpanContext

		if parentSpan := zipkin.SpanFromContext(r.Context()); parentSpan != nil {
			sc = parentSpan.Context()
		}
		sp := tracer.Tracer.StartSpan(tracer.OperationName, zipkin.Parent(sc))
		defer sp.Finish()

		r.WithContext(zipkin.NewContext(r.Context(), sp))
		next.ServeHTTP(w, r)
	})
}