package usecase

import (
	"qlsvgo/internal/domain/event"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/mapping"
	"qlsvgo/internal/repository"
	"time"

	"github.com/google/uuid"
)

type StudentUsecase struct {
	CommandRepo repository.StudentCommandRepository
	QueryRepo   repository.StudentQueryRepository
	EventBus    func(evt *event.Event) error // publish event
}

func (uc *StudentUsecase) Create(student *model.Student) error {
	if student.ID == "" {
		student.ID = uuid.New().String()
	}
	err := uc.CommandRepo.Create(student)
	if err != nil {
		return err
	}
	// Publish event
	evt := mapping.StudentToEvent(student, event.StudentCreated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *StudentUsecase) Update(student *model.Student) error {
	err := uc.CommandRepo.Update(student)
	if err != nil {
		return err
	}
	evt := mapping.StudentToEvent(student, event.StudentUpdated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *StudentUsecase) Delete(id string) error {
	err := uc.CommandRepo.Delete(id)
	if err != nil {
		return err
	}
	evt := &event.Event{
		Type:      event.StudentDeleted,
		Payload:   map[string]interface{}{"id": id},
		CreatedAt: time.Now(),
	}
	return uc.EventBus(evt)
}

func (uc *StudentUsecase) GetByID(id string) (*model.Student, error) {
	return uc.QueryRepo.GetByID(id)
}

func (uc *StudentUsecase) GetAll() ([]*model.Student, error) {
	return uc.QueryRepo.GetAll()
}
