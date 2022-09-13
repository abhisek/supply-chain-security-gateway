package obs

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerControlEnvKey         = "APP_SERVICE_OBS_ENABLED"
	tracerServiceNameEnvKey     = "APP_SERVICE_NAME"
	tracerServiceEnvEnvKey      = "APP_SERVICE_ENV"
	tracerServiceLabelEnvKey    = "APP_SERVICE_LABELS"
	tracerOtelExporterUrlEnvKey = "APP_OTEL_EXPORTER_OTLP_ENDPOINT"
)

var (
	globalTracer = otel.Tracer("NOP")
)

func InitTracing() func(context.Context) error {
	if !isTracingEnabled() {
		logger.Infof("Tracer is disabled")
		return func(ctx context.Context) error { return nil }
	}

	serviceName := os.Getenv(tracerServiceNameEnvKey)
	serviceEnv := os.Getenv(tracerServiceEnvEnvKey)
	otlpExporterUrl := os.Getenv(tracerOtelExporterUrlEnvKey)

	if utils.IsEmptyString(serviceName) || utils.IsEmptyString(otlpExporterUrl) {
		panic("tracer is enable but required environment is not defined")
	}

	logger.With(map[string]any{
		"ServiceName":     serviceName,
		"ServiceEnv":      serviceEnv,
		"OtlpExporterUrl": otlpExporterUrl,
	}).Infof("Enabling tracer")

	// NOTE: We expect the collector to be a sidecar
	// TODO: Revisit this for using a secure channel
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(otlpExporterUrl),
		),
	)

	if err != nil {
		panic(fmt.Sprintf("error creating otlp exporter: %v", err))
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.environment", serviceEnv),
			attribute.String("service.language", "go"),
		),
	)

	if err != nil {
		panic(fmt.Sprintf("error creating otlp resource: %v", err))
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	globalTracer = otel.Tracer(serviceName)

	logger.Infof("Global tracer enabled")
	return exporter.Shutdown
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

func LoggerTags(ctx context.Context) map[string]any {
	tags := map[string]any{}
	span := trace.SpanFromContext(ctx)

	if span.IsRecording() {
		tags["span_id"] = span.SpanContext().SpanID()
		tags["trace_id"] = span.SpanContext().TraceID()
		tags["trace_flags"] = span.SpanContext().TraceFlags()
	}

	return tags
}

func isTracingEnabled() bool {
	bRet, err := strconv.ParseBool(os.Getenv(tracerControlEnvKey))
	if err != nil {
		return false
	}

	return bRet
}
