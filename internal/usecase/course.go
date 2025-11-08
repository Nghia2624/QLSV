package usecase

import (
	"qlsvgo/internal/domain/event"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/mapping"
	"qlsvgo/internal/repository"
	"time"

	"github.com/google/uuid"
)

type CourseUsecase struct {
	CommandRepo repository.CourseCommandRepository
	QueryRepo   repository.CourseQueryRepository
	EventBus    func(evt *event.Event) error // publish event
}

func (uc *CourseUsecase) Create(course *model.Course) error {
	if course.ID == "" {
		course.ID = uuid.New().String()
	}
	err := uc.CommandRepo.Create(course)
	if err != nil {
		return err
	}
	evt := mapping.CourseToEvent(course, event.CourseCreated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *CourseUsecase) Update(course *model.Course) error {
	err := uc.CommandRepo.Update(course)
	if err != nil {
		return err
	}
	evt := mapping.CourseToEvent(course, event.CourseUpdated)
	evt.CreatedAt = time.Now()
	return uc.EventBus(evt)
}

func (uc *CourseUsecase) Delete(id string) error {
	err := uc.CommandRepo.Delete(id)
	if err != nil {
		return err
	}
	evt := &event.Event{
		Type:      event.CourseDeleted,
		Payload:   map[string]interface{}{"id": id},
		CreatedAt: time.Now(),
	}
	return uc.EventBus(evt)
}

func (uc *CourseUsecase) GetByID(id string) (*model.Course, error) {
	return uc.QueryRepo.GetByID(id)
}

func (uc *CourseUsecase) GetAll() ([]*model.Course, error) {
	return uc.QueryRepo.GetAll()
}
