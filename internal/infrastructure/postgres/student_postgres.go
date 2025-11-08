package postgres

import (
	"database/sql"
	"qlsvgo/internal/domain/model"
)

type StudentPostgresRepository struct {
	DB *sql.DB
}

func (r *StudentPostgresRepository) Create(student *model.Student) error {
	_, err := r.DB.Exec(
		`INSERT INTO students (id, name, email, dob, gender, address, phone, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		student.ID, student.Name, student.Email, student.DOB, student.Gender, student.Address, student.Phone, student.CreatedAt, student.UpdatedAt,
	)
	return err
}

func (r *StudentPostgresRepository) Update(student *model.Student) error {
	_, err := r.DB.Exec(
		`UPDATE students SET student_code=$1, name=$2, email=$3, dob=$4, gender=$5, address=$6, phone=$7, updated_at=$8 WHERE id=$9`,
		student.StudentCode, student.Name, student.Email, student.DOB, student.Gender, student.Address, student.Phone, student.UpdatedAt, student.ID,
	)
	return err
}

func (r *StudentPostgresRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM students WHERE id=$1`, id)
	return err
}

func (r *StudentPostgresRepository) GetByID(id string) (*model.Student, error) {
	row := r.DB.QueryRow(`SELECT id, student_code, name, email, dob, gender, address, phone, created_at, updated_at FROM students WHERE id=$1`, id)
	var s model.Student
	err := row.Scan(&s.ID, &s.StudentCode, &s.Name, &s.Email, &s.DOB, &s.Gender, &s.Address, &s.Phone, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
