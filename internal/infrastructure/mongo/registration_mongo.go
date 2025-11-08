package mongo

import (
	"context"
	"qlsvgo/internal/domain/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegistrationMongoRepository struct {
	Collection *mongo.Collection
}

func (r *RegistrationMongoRepository) GetByID(id string) (*model.Registration, error) {
	var reg model.Registration
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&reg)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *RegistrationMongoRepository) GetAll() ([]*model.Registration, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var regs []*model.Registration
	for cursor.Next(context.Background()) {
		var reg model.Registration
		if err := cursor.Decode(&reg); err != nil {
			return nil, err
		}
		regs = append(regs, &reg)
	}
	return regs, nil
}

func (r *RegistrationMongoRepository) FindByStudentID(studentID string) ([]*model.Registration, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var regs []*model.Registration
	for cursor.Next(context.Background()) {
		var reg model.Registration
		if err := cursor.Decode(&reg); err != nil {
			return nil, err
		}
		regs = append(regs, &reg)
	}
	return regs, nil
}

func (r *RegistrationMongoRepository) FindByClassID(classID string) ([]*model.Registration, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{"class_id": classID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var regs []*model.Registration
	for cursor.Next(context.Background()) {
		var reg model.Registration
		if err := cursor.Decode(&reg); err != nil {
			return nil, err
		}
		regs = append(regs, &reg)
	}
	return regs, nil
}

func (r *RegistrationMongoRepository) FindByStatus(status string) ([]*model.Registration, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var regs []*model.Registration
	for cursor.Next(context.Background()) {
		var reg model.Registration
		if err := cursor.Decode(&reg); err != nil {
			return nil, err
		}
		regs = append(regs, &reg)
	}
	return regs, nil
}
