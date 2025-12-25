package log

import "context"

type noopLogger struct{}

func (noopLogger) Debug(string, ...Field) {}
func (noopLogger) Info(string, ...Field)  {}
func (noopLogger) Warn(string, ...Field)  {}
func (noopLogger) Error(string, ...Field) {}

func (n noopLogger) With(...Field) Logger               { return n }
func (n noopLogger) WithContext(context.Context) Logger { return n }
func (n noopLogger) Unsafe() UnsafeLogger               { return noopUnsafeLogger{} }

type noopUnsafeLogger struct{}

func (noopUnsafeLogger) Debug(string, ...any) {}
func (noopUnsafeLogger) Info(string, ...any)  {}
func (noopUnsafeLogger) Warn(string, ...any)  {}
func (noopUnsafeLogger) Error(string, ...any) {}
