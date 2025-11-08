package postgres

import (
	"database/sql"

	"qlsvgo/internal/domain/model"
)

type RegistrationPostgresRepository struct {
	DB *sql.DB
}

func (r *RegistrationPostgresRepository) Create(reg *model.Registration) error {
	_, err := r.DB.Exec(
		`INSERT INTO registrations (id, student_id, class_id, registered_at, status, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		reg.ID, reg.StudentID, reg.ClassID, reg.RegisteredAt, reg.Status, reg.CreatedAt, reg.UpdatedAt,
	)
	return err
}

func (r *RegistrationPostgresRepository) Update(reg *model.Registration) error {
	_, err := r.DB.Exec(
		`UPDATE registrations SET student_id=$1, class_id=$2, registered_at=$3, status=$4, updated_at=$5 WHERE id=$6`,
		reg.StudentID, reg.ClassID, reg.RegisteredAt, reg.Status, reg.UpdatedAt, reg.ID,
	)
	return err
}

func (r *RegistrationPostgresRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM registrations WHERE id=$1`, id)
	return err
}

func (r *RegistrationPostgresRepository) GetByID(id string) (*model.Registration, error) {
	row := r.DB.QueryRow(`SELECT id, student_id, class_id, registered_at, status, created_at, updated_at FROM registrations WHERE id=$1`, id)
	var reg model.Registration
	err := row.Scan(&reg.ID, &reg.StudentID, &reg.ClassID, &reg.RegisteredAt, &reg.Status, &reg.CreatedAt, &reg.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}
