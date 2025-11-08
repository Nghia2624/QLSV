package repository

import "qlsvgo/internal/domain/model"

type CourseCommandRepository interface {
	Create(course *model.Course) error
	Update(course *model.Course) error
	Delete(id string) error
	GetByID(id string) (*model.Course, error)
}

type CourseQueryRepository interface {
	GetByID(id string) (*model.Course, error)
	GetAll() ([]*model.Course, error)
	FindByName(name string) ([]*model.Course, error)
}
