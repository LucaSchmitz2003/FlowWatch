# FlowWatch
An abstraction layer of Open Telemetry to standardize Tracing, Logging (also local), Metrics and Exception Handling for every program of the companies eco system. 

## Usage
### OpenTelemetry connection
- `otelHelper.SetupOtelHelper()`: To be called at the beginning of your program to set up the OpenTelemetry connection.
- `otelHelper.Shutdown()`: Needs to be called at the end of your program to ensure all data is sent to the collector, 
before the program exits.\
→ Defer call is recommended.
- `otel.Tracer(<tracer_name>)`: To start a trace, the tracer needs to be initialized with a name.\
***Recommendation:** Set the tracer as a global variable in your program.*
- `tracer.Start(<context>, <span_name>)`: Start a span with the context and a name. It returns the updated context and the span.
- `<span_name>.End()`: To end the span.\
  → Defer call is recommended.
- *Every log message will be attached to the span, if the updated context is used.*
- *If a span is created inside an existing span (using its context) the new span will be a child of the existing span.*

### Logging
- `loggingHelper.GetLogHelper()`: To get the logging helper.
- `<loggingHelper_name>.<log_level>(<context>, [<string-1>, ...])`: Log a message with the given log level.\
  → Each log level is represented by a logging function:
  - Debug
  - Info
  - Warn
  - Error
  - Fatal

### Exceptions
- Recommended to use the `pkg/errors` package to create custom errors, since it allows wrapping.
- The LoggingHelper can handle these errors.
- Errors should be declared as global variables in the program, to ensure that the error message is consistent throughout the program.

## Example

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

  // Initialize the logging helper
  lh := loggingHelper.GetLogHelper()

  // Create a sub-span
  ctx, span := tracer.Start(ctx, "Test span")
  defer span.End()
  
  // Call function, catch error and log it
  err := errorTest()
  if err != nil {
    lh.Warn(ctx, err)
  }
}
```