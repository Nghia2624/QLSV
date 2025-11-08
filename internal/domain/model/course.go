package model

import "time"

type Course struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	CourseCode  string    `bson:"course_code" json:"course_code"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Credits     int       `bson:"credits" json:"credits"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
