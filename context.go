package trace

import "context"

import xtr "golang.org/x/net/trace"

type ctxKey int

var contextKey = ctxKey(0)

// TraceIDKey is the context key used to get the Trace's title out of a context if it exists
var TraceIDKey = ctxKey(1)

func contextWithTrace(ctx context.Context, trace *trace) context.Context {
	ctx = context.WithValue(ctx, contextKey, trace)
	ctx = context.WithValue(ctx, TraceIDKey, trace.id)
	return ctx
}

func traceFromContext(ctx context.Context) *trace {
	val := ctx.Value(contextKey)
	if trace, ok := val.(*trace); ok {
		return trace
	}
	return nil
}

func parentOrChildFromContext(ctx context.Context, family, title string) *trace {
	sp := traceFromContext(ctx)
	if sp == nil {
		return newTrace(family, title)
	}
	return sp.child(title)
}

// TitleFromContext is a convenience function that returns the Trace's title from a context
// or an empty string if none exists
func TitleFromContext(ctx context.Context) string {
	id, ok := ctx.Value(TraceIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

// NewContext returns a new context.Context and Trace with the given family and title.  The trace will
// be stored in the context.
func NewContext(ctx context.Context, family, title string) (xtr.Trace, context.Context) {
	sp := parentOrChildFromContext(ctx, family, title)
	return sp, contextWithTrace(ctx, sp)
}
