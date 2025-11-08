package usecase

import (
	"qlsvgo/internal/domain/event"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/mapping"
	"qlsvgo/internal/repository"
	"time"

	"github.com/google/uuid"
)

type RegistrationUsecase struct {
	CommandRepo repository.RegistrationCommandRepository
	QueryRepo   repository.RegistrationQueryRepository
	EventBus    func(evt *event.Event) error // publish event
}

func (uc *RegistrationUsecase) Create(reg *model.Registration) error {
	if reg.ID == "" {
		reg.ID = uuid.New().String()
	}
	err := uc.CommandRepo.Create(reg)
	if err != nil {
		return err
	}
	evt := mapping.RegistrationToEvent(reg, event.RegistrationCreated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *RegistrationUsecase) Update(reg *model.Registration) error {
	err := uc.CommandRepo.Update(reg)
	if err != nil {
		return err
	}
	evt := mapping.RegistrationToEvent(reg, event.RegistrationUpdated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *RegistrationUsecase) Delete(id string) error {
	err := uc.CommandRepo.Delete(id)
	if err != nil {
		return err
	}
	evt := &event.Event{
		Type:      event.RegistrationDeleted,
		Payload:   map[string]interface{}{"id": id},
		CreatedAt: time.Now(),
	}
	return uc.EventBus(evt)
}

func (uc *RegistrationUsecase) GetByID(id string) (*model.Registration, error) {
	return uc.QueryRepo.GetByID(id)
}

func (uc *RegistrationUsecase) GetAll() ([]*model.Registration, error) {
	return uc.QueryRepo.GetAll()
}

// Các hàm Update, Delete, GetByID, GetAll tương tự, publish event khi cần
