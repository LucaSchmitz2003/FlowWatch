package loggingHelper

import "context"

// Abstraction for log functions to enable simpler switching between logging libraries.
// Context is required to add the event to the span (if possible).

// Debug logs a message at the debug level.
func (lh *LogHelper) Debug(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Debug(args...)
}

// Info logs a message at the info level.
func (lh *LogHelper) Info(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Info(args...)
}

// Warn logs a message at the warning level.
func (lh *LogHelper) Warn(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Warn(args...)
}

// Error logs a message at the error level.
func (lh *LogHelper) Error(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Error(args...)
}

// Fatal logs a message at the fatal level.
func (lh *LogHelper) Fatal(ctx context.Context, args ...interface{}) {
	lh.Logger.WithContext(ctx).Fatal(args...)
}
