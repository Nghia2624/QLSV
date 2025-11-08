package mongo

import (
	"context"
	"qlsvgo/internal/domain/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoRepository struct {
	Collection *mongo.Collection
}

func (r *UserMongoRepository) Create(user *model.User) error {
	_, err := r.Collection.InsertOne(context.Background(), user)
	return err
}

func (r *UserMongoRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.Collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserMongoRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserMongoRepository) GetAll() ([]*model.User, error) {
	cursor, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*model.User
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserMongoRepository) Update(user *model.User) error {
	_, err := r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func (r *UserMongoRepository) Delete(id string) error {
	_, err := r.Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
