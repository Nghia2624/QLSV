package repository

import "qlsvgo/internal/domain/model"

type StudentCommandRepository interface {
	Create(student *model.Student) error
	Update(student *model.Student) error
	Delete(id string) error
	GetByID(id string) (*model.Student, error)
}

type StudentQueryRepository interface {
	GetByID(id string) (*model.Student, error)
	GetAll() ([]*model.Student, error)
	FindByEmail(email string) (*model.Student, error)
}
