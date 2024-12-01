package otelHelper

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"log"
)

func initTraceProvider(serviceName, collectorURL string, supportTLS bool) error {
	// Create a slice to hold the exporter options
	var opts []otlptracegrpc.Option

	// Add the collector URL to the exporter options
	opts = append(opts, otlptracegrpc.WithEndpoint(collectorURL))

	// If the connection is insecure, add the insecure option to the exporter options
	if !supportTLS { // Thanks to Levin for pointing out the missing exclamation mark
		opts = append(opts, otlptracegrpc.WithInsecure())
		log.Println("Insecure connection to the collector")
	} else {
		log.Fatal("TLS is not implemented yet")
		// TODO: Implement TLS connection
	}

	// Create a slice to hold the trace provider options
	var tpOptions []trace.TracerProviderOption

	// Create an OTLP trace exporter
	sigNozTraceExporter, err := otlptracegrpc.New(context.Background(), opts...)
	if err != nil {
		err = errors.Wrap(err, "Failed to create OTLP exporter")
		return err
	}
	tpOptions = append(tpOptions, trace.WithBatcher(sigNozTraceExporter))

	// Set the service name
	tpOptions = append(tpOptions, trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))))

	// Create a new trace provider with the configured options
	tp := trace.NewTracerProvider(tpOptions...)

	// Set the trace provider to the global provider
	otel.SetTracerProvider(tp)

	// Add the shutdown function to the global slice
	shutdown := func() error {
		// Shutdown the tracer provider to flush any remaining spans
		err1 := tp.Shutdown(context.Background())
		if err1 != nil {
			err1 = errors.Wrap(err1, "Failed to shut down the tracer provider.")
		}

		// Shutdown the SigNoz exporter to ensure all spans are sent
		err2 := sigNozTraceExporter.Shutdown(context.Background())
		if err2 != nil {
			err2 = errors.Wrap(err2, "Failed to shut down the SigNoz exporter.")
		}

		if err1 != nil && err2 != nil {
			err := errors.Wrap(err1, err2.Error())
			return err
		} else if err1 != nil {
			return err1
		}

		return err2
	}

	shutdownFuncs = append(shutdownFuncs, shutdown)

	return nil
}
