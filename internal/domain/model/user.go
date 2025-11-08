package model

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
	Role     string `bson:"role" json:"role"`     // admin, student, teacher
	RefID    string `bson:"ref_id" json:"ref_id"` // Liên kết với id của student/teacher nếu có
}
