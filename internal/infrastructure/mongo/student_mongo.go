package mongo

import (
	"context"
	"qlsvgo/internal/domain/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentMongoRepository struct {
	Collection *mongo.Collection
}

func (r *StudentMongoRepository) GetByID(id string) (*model.Student, error) {
	var s model.Student
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentMongoRepository) GetAll() ([]*model.Student, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var students []*model.Student
	for cursor.Next(context.Background()) {
		var s model.Student
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		students = append(students, &s)
	}
	return students, nil
}

func (r *StudentMongoRepository) FindByEmail(email string) (*model.Student, error) {
	var s model.Student
	err := r.Collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
