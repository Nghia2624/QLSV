package postgres

import (
	"database/sql"
	"qlsvgo/internal/domain/model"
)

type CoursePostgresRepository struct {
	DB *sql.DB
}

func (r *CoursePostgresRepository) Create(course *model.Course) error {
	_, err := r.DB.Exec(
		`INSERT INTO courses (id, name, description, credits, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		course.ID, course.Name, course.Description, course.Credits, course.CreatedAt, course.UpdatedAt,
	)
	return err
}

func (r *CoursePostgresRepository) Update(course *model.Course) error {
	_, err := r.DB.Exec(
		`UPDATE courses SET course_code=$1, name=$2, description=$3, credits=$4, updated_at=$5 WHERE id=$6`,
		course.CourseCode, course.Name, course.Description, course.Credits, course.UpdatedAt, course.ID,
	)
	return err
}

func (r *CoursePostgresRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM courses WHERE id=$1`, id)
	return err
}

func (r *CoursePostgresRepository) GetByID(id string) (*model.Course, error) {
	row := r.DB.QueryRow(`SELECT id, course_code, name, description, credits, created_at, updated_at FROM courses WHERE id=$1`, id)
	var c model.Course
	err := row.Scan(&c.ID, &c.CourseCode, &c.Name, &c.Description, &c.Credits, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
