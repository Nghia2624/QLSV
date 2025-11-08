package repository

import "qlsvgo/internal/domain/model"

type UserRepository interface {
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
	GetByID(id string) (*model.User, error)
	GetAll() ([]*model.User, error)
	Update(user *model.User) error
	Delete(id string) error
}
