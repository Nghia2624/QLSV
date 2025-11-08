package postgres

import (
	"database/sql"

	"qlsvgo/internal/domain/model"
)

type ClassPostgresRepository struct {
	DB *sql.DB
}

func (r *ClassPostgresRepository) Create(class *model.Class) error {
	_, err := r.DB.Exec(
		`INSERT INTO classes (id, course_id, teacher_id, name, schedule, semester, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		class.ID, class.CourseID, class.TeacherID, class.Name, class.Schedule, class.Semester, class.CreatedAt, class.UpdatedAt,
	)
	return err
}

func (r *ClassPostgresRepository) Update(class *model.Class) error {
	_, err := r.DB.Exec(
		`UPDATE classes SET class_code=$1, course_id=$2, teacher_id=$3, name=$4, schedule=$5, semester=$6, updated_at=$7 WHERE id=$8`,
		class.ClassCode, class.CourseID, class.TeacherID, class.Name, class.Schedule, class.Semester, class.UpdatedAt, class.ID,
	)
	return err
}

func (r *ClassPostgresRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM classes WHERE id=$1`, id)
	return err
}

func (r *ClassPostgresRepository) GetByID(id string) (*model.Class, error) {
	row := r.DB.QueryRow(`SELECT id, class_code, course_id, teacher_id, name, schedule, semester, created_at, updated_at FROM classes WHERE id=$1`, id)
	var cl model.Class
	err := row.Scan(&cl.ID, &cl.ClassCode, &cl.CourseID, &cl.TeacherID, &cl.Name, &cl.Schedule, &cl.Semester, &cl.CreatedAt, &cl.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}
