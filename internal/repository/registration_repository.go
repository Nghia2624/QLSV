package repository

import "qlsvgo/internal/domain/model"

type RegistrationCommandRepository interface {
	Create(reg *model.Registration) error
	Update(reg *model.Registration) error
	Delete(id string) error
	GetByID(id string) (*model.Registration, error)
}

type RegistrationQueryRepository interface {
	GetByID(id string) (*model.Registration, error)
	GetAll() ([]*model.Registration, error)
	FindByStudentID(studentID string) ([]*model.Registration, error)
	FindByClassID(classID string) ([]*model.Registration, error)
	FindByStatus(status string) ([]*model.Registration, error)
}
