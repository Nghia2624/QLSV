package model

import "time"

type Registration struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	StudentID    string    `bson:"student_id" json:"student_id"`
	ClassID      string    `bson:"class_id" json:"class_id"`
	RegisteredAt time.Time `bson:"registered_at" json:"registered_at"`
	Status       string    `bson:"status" json:"status"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}
