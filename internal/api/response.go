package api

import (
	"encoding/json"
	//"log"
	"github.com/kararnab/authdemo/pkg/log"

	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		// At this point headers are already written
		// Best effort logging only
		//log.Printf("failed to write JSON response: %v", err)
		log.Unsafe().Error(
			"failed to write JSON response",
			log.F("error", err, log.RedactNone),
		)
	}
}
