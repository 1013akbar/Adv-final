package domain

type Student struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type Instructor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Course struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Capacity    int    `json:"capacity"`
	InstructorID string `json:"instructor_id,omitempty"`
}

type EnrollmentStatus string

const (
	StatusEnrolled  EnrollmentStatus = "ENROLLED"
	StatusWaitlisted EnrollmentStatus = "WAITLISTED"
	StatusDropped   EnrollmentStatus = "DROPPED"
)

type Enrollment struct {
	ID        string           `json:"id"`
	StudentID string           `json:"student_id"`
	CourseID  string           `json:"course_id"`
	Status    EnrollmentStatus `json:"status"`
}
