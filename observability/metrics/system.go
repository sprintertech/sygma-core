package metrics

import (
	"context"
	"runtime"
	"runtime/debug"

	"go.opentelemetry.io/otel/metric"
)

type SystemMetrics struct {
	goRoutinesGauge     metric.Int64ObservableGauge
	gcDurationHistogram metric.Float64Histogram
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

	gcDurationHistogram, err := meter.Float64Histogram(
		"relayer.GcDurationSeconds",
		metric.WithFloat64Callback(func(context context.Context, result metric.Float64Observer) error {
			var gcStats debug.GCStats
			debug.ReadGCStats(&gcStats)
			if len(gcStats.Pause) > 0 {
				recentPauseDuration := gcStats.Pause[0].Seconds()
				result.Observe(recentPauseDuration, opts)
			}
		}),
		metric.WithDescription("Duration of garbage collection cycles."),
	)
	if err != nil {
		return nil, err
	}

	return &SystemMetrics{
		goRoutinesGauge:     goRoutinesGauge,
		gcDurationHistogram: gcDurationHistogram,
	}, err
}
