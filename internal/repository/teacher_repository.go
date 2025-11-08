package repository

import "qlsvgo/internal/domain/model"

type TeacherCommandRepository interface {
	Create(teacher *model.Teacher) error
	Update(teacher *model.Teacher) error
	Delete(id string) error
	GetByID(id string) (*model.Teacher, error)
}

type TeacherQueryRepository interface {
	GetByID(id string) (*model.Teacher, error)
	GetAll() ([]*model.Teacher, error)
	FindByEmail(email string) (*model.Teacher, error)
}
