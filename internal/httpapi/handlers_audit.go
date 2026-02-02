package httpapi

import (
	"net/http"

	"course-registration/internal/service"
)

func AuditHandler(audit *service.AuditWorker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, audit.List())
	}
}
