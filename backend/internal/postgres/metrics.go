package postgres

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const meterName = "github.com/arowden/augment-fund/internal/postgres"

func RegisterMetrics(pool *Pool) error {
	meter := otel.Meter(meterName)

	_, err := meter.Int64ObservableGauge(
		"db_pool_size",
		metric.WithDescription("Maximum number of connections in the pool"),
		metric.WithUnit("{connections}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(pool.Stat().MaxConns()))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	_, err = meter.Int64ObservableGauge(
		"db_pool_active_connections",
		metric.WithDescription("Number of connections currently in use"),
		metric.WithUnit("{connections}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(pool.Stat().AcquiredConns()))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	_, err = meter.Int64ObservableGauge(
		"db_pool_idle_connections",
		metric.WithDescription("Number of idle connections in the pool"),
		metric.WithUnit("{connections}"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(pool.Stat().IdleConns()))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	return nil
}
