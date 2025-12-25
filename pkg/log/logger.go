package log

import "context"

// ---------- public API ----------

var L Logger = noopLogger{}

func Init(l Logger) {
	if l == nil {
		panic("log.Init: nil logger")
	}
	L = l
}

func Debug(msg string, fields ...Field) { L.Debug(msg, fields...) }
func Info(msg string, fields ...Field)  { L.Info(msg, fields...) }
func Warn(msg string, fields ...Field)  { L.Warn(msg, fields...) }
func Error(msg string, fields ...Field) { L.Error(msg, fields...) }

func With(fields ...Field) Logger {
	return L.With(fields...)
}

func WithContext(ctx context.Context) Logger {
	return L.WithContext(ctx)
}

// Deprecated: Unsafe bypasses redaction and must not be used in production.
func Unsafe() UnsafeLogger {
	return L.Unsafe()
}

// ---------- types ----------

type Field struct {
	Key       string
	Value     any
	Redaction RedactionType
}

func F(key string, value any, r RedactionType) Field {
	return Field{Key: key, Value: value, Redaction: r}
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)

	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
	Unsafe() UnsafeLogger
}

type UnsafeLogger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)
}
