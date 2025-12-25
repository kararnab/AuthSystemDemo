package stdout

import (
	"context"

	"github.com/kararnab/authdemo/pkg/iam/audit"
	"github.com/kararnab/authdemo/pkg/log"
)

//
// ================================
// AUDIT LOGGER
// ================================
//
// NOTE:
//   - Logs provide context
//   - Rates/counts must come from metrics
//

type AuditLogger struct{}

func (l *AuditLogger) Log(
	ctx context.Context,
	event audit.Event,
) error {
	//log.Printf(
	//	"[AUDIT] type=%s subject=%s provider=%s msg=%s attrs=%v",
	//	event.Type,
	//	event.SubjectID,
	//	event.Provider,
	//	event.Message,
	//	event.Attrs,
	//)
	log.Info(
		event.Message,
		log.F("log_type", "audit", log.RedactNone),
		log.F("event_type", string(event.Type), log.RedactNone),
		log.F("subject_id", event.SubjectID, log.RedactFull),
		log.F("provider", event.Provider, log.RedactNone),
		log.F("attrs", "[omitted]", log.RedactFull),
	)
	return nil
}
