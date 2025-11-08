package usecase

import (
	"qlsvgo/internal/domain/event"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/mapping"
	"qlsvgo/internal/repository"

	"github.com/google/uuid"
)

type UserUsecase struct {
	Repo     repository.UserRepository
	EventBus func(*event.Event) error
}

func (uc *UserUsecase) Register(user *model.User) error {
	// Tạo UUID cho user nếu chưa có
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Ghi vào database
	if err := uc.Repo.Create(user); err != nil {
		return err
	}

	// Publish event nếu có EventBus
	if uc.EventBus != nil {
		evt := mapping.UserToEvent(user, event.UserCreated)
		if err := uc.EventBus(evt); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UserUsecase) GetByUsername(username string) (*model.User, error) {
	return uc.Repo.GetByUsername(username)
}

func (uc *UserUsecase) GetByID(id string) (*model.User, error) {
	return uc.Repo.GetByID(id)
}
