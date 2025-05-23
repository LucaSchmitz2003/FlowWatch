package otelHelper

import (
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	shutdownFuncs []func() error
	once          sync.Once
)

func Shutdown() {
	for _, shutdown := range shutdownFuncs {
		err := shutdown()
		if err != nil {
			log.Printf("Failed to shut down the service. %v", err)
		}
	}
}

// initOtelHelper initializes the trace-, metric- & log-provider.
func initOtelHelper() {
	// Set the global text map propagator
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	// Load the environment variables to make sure that the settings have already been loaded
	_ = godotenv.Load(".env")

	// Get the service name from the environment variables
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "TestService"
		log.Println("OTEL_SERVICE_NAME not set, using default")
	}

	// Get the collector URL from the environment variables
	collectorURL := os.Getenv("OTEL_COLLECTOR_URL")
	if collectorURL == "" {
		log.Println("OTEL_COLLECTOR_URL not set, trace export will be skipped")
	}

	// Get the tls support state from the environment variables
	supportTLS, err := strconv.ParseBool(os.Getenv("OTEL_SUPPORT_TLS"))
	if err != nil {
		supportTLS = false
		log.Printf("Failed to parse OTEL_SUPPORT_TLS, using default. %v", err)
	}

	// Initialize the trace provider
	err = initTraceProvider(serviceName, collectorURL, supportTLS)
	if err != nil {
		log.Fatalf("Failed to set up the trace provider. %v", err)
	}
}

// SetupOtelHelper initializes the OpenTelemetry SDK connection to the backend if it has not been initialized yet according to the singleton pattern.
func SetupOtelHelper() {
	// Create a new LogHelper instance if it does not exist
	once.Do(func() {
		initOtelHelper()
	})
}
