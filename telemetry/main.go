package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type loggerContextKey struct{}

func NewLogger(serviceName string) *slog.Logger {
	return slog.New(slog.
		NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {

				switch a.Value.Kind() {
				case slog.KindAny:
					switch v := a.Value.Any().(type) {
					case *http.Request:
						return slog.Group(a.Key,
							slog.String("method", v.Method),
							slog.String("path", v.URL.String()))
					}
				}

				return a
			},
		}).
		WithAttrs([]slog.Attr{
			slog.String("service_name", serviceName),
		}))
}

func StartTracerProvider(ctx context.Context, logger *slog.Logger, exportTraces bool) (func(), error) {
	resource, err := ecs.NewResourceDetector().Detect(ctx)
	if err != nil {
		return nil, err
	}

	var traceExporter trace.SpanExporter
	if exportTraces {
		traceExporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("0.0.0.0:4317"),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		)

		if err != nil {
			return nil, err
		}
	}

	idg := xray.NewIDGenerator()
	tp := trace.NewTracerProvider(
		trace.WithResource(resource),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
		trace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("tracer shutdown error: %s", err.Error()))
		}
	}, nil
}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

func Middleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := oteltrace.SpanFromContext(r.Context())

			span.SetAttributes(
				attribute.String("http.target", r.URL.Path),
			)

			loggerWithRequest := logger.With(
				slog.String("trace_id", span.SpanContext().TraceID().String()),
				slog.Any("request", r),
			)

			r = r.WithContext(ContextWithLogger(r.Context(), loggerWithRequest))

			next.ServeHTTP(w, r)
		})
	}
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	return ctx.Value(loggerContextKey{}).(*slog.Logger)
}
