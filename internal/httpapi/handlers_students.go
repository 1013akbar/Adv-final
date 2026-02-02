package httpapi

import (
	"net/http"
	"strings"

	"course-registration/internal/store"
)

type createStudentReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func StudentsHandler(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			var req createStudentReq
			if err := readJSON(r, &req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			req.FullName = strings.TrimSpace(req.FullName)
			req.Email = strings.TrimSpace(req.Email)
			if req.FullName == "" || req.Email == "" {
				http.Error(w, "full_name and email are required", http.StatusBadRequest)
				return
			}

			s := st.CreateStudent(req.FullName, req.Email)
			writeJSON(w, http.StatusCreated, s)

		case http.MethodGet:
			writeJSON(w, http.StatusOK, st.ListStudents())

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
