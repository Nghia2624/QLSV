package usecase

import (
	"qlsvgo/internal/domain/event"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/mapping"
	"qlsvgo/internal/repository"
	"time"
	"github.com/google/uuid"
)

type ClassUsecase struct {
	CommandRepo repository.ClassCommandRepository
	QueryRepo   repository.ClassQueryRepository
	EventBus    func(evt *event.Event) error // publish event
}

func (uc *ClassUsecase) Create(class *model.Class) error {
	if class.ID == "" {
		class.ID = uuid.New().String()
	}
	err := uc.CommandRepo.Create(class)
	if err != nil {
		return err
	}
	evt := mapping.ClassToEvent(class, event.ClassCreated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *ClassUsecase) Update(class *model.Class) error {
	err := uc.CommandRepo.Update(class)
	if err != nil {
		return err
	}
	evt := mapping.ClassToEvent(class, event.ClassUpdated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *ClassUsecase) Delete(id string) error {
	err := uc.CommandRepo.Delete(id)
	if err != nil {
		return err
	}
	evt := &event.Event{
		Type:      event.ClassDeleted,
		Payload:   map[string]interface{}{"id": id},
		CreatedAt: time.Now(),
	}
	return uc.EventBus(evt)
}

func (uc *ClassUsecase) GetByID(id string) (*model.Class, error) {
	return uc.QueryRepo.GetByID(id)
}

func (uc *ClassUsecase) GetAll() ([]*model.Class, error) {
	return uc.QueryRepo.GetAll()
}

// Các hàm Update, Delete, GetByID, GetAll tương tự, publish event khi cần
