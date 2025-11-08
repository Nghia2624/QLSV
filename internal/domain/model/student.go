package model

import "time"

type Student struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	StudentCode string    `bson:"student_code" json:"student_code"`
	Name        string    `bson:"name" json:"name"`
	Email       string    `bson:"email" json:"email"`
	DOB         string    `bson:"dob" json:"dob"`
	Gender      string    `bson:"gender" json:"gender"`
	Address     string    `bson:"address" json:"address"`
	Phone       string    `bson:"phone" json:"phone"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
