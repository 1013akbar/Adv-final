package httpapi

import (
	"net/http"
	"strings"

	"course-registration/internal/store"
)

type createCourseReq struct {
	Code         string `json:"code"`
	Title        string `json:"title"`
	Capacity     int    `json:"capacity"`
	InstructorID string `json:"instructor_id,omitempty"`
}

func CoursesHandler(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			var req createCourseReq
			if err := readJSON(r, &req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			req.Code = strings.TrimSpace(req.Code)
			req.Title = strings.TrimSpace(req.Title)
			if req.Code == "" || req.Title == "" || req.Capacity <= 0 {
				http.Error(w, "code, title required; capacity must be > 0", http.StatusBadRequest)
				return
			}

			c := st.CreateCourse(req.Code, req.Title, req.Capacity, req.InstructorID)
			writeJSON(w, http.StatusCreated, c)

		case http.MethodGet:
			writeJSON(w, http.StatusOK, st.ListCourses())

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func CourseByIDHandler(st *store.Store, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id = strings.TrimSpace(id)
		if id == "" {
			http.Error(w, "missing course id", http.StatusBadRequest)
			return
		}

		switch r.Method {

		case http.MethodGet:
			c, ok := st.GetCourse(id)
			if !ok {
				http.Error(w, "course not found", http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, c)

		case http.MethodPut:
			var req createCourseReq
			if err := readJSON(r, &req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			req.Code = strings.TrimSpace(req.Code)
			req.Title = strings.TrimSpace(req.Title)
			if req.Code == "" || req.Title == "" || req.Capacity <= 0 {
				http.Error(w, "code, title required; capacity must be > 0", http.StatusBadRequest)
				return
			}
			updated, err := st.UpdateCourse(id, req.Code, req.Title, req.Capacity, req.InstructorID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, updated)

		case http.MethodDelete:
			if err := st.DeleteCourse(id); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
