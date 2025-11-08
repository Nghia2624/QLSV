package repository

import "qlsvgo/internal/domain/model"

type ClassCommandRepository interface {
	Create(class *model.Class) error
	Update(class *model.Class) error
	Delete(id string) error
	GetByID(id string) (*model.Class, error)
}

type ClassQueryRepository interface {
	GetByID(id string) (*model.Class, error)
	GetAll() ([]*model.Class, error)
	FindByCourseID(courseID string) ([]*model.Class, error)
	FindByTeacherID(teacherID string) ([]*model.Class, error)
	FindBySemester(semester string) ([]*model.Class, error)
}
