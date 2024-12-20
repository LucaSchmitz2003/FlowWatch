# FlowWatch

FlowWatch is an abstraction layer for OpenTelemetry, standardizing tracing, logging (both local and remote), metrics, and exception handling across the company's ecosystem programs. Its goal is to simplify implementation and ensure consistent usage.

---

## 1. OpenTelemetry Integration

### Setup
To set up OpenTelemetry, initialize it at the start of your program:
```go
otelHelper.SetupOtelHelper()
defer otelHelper.Shutdown() // Recommended: Graceful shutdown at program end
```

### Tracing
To start a trace, use the following methods:
```go
tracer := otel.Tracer("TracerName")          // Initialize the tracer
ctx, span := tracer.Start(ctx, "SpanName")  // Start a new span
defer span.End()                            // End the span (defer recommended)
```

> **Note:** Use the updated context `ctx` in all subsequent operations to ensure that logs and spans are properly associated.

---

## 2. Logging

### Example
To log messages, retrieve the logging helper and use the appropriate log level:
```go
lh := loggingHelper.GetLogHelper()
lh.Info(ctx, "Info log message")
lh.Warn(ctx, "Warning log message")
```

> **Note:** Supported log levels are `Debug`, `Info`, `Warn`, `Error`, and `Fatal`.

---

## 3. Exception Handling

- **Recommendation:** Use `pkg/errors` for creating and wrapping errors:
```go
err := errors.Wrap(CustomError1, "Additional context")
```

- **Global Variables:** Declare errors as global variables to ensure consistent error messages:
```go
var CustomError1 = errors.New("Error message")
```

### Logging an Error
```go
if err != nil {
  lh.Error(ctx, err)
}
```

---

## 4. Example
```go
package main

import (
  "context"
  "github.com/LucaSchmitz2003/FlowWatch/loggingHelper"
  "github.com/LucaSchmitz2003/FlowWatch/otelHelper"
  "github.com/pkg/errors"
  "go.opentelemetry.io/otel"
)

var (
  CustomError1 = errors.New("Error message")
  tracer       = otel.Tracer("TestTracer")
  logger       = loggingHelper.GetLogHelper()
)

func errorTest() error {
  // Something went wrong
  err := errors.Wrap(CustomError1, "Error in errorTest()")
  return err
}

func main() {
  ctx := context.Background()

  // Initialize the OpenTelemetry SDK connection to the backend
  otelHelper.SetupOtelHelper()
  defer otelHelper.Shutdown() // Defer the shutdown function to ensure a graceful shutdown of the SDK connection at the end

  // Create a sub-span
  ctx, span := tracer.Start(ctx, "Test span")
  defer span.End()

  // Call function, catch error and log it
  err := errorTest()
  if err != nil {
	  logger.Warn(ctx, err)
  }
}
```

---

## 5. Import in other projects
```commandline
export GOPRIVATE=github.com/LucaSchmitz2003/*
GIT_SSH_COMMAND="ssh -v" go get github.com/LucaSchmitz2003/FlowWatch@main
```

## 6. Environment variables
```dotenv
OTEL_SERVICE_NAME="<name>"
OTEL_COLLECTOR_URL="<url>:<port>"
OTEL_SUPPORT_TLS=<bool>
```
