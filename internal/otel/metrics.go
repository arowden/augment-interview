package otel

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const meterName = "github.com/arowden/augment-fund/internal/otel"

var (
	fundCreatedTotal      metric.Int64Counter
	transferExecutedTotal metric.Int64Counter
	transferUnitsTotal    metric.Int64Counter

	httpRequestDuration metric.Float64Histogram

	initOnce sync.Once
	initErr  error
)

func InitMetrics() error {
	initOnce.Do(func() {
		initErr = initMetricsInternal()
	})
	return initErr
}

func initMetricsInternal() error {
	meter := otel.Meter(meterName)
	var err error

	fundCreatedTotal, err = meter.Int64Counter(
		"fund_created_total",
		metric.WithDescription("Total number of funds created"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	transferExecutedTotal, err = meter.Int64Counter(
		"transfer_executed_total",
		metric.WithDescription("Total number of transfers executed"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	transferUnitsTotal, err = meter.Int64Counter(
		"transfer_units_total",
		metric.WithDescription("Total units transferred"),
		metric.WithUnit("{units}"),
	)
	if err != nil {
		return err
	}

	httpRequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request latency in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(
			0.001, 0.005, 0.01, 0.025, 0.05,
			0.1, 0.25, 0.5, 1, 2.5, 5, 10,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func RecordFundCreated(ctx context.Context) {
	if fundCreatedTotal != nil {
		fundCreatedTotal.Add(ctx, 1)
	}
}

func RecordTransferExecuted(ctx context.Context, units int64, status string) {
	attrs := metric.WithAttributes(attribute.String("status", status))

	if transferExecutedTotal != nil {
		transferExecutedTotal.Add(ctx, 1, attrs)
	}
	if transferUnitsTotal != nil && units > 0 {
		transferUnitsTotal.Add(ctx, units, attrs)
	}
}

func RecordHTTPRequestDuration(ctx context.Context, duration float64, method, path string, statusCode int) {
	if httpRequestDuration != nil {
		httpRequestDuration.Record(ctx, duration,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int("status_code", statusCode),
			),
		)
	}
}
