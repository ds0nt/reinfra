package logmw

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxMarker struct{}

var (
	// ctxMarkerKey is the Context value marker used by *all* logging middleware.
	// The logging middleware object must interf
	ctxMarkerKey = &ctxMarker{}
)

// Tags is the struct used for storing request tags between Context calls.
// This object is *not* thread safe, and should be handled only in the context of the request.
type Tags struct {
	values map[string]interface{}
}

// Set sets the given key in the metadata tags.
func (t *Tags) Set(key string, value interface{}) *Tags {
	t.values[key] = value
	return t
}

// Has checks if the given key exists.
func (t *Tags) Has(key string) bool {
	_, ok := t.values[key]
	return ok
}

// Values returns a map of key to values.
// Do not modify the underlying map, please use Set instead.
func (t *Tags) Values() map[string]interface{} {
	return t.values
}

// Extracts returns a pre-existing Tags object in the Context.
// If the context wasn't set in a tag interceptor, a no-op Tag storage is returned that will *not* be propagated in context.
func Extract(ctx context.Context) *Tags {
	t, ok := ctx.Value(ctxMarkerKey).(*Tags)
	if !ok {
		return &Tags{values: make(map[string]interface{})}
	}
	return t
}

func setInContext(ctx context.Context, tags *Tags) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, tags)
}

type logrusCtxMarker struct{}

var (
	logrusCtxMarkerKey = &logrusCtxMarker{}
)

// Extract takes the call-scoped logrus.Entry from grpc_logrus middleware.
//
// If the grpc_logrus middleware wasn't used, a no-op `logrus.Entry` is returned. This makes it safe to
// use regardless.
func ExtractEntry(ctx context.Context) *logrus.Entry {
	l, ok := ctx.Value(logrusCtxMarkerKey).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(nullLogger)
	}
	// Add grpc_ctxtags tags metadata until now.
	return l.WithFields(logrus.Fields(Extract(ctx).Values()))
}

func EntryToContext(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, logrusCtxMarkerKey, entry)

}
