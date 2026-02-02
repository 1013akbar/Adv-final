package service

import (
	"errors"
	"sync"
	"time"

	"course-registration/internal/domain"
	"course-registration/internal/store"
)

/* -------------------- Audit Worker (goroutine) -------------------- */

type AuditEvent struct {
	Time      time.Time `json:"time"`
	Action    string    `json:"action"`
	StudentID string    `json:"student_id"`
	CourseID  string    `json:"course_id"`
	Result    string    `json:"result"`
}

type AuditWorker struct {
	ch     chan AuditEvent
	stopCh chan struct{}

	mu     sync.Mutex
	events []AuditEvent
}

func NewAuditWorker(buffer int) *AuditWorker {
	return &AuditWorker{
		ch:     make(chan AuditEvent, buffer),
		stopCh: make(chan struct{}),
		events: make([]AuditEvent, 0),
	}
}

func (a *AuditWorker) Start() {
	go func() {
		for {
			select {
			case ev := <-a.ch:
				a.mu.Lock()
				a.events = append(a.events, ev)
				a.mu.Unlock()
			case <-a.stopCh:
				return
			}
		}
	}()
}

func (a *AuditWorker) Stop() {
	close(a.stopCh)
}

func (a *AuditWorker) Publish(ev AuditEvent) {
	// Non-blocking best effort
	select {
	case a.ch <- ev:
	default:
	}
}

func (a *AuditWorker) List() []AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()

	out := make([]AuditEvent, 0, len(a.events))
	out = append(out, a.events...)
	return out
}

/* -------------------- Enrollment Service -------------------- */

type EnrollmentService struct {
	st    *store.Store
	audit *AuditWorker
}

func NewEnrollmentService(st *store.Store, audit *AuditWorker) *EnrollmentService {
	return &EnrollmentService{st: st, audit: audit}
}

// Enroll enforces capacity and prevents duplicates (core problem solved)
func (s *EnrollmentService) Enroll(studentID, courseID string) (domain.Enrollment, error) {
	// Validate student exists
	if _, ok := s.st.GetStudent(studentID); !ok {
		return domain.Enrollment{}, errors.New("student not found")
	}

	// Validate course exists
	c, ok := s.st.GetCourse(courseID)
	if !ok {
		return domain.Enrollment{}, errors.New("course not found")
	}

	// Prevent duplicate enrollment
	if s.st.HasEnrollment(studentID, courseID) {
		s.audit.Publish(AuditEvent{
			Time: time.Now(), Action: "ENROLL",
			StudentID: studentID, CourseID: courseID, Result: "DUPLICATE_BLOCKED",
		})
		return domain.Enrollment{}, errors.New("student already enrolled or waitlisted")
	}

	// Capacity enforcement
	enrolled := s.st.CountEnrolled(courseID)
	status := domain.StatusEnrolled
	if enrolled >= c.Capacity {
		status = domain.StatusWaitlisted
	}

	e := s.st.CreateEnrollment(studentID, courseID, status)

	// Async audit log (goroutine requirement)
	s.audit.Publish(AuditEvent{
		Time: time.Now(), Action: "ENROLL",
		StudentID: studentID, CourseID: courseID, Result: string(status),
	})

	return e, nil
}
