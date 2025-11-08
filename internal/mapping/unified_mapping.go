package mapping

import (
	"encoding/json"
	"qlsvgo/internal/domain/event"
	"qlsvgo/internal/domain/model"
	"reflect"
	"time"
)

// EntityMapper provides unified mapping between entities and events
type EntityMapper struct{}

// NewEntityMapper creates a new entity mapper
func NewEntityMapper() *EntityMapper {
	return &EntityMapper{}
}

// ToEvent converts any entity to event
func (em *EntityMapper) ToEvent(entity interface{}, eventType string) *event.Event {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	payload := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get JSON tag name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = fieldType.Name
		}

		// Skip empty JSON tags
		if jsonTag == "-" {
			continue
		}

		// Convert field value to interface{}
		var value interface{}
		switch field.Kind() {
		case reflect.String:
			value = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = field.Int()
		case reflect.Float32, reflect.Float64:
			value = field.Float()
		case reflect.Bool:
			value = field.Bool()
		case reflect.Struct:
			if field.Type() == reflect.TypeOf(time.Time{}) {
				value = field.Interface().(time.Time).Format(time.RFC3339)
			} else {
				value = field.Interface()
			}
		default:
			value = field.Interface()
		}

		payload[jsonTag] = value
	}

	return &event.Event{
		Type:    eventType,
		Payload: payload,
	}
}

// FromEvent converts event to entity
func (em *EntityMapper) FromEvent(evt *event.Event, entityType reflect.Type) interface{} {
	// Create new instance of entity
	entity := reflect.New(entityType).Interface()

	// Convert payload to JSON and back to entity
	payloadBytes, err := json.Marshal(evt.Payload)
	if err != nil {
		return nil
	}

	if err := json.Unmarshal(payloadBytes, entity); err != nil {
		return nil
	}

	return entity
}

// Specific mapping functions for type safety
func StudentToEvent(s *model.Student, eventType string) *event.Event {
	mapper := NewEntityMapper()
	return mapper.ToEvent(s, eventType)
}

func EventToStudent(evt *event.Event) *model.Student {
	mapper := NewEntityMapper()
	result := mapper.FromEvent(evt, reflect.TypeOf(model.Student{}))
	if student, ok := result.(*model.Student); ok {
		return student
	}
	return nil
}

func TeacherToEvent(t *model.Teacher, eventType string) *event.Event {
	mapper := NewEntityMapper()
	return mapper.ToEvent(t, eventType)
}

func EventToTeacher(evt *event.Event) *model.Teacher {
	mapper := NewEntityMapper()
	result := mapper.FromEvent(evt, reflect.TypeOf(model.Teacher{}))
	if teacher, ok := result.(*model.Teacher); ok {
		return teacher
	}
	return nil
}

func CourseToEvent(c *model.Course, eventType string) *event.Event {
	mapper := NewEntityMapper()
	return mapper.ToEvent(c, eventType)
}

func EventToCourse(evt *event.Event) *model.Course {
	mapper := NewEntityMapper()
	result := mapper.FromEvent(evt, reflect.TypeOf(model.Course{}))
	if course, ok := result.(*model.Course); ok {
		return course
	}
	return nil
}

func ClassToEvent(cl *model.Class, eventType string) *event.Event {
	mapper := NewEntityMapper()
	return mapper.ToEvent(cl, eventType)
}

func EventToClass(evt *event.Event) *model.Class {
	mapper := NewEntityMapper()
	result := mapper.FromEvent(evt, reflect.TypeOf(model.Class{}))
	if class, ok := result.(*model.Class); ok {
		return class
	}
	return nil
}

func RegistrationToEvent(r *model.Registration, eventType string) *event.Event {
	mapper := NewEntityMapper()
	return mapper.ToEvent(r, eventType)
}

func EventToRegistration(evt *event.Event) *model.Registration {
	mapper := NewEntityMapper()
	result := mapper.FromEvent(evt, reflect.TypeOf(model.Registration{}))
	if registration, ok := result.(*model.Registration); ok {
		return registration
	}
	return nil
}

func UserToEvent(u *model.User, eventType string) *event.Event {
	mapper := NewEntityMapper()
	return mapper.ToEvent(u, eventType)
}

func EventToUser(evt *event.Event) *model.User {
	mapper := NewEntityMapper()
	result := mapper.FromEvent(evt, reflect.TypeOf(model.User{}))
	if user, ok := result.(*model.User); ok {
		return user
	}
	return nil
}

// Serialization functions
func SerializeEvent(evt *event.Event) ([]byte, error) {
	return json.Marshal(evt)
}

func DeserializeEvent(data []byte) (*event.Event, error) {
	var evt event.Event
	err := json.Unmarshal(data, &evt)
	return &evt, err
}
