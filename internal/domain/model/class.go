package model

import "time"

type Class struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	ClassCode string    `bson:"class_code" json:"class_code"`
	CourseID  string    `bson:"course_id" json:"course_id"`
	TeacherID string    `bson:"teacher_id" json:"teacher_id"`
	Name      string    `bson:"name" json:"name"`
	Schedule  string    `bson:"schedule" json:"schedule"`
	Semester  string    `bson:"semester" json:"semester"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
