package mongo

import (
	"context"
	"qlsvgo/internal/domain/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseMongoRepository struct {
	Collection *mongo.Collection
}

func (r *CourseMongoRepository) GetByID(id string) (*model.Course, error) {
	var c model.Course
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CourseMongoRepository) GetAll() ([]*model.Course, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var courses []*model.Course
	for cursor.Next(context.Background()) {
		var c model.Course
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}
	return courses, nil
}

func (r *CourseMongoRepository) FindByName(name string) ([]*model.Course, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{"name": name})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var courses []*model.Course
	for cursor.Next(context.Background()) {
		var c model.Course
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}
	return courses, nil
}
