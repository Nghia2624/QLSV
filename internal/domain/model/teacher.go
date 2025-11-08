package model

import "time"

type Teacher struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	TeacherCode string    `bson:"teacher_code" json:"teacher_code"`
	Name        string    `bson:"name" json:"name"`
	Email       string    `bson:"email" json:"email"`
	Department  string    `bson:"department" json:"department"`
	Address     string    `bson:"address" json:"address"`
	Phone       string    `bson:"phone" json:"phone"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
