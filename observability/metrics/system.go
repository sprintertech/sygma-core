package metrics

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

type SystemMetrics struct {
	goRoutinesGauge metric.Int64ObservableGauge
}

// NewSystemMetrics initializes system performance and resource utilization metrics
func NewSystemMetrics(meter metric.Meter, opts metric.MeasurementOption) (*SystemMetrics, error) {
	goRoutinesGauge, err := meter.Int64ObservableGauge(
		"relayer.GoRoutines",
		metric.WithInt64Callback(func(context context.Context, result metric.Int64Observer) error {
			result.Observe(int64(runtime.NumGoroutine()), opts)
			return nil
		}),
		metric.WithDescription("Number of Go routines running."),
	)
	if err != nil {
		return nil, err
	}

	return &SystemMetrics{
		goRoutinesGauge: goRoutinesGauge,
	}, err
}
