package postgres

import (
	"database/sql"
	"qlsvgo/internal/domain/model"
)

type TeacherPostgresRepository struct {
	DB *sql.DB
}

func (r *TeacherPostgresRepository) Create(teacher *model.Teacher) error {
	_, err := r.DB.Exec(
		`INSERT INTO teachers (id, name, email, department, address, phone, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		teacher.ID, teacher.Name, teacher.Email, teacher.Department, teacher.Address, teacher.Phone, teacher.CreatedAt, teacher.UpdatedAt,
	)
	return err
}

func (r *TeacherPostgresRepository) Update(teacher *model.Teacher) error {
	_, err := r.DB.Exec(
		`UPDATE teachers SET teacher_code=$1, name=$2, email=$3, department=$4, address=$5, phone=$6, updated_at=$7 WHERE id=$8`,
		teacher.TeacherCode, teacher.Name, teacher.Email, teacher.Department, teacher.Address, teacher.Phone, teacher.UpdatedAt, teacher.ID,
	)
	return err
}

func (r *TeacherPostgresRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM teachers WHERE id=$1`, id)
	return err
}

func (r *TeacherPostgresRepository) GetByID(id string) (*model.Teacher, error) {
	row := r.DB.QueryRow(`SELECT id, teacher_code, name, email, department, address, phone, created_at, updated_at FROM teachers WHERE id=$1`, id)
	var t model.Teacher
	err := row.Scan(&t.ID, &t.TeacherCode, &t.Name, &t.Email, &t.Department, &t.Address, &t.Phone, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
