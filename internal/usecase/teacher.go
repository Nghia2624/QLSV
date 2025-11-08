package usecase

import (
	"qlsvgo/internal/domain/event"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/mapping"
	"qlsvgo/internal/repository"
	"time"

	"github.com/google/uuid"
)

type TeacherUsecase struct {
	CommandRepo repository.TeacherCommandRepository
	QueryRepo   repository.TeacherQueryRepository
	EventBus    func(evt *event.Event) error // publish event
}

func (uc *TeacherUsecase) Create(teacher *model.Teacher) error {
	if teacher.ID == "" {
		teacher.ID = uuid.New().String()
	}
	err := uc.CommandRepo.Create(teacher)
	if err != nil {
		return err
	}
	evt := mapping.TeacherToEvent(teacher, event.TeacherCreated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *TeacherUsecase) Update(teacher *model.Teacher) error {
	err := uc.CommandRepo.Update(teacher)
	if err != nil {
		return err
	}
	evt := mapping.TeacherToEvent(teacher, event.TeacherUpdated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *TeacherUsecase) Delete(id string) error {
	err := uc.CommandRepo.Delete(id)
	if err != nil {
		return err
	}
	evt := &event.Event{
		Type:      event.TeacherDeleted,
		Payload:   map[string]interface{}{"id": id},
		CreatedAt: time.Now(),
	}
	return uc.EventBus(evt)
}

func (uc *TeacherUsecase) GetByID(id string) (*model.Teacher, error) {
	return uc.QueryRepo.GetByID(id)
}

func (uc *TeacherUsecase) GetAll() ([]*model.Teacher, error) {
	return uc.QueryRepo.GetAll()
}
