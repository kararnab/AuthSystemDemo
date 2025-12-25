package log

func SanitizeFields(fields []Field) map[string]any {
	out := make(map[string]any, len(fields))
	for _, f := range fields {
		out[f.Key] = applyRedaction(f.Redaction, f.Value)
	}
	return out
}
