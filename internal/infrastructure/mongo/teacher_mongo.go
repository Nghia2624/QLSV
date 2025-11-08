package mongo

import (
	"context"
	"qlsvgo/internal/domain/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeacherMongoRepository struct {
	Collection *mongo.Collection
}

func (r *TeacherMongoRepository) GetByID(id string) (*model.Teacher, error) {
	var t model.Teacher
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TeacherMongoRepository) GetAll() ([]*model.Teacher, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var teachers []*model.Teacher
	for cursor.Next(context.Background()) {
		var t model.Teacher
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		teachers = append(teachers, &t)
	}
	return teachers, nil
}

func (r *TeacherMongoRepository) FindByEmail(email string) (*model.Teacher, error) {
	var t model.Teacher
	err := r.Collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
