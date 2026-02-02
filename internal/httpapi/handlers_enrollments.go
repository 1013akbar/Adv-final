package httpapi

import (
	"net/http"
	"strings"

	"course-registration/internal/service"
	"course-registration/internal/store"
)

type enrollReq struct {
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
}

func EnrollmentsHandler(enrollSvc *service.EnrollmentService, st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			var req enrollReq
			if err := readJSON(r, &req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			req.StudentID = strings.TrimSpace(req.StudentID)
			req.CourseID = strings.TrimSpace(req.CourseID)
			if req.StudentID == "" || req.CourseID == "" {
				http.Error(w, "student_id and course_id are required", http.StatusBadRequest)
				return
			}

			e, err := enrollSvc.Enroll(req.StudentID, req.CourseID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, http.StatusCreated, e)

		case http.MethodGet:
			writeJSON(w, http.StatusOK, st.ListEnrollments())

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
