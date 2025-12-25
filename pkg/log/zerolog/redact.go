package zerolog

import "strings"

// RedactedValue marks sensitive data.
type RedactedValue struct {
	Original string
}

func Redacted(v string) RedactedValue {
	return RedactedValue{Original: v}
}

func redactValue(v any) any {
	switch t := v.(type) {

	case RedactedValue:
		return "[REDACTED]"

	case string:
		// Best-effort fallback protection
		if looksSensitive(t) {
			return "[REDACTED]"
		}
		return t

	default:
		return v
	}
}

func looksSensitive(s string) bool {
	l := strings.ToLower(s)

	return strings.Contains(l, "secret") ||
		strings.Contains(l, "token") ||
		strings.Contains(l, "password")
}
