package obs

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

var (
	globalTracer = otel.Tracer("NOP")
)

func InitTracing() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func ShutdownTracing() {
	// Explicitly flush and shutdown tracers
}

func Spanned(current context.Context, name string,
	f func(context.Context) error) error {
	newCtx, span := globalTracer.Start(current, name)
	defer span.End()

	err := f(newCtx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
