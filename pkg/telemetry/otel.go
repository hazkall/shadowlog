package telemetry

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	Tracer trace.Tracer
	Meter  metric.Meter
)

func newOTLPGRPCTraceExporter(ctx context.Context) (oteltrace.SpanExporter, error) {
	return otlptracegrpc.New(ctx)
}

func newOTLPGRPCMetricExporter(ctx context.Context) (*otlpmetricgrpc.Exporter, error) {
	return otlpmetricgrpc.New(ctx)
}

func newTraceProvider(ctx context.Context, exp oteltrace.SpanExporter) *oteltrace.TracerProvider {
	r := getOTLPResource(ctx)

	return oteltrace.NewTracerProvider(
		oteltrace.WithResource(r),
		oteltrace.WithSpanProcessor(oteltrace.NewBatchSpanProcessor(exp)),
		oteltrace.WithSampler(oteltrace.AlwaysSample()),
	)
}

func newMeterProvider(ctx context.Context, exp otelmetric.Exporter) *otelmetric.MeterProvider {

	r := getOTLPResource(ctx)

	return otelmetric.NewMeterProvider(
		otelmetric.WithResource(r),
		otelmetric.WithReader(otelmetric.NewPeriodicReader(exp)),
	)

}

func GenerateCommonAttributes(attrs ...attribute.KeyValue) []attribute.KeyValue {
	return attrs
}

func getOTLPResource(ctx context.Context) *resource.Resource {
	r, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithSchemaURL(
			semconv.SchemaURL,
		),
	)

	if err != nil {
		panic(err)
	}

	return r

}

func getOTLPPropagators() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{},
		propagation.Baggage{},
		jaeger.Jaeger{})
}

func TraceInit(ctx context.Context, name string) oteltrace.SpanExporter {
	exp, err := newOTLPGRPCTraceExporter(ctx)

	if err != nil {
		panic(err)
	}

	tp := newTraceProvider(ctx, exp)

	otel.SetTextMapPropagator(getOTLPPropagators())

	Tracer = tp.Tracer(name)

	return exp
}

func MetricInit(ctx context.Context, name string) (*otelmetric.MeterProvider, otelmetric.Exporter) {
	exp, err := newOTLPGRPCMetricExporter(ctx)

	if err != nil {
		panic(err)
	}

	mp := newMeterProvider(ctx, exp)

	otel.SetMeterProvider(mp)

	otel.SetTextMapPropagator(getOTLPPropagators())

	Meter = mp.Meter(name)

	return mp, exp
}

func RuntimeStart(m metric.MeterProvider) {
	if err := runtime.Start(
		runtime.WithMeterProvider(m),
	); err != nil {
		panic(err)
	}
}
