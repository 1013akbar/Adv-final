package store

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"

	"course-registration/internal/domain"
)

type Store struct {
	mu sync.RWMutex

	students    map[string]domain.Student
	instructors map[string]domain.Instructor
	courses     map[string]domain.Course
	enrollments map[string]domain.Enrollment

	// Helps prevent duplicates quickly: key = studentID + ":" + courseID
	enrollmentIndex map[string]string // -> enrollmentID

	idCounter uint64
}

func NewStore() *Store {
	return &Store{
		students:         make(map[string]domain.Student),
		instructors:      make(map[string]domain.Instructor),
		courses:          make(map[string]domain.Course),
		enrollments:      make(map[string]domain.Enrollment),
		enrollmentIndex:  make(map[string]string),
	}
}

func (s *Store) nextID(prefix string) string {
	n := atomic.AddUint64(&s.idCounter, 1)
	return prefix + "-" + strconv.FormatUint(n, 10)
}

/* -------------------- Students -------------------- */

func (s *Store) CreateStudent(fullName, email string) domain.Student {
	s.mu.Lock()
	defer s.mu.Unlock()

	st := domain.Student{
		ID:       s.nextID("stu"),
		FullName: fullName,
		Email:    email,
	}
	s.students[st.ID] = st
	return st
}

func (s *Store) ListStudents() []domain.Student {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]domain.Student, 0, len(s.students))
	for _, v := range s.students {
		out = append(out, v)
	}
	return out
}

func (s *Store) GetStudent(id string) (domain.Student, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.students[id]
	return st, ok
}

/* -------------------- Courses (CRUD) -------------------- */

func (s *Store) CreateCourse(code, title string, capacity int, instructorID string) domain.Course {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := domain.Course{
		ID:           s.nextID("crs"),
		Code:         code,
		Title:        title,
		Capacity:     capacity,
		InstructorID: instructorID,
	}
	s.courses[c.ID] = c
	return c
}

func (s *Store) ListCourses() []domain.Course {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]domain.Course, 0, len(s.courses))
	for _, v := range s.courses {
		out = append(out, v)
	}
	return out
}

func (s *Store) GetCourse(id string) (domain.Course, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.courses[id]
	return c, ok
}

func (s *Store) UpdateCourse(id, code, title string, capacity int, instructorID string) (domain.Course, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.courses[id]
	if !ok {
		return domain.Course{}, errors.New("course not found")
	}

	// simple update
	c.Code = code
	c.Title = title
	c.Capacity = capacity
	c.InstructorID = instructorID
	s.courses[id] = c
	return c, nil
}

func (s *Store) DeleteCourse(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.courses[id]; !ok {
		return errors.New("course not found")
	}
	delete(s.courses, id)
	return nil
}

/* -------------------- Enrollments -------------------- */

func pairKey(studentID, courseID string) string {
	return studentID + ":" + courseID
}

func (s *Store) CountEnrolled(courseID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, e := range s.enrollments {
		if e.CourseID == courseID && e.Status == domain.StatusEnrolled {
			count++
		}
	}
	return count
}

func (s *Store) HasEnrollment(studentID, courseID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.enrollmentIndex[pairKey(studentID, courseID)]
	return ok
}

func (s *Store) CreateEnrollment(studentID, courseID string, status domain.EnrollmentStatus) domain.Enrollment {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := domain.Enrollment{
		ID:        s.nextID("enr"),
		StudentID: studentID,
		CourseID:  courseID,
		Status:    status,
	}
	s.enrollments[e.ID] = e
	s.enrollmentIndex[pairKey(studentID, courseID)] = e.ID
	return e
}

func (s *Store) ListEnrollments() []domain.Enrollment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]domain.Enrollment, 0, len(s.enrollments))
	for _, v := range s.enrollments {
		out = append(out, v)
	}
	return out
}
