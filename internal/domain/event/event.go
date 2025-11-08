package event

import "time"

// EventType định nghĩa các loại event
const (
	UserCreated         = "UserCreated"
	UserUpdated         = "UserUpdated"
	UserDeleted         = "UserDeleted"
	StudentCreated      = "StudentCreated"
	StudentUpdated      = "StudentUpdated"
	StudentDeleted      = "StudentDeleted"
	TeacherCreated      = "TeacherCreated"
	TeacherUpdated      = "TeacherUpdated"
	TeacherDeleted      = "TeacherDeleted"
	CourseCreated       = "CourseCreated"
	CourseUpdated       = "CourseUpdated"
	CourseDeleted       = "CourseDeleted"
	ClassCreated        = "ClassCreated"
	ClassUpdated        = "ClassUpdated"
	ClassDeleted        = "ClassDeleted"
	RegistrationCreated = "RegistrationCreated"
	RegistrationUpdated = "RegistrationUpdated"
	RegistrationDeleted = "RegistrationDeleted"
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	CreatedAt time.Time   `json:"created_at"`
}
