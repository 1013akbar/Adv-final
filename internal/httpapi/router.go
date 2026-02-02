package httpapi

import (
	"net/http"
	"strings"

	"course-registration/internal/service"
	"course-registration/internal/store"
)

func NewRouter(st *store.Store, enrollSvc *service.EnrollmentService, audit *service.AuditWorker) http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Students
	mux.HandleFunc("/students", StudentsHandler(st))

	// Courses (CRUD)
	mux.HandleFunc("/courses", CoursesHandler(st))
	mux.HandleFunc("/courses/", func(w http.ResponseWriter, r *http.Request) {
		// /courses/{id}
		id := strings.TrimPrefix(r.URL.Path, "/courses/")
		CourseByIDHandler(st, id)(w, r)
	})

	// Enrollments
	mux.HandleFunc("/enrollments", EnrollmentsHandler(enrollSvc, st))

	// Audit (extra endpoint helps demo goroutine)
	mux.HandleFunc("/audit", AuditHandler(audit))

	return mux
}
