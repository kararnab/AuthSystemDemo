package log

import "sync"

type RedactionType int

const (
	RedactNone RedactionType = iota
	RedactFull
)

type Redactor func(any) any

var (
	mu        sync.RWMutex
	redactors = map[RedactionType]Redactor{}
)

func RegisterRedactor(rt RedactionType, fn Redactor) {
	mu.Lock()
	defer mu.Unlock()
	redactors[rt] = fn
}

func applyRedaction(rt RedactionType, v any) any {
	if rt == RedactNone {
		return v
	}
	mu.RLock()
	fn := redactors[rt]
	mu.RUnlock()
	if fn == nil {
		return "[REDACTION_NOT_REGISTERED]"
	}
	return fn(v)
}

func init() {
	RegisterRedactor(RedactFull, func(any) any {
		return "[REDACTED]"
	})
}
