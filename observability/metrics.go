package observability

import (
	"context"
	"net/url"

	"github.com/sygmaprotocol/sygma-core/observability/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	api "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initResource() *sdkresource.Resource {
	res, _ := sdkresource.New(context.Background(),
		sdkresource.WithProcess(),
		sdkresource.WithTelemetrySDK(),
		sdkresource.WithHost(),
		sdkresource.WithAttributes(
			semconv.ServiceName("relayer"),
		),
	)
	return res
}

func InitMetricProvider(ctx context.Context, agentURL string) (*sdkmetric.MeterProvider, error) {
	collectorURL, err := url.Parse(agentURL)
	if err != nil {
		return nil, err
	}

	metricOptions := []otlpmetrichttp.Option{
		otlpmetrichttp.WithURLPath(collectorURL.Path),
		otlpmetrichttp.WithEndpoint(collectorURL.Host),
	}
	if collectorURL.Scheme == "http" {
		metricOptions = append(metricOptions, otlpmetrichttp.WithInsecure())
	}

	metricHTTPExporter, err := otlpmetrichttp.New(ctx, metricOptions...)
	if err != nil {
		return nil, err
	}

	httpMetricReader := sdkmetric.NewPeriodicReader(metricHTTPExporter)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(httpMetricReader),
		sdkmetric.WithResource(initResource()),
	)

	return meterProvider, nil
}

type RelayerMetrics struct {
	*metrics.SystemMetrics
	*metrics.NetworkMetrics

	Opts api.MeasurementOption
}

// NewRelayerMetrics initializes OpenTelemetry metrics
func NewRelayerMetrics(ctx context.Context, meter metric.Meter, attributes ...attribute.KeyValue) (*RelayerMetrics, error) {
	opts := api.WithAttributes(attributes...)

	systemMetrics, err := metrics.NewSystemMetrics(ctx, meter, opts)
	if err != nil {
		return nil, err
	}

	networkMetrics, err := metrics.NewNetworkMetrics(ctx, meter, opts)
	if err != nil {
		return nil, err
	}

	return &RelayerMetrics{
		SystemMetrics:  systemMetrics,
		NetworkMetrics: networkMetrics,
		Opts:           opts,
	}, err
}
