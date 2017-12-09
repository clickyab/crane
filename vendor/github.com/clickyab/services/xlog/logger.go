package xlog

import (
	"context"

	"sync"

	"github.com/sirupsen/logrus"
)

type contextKey int

const ctxKey contextKey = iota

type concurrentFields struct {
	fields logrus.Fields
	lock   sync.RWMutex
}

// Get return the logger and initialize it based on context ctxKey
func Get(ctx context.Context) *logrus.Entry {
	fields, ok := ctx.Value(ctxKey).(*concurrentFields)
	entry := logrus.NewEntry(logrus.StandardLogger())
	if ok {
		fields.lock.RLock()
		defer fields.lock.RUnlock()
		return entry.WithFields(fields.fields)
	}

	return entry
}

// GetWithError is a shorthand for Get(ctx).WithError(err)
func GetWithError(ctx context.Context, err error) *logrus.Entry {
	return Get(ctx).WithError(err)
}

// GetWithField is a shorthand for Get().WithField()
func GetWithField(ctx context.Context, key string, val interface{}) *logrus.Entry {
	return Get(ctx).WithField(key, val)
}

// GetWithFields is a shorthand for Get().WithFields()
func GetWithFields(ctx context.Context, f logrus.Fields) *logrus.Entry {
	return Get(ctx).WithFields(f)
}

// SetField in the context
func SetField(ctx context.Context, key string, val interface{}) context.Context {
	fields, ok := ctx.Value(ctxKey).(*concurrentFields)
	if !ok {
		fields = &concurrentFields{
			fields: make(logrus.Fields),
			lock:   sync.RWMutex{},
		}
	}
	fields.lock.Lock()
	defer fields.lock.Unlock()
	fields.fields[key] = val

	return context.WithValue(ctx, ctxKey, fields)
}

// SetFields set the fields for logger
func SetFields(ctx context.Context, fl logrus.Fields) context.Context {
	fields, ok := ctx.Value(ctxKey).(*concurrentFields)
	if !ok {
		fields = &concurrentFields{
			fields: make(logrus.Fields),
			lock:   sync.RWMutex{},
		}
	}
	fields.lock.Lock()
	defer fields.lock.Unlock()
	for i := range fl {
		fields.fields[i] = fl[i]
	}
	return context.WithValue(ctx, ctxKey, fields)
}
