package mongo

import (
	"context"
	"qlsvgo/internal/domain/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassMongoRepository struct {
	Collection *mongo.Collection
}

func (r *ClassMongoRepository) GetByID(id string) (*model.Class, error) {
	var cl model.Class
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&cl)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

func (r *ClassMongoRepository) GetAll() ([]*model.Class, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var classes []*model.Class
	for cursor.Next(context.Background()) {
		var cl model.Class
		if err := cursor.Decode(&cl); err != nil {
			return nil, err
		}
		classes = append(classes, &cl)
	}
	return classes, nil
}

func (r *ClassMongoRepository) FindByCourseID(courseID string) ([]*model.Class, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{"course_id": courseID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var classes []*model.Class
	for cursor.Next(context.Background()) {
		var cl model.Class
		if err := cursor.Decode(&cl); err != nil {
			return nil, err
		}
		classes = append(classes, &cl)
	}
	return classes, nil
}

func (r *ClassMongoRepository) FindByTeacherID(teacherID string) ([]*model.Class, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{"teacher_id": teacherID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var classes []*model.Class
	for cursor.Next(context.Background()) {
		var cl model.Class
		if err := cursor.Decode(&cl); err != nil {
			return nil, err
		}
		classes = append(classes, &cl)
	}
	return classes, nil
}

func (r *ClassMongoRepository) FindBySemester(semester string) ([]*model.Class, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{"semester": semester})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var classes []*model.Class
	for cursor.Next(context.Background()) {
		var cl model.Class
		if err := cursor.Decode(&cl); err != nil {
			return nil, err
		}
		classes = append(classes, &cl)
	}
	return classes, nil
}
