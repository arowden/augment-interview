package otel

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func newSampler(sampleRate float64) sdktrace.Sampler {
	if sampleRate < 0 {
		sampleRate = 0
	}
	if sampleRate > 1 {
		sampleRate = 1
	}

	return sdktrace.ParentBased(
		sdktrace.TraceIDRatioBased(sampleRate),
		sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
		sdktrace.WithRemoteParentNotSampled(sdktrace.TraceIDRatioBased(sampleRate)),
		sdktrace.WithLocalParentSampled(sdktrace.AlwaysSample()),
		sdktrace.WithLocalParentNotSampled(sdktrace.TraceIDRatioBased(sampleRate)),
	)
}
