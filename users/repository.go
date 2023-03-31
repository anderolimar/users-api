package users

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection string = "users"
const INVALID_OBJECT_ID string = "INVALID_OBJECT_ID"

type UserRepository interface {
	InsertUser(user User) (string, error)
	FindUserByEmail(email string, projection Projection) (*User, error)
	FindUserByID(ID string, projection Projection) (*User, error)
	UpdateUser(userID string, user User) error
	DeleteUser(userID string) error
}

type userRepository struct {
	client   *mongo.Client
	database string
}

func NewUserRepository(client *mongo.Client, database string) UserRepository {
	return &userRepository{
		client:   client,
		database: database,
	}
}

func (repo *userRepository) InsertUser(user User) (string, error) {
	user.ID = ""
	coll := repo.client.Database(repo.database).Collection(userCollection)
	result, err := coll.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}
	var objID primitive.ObjectID = result.InsertedID.(primitive.ObjectID)

	return objID.Hex(), nil
}

func (repo *userRepository) FindUserByID(ID string, projection Projection) (*User, error) {
	coll := repo.client.Database(repo.database).Collection(userCollection)
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf(INVALID_OBJECT_ID)
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	opts := options.FindOne().SetProjection(projection.toBSON())

	var user User
	err = coll.FindOne(context.Background(), filter, opts).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
	}
	return &user, err
}

func (repo *userRepository) FindUserByEmail(email string, projection Projection) (*User, error) {
	coll := repo.client.Database(repo.database).Collection(userCollection)
	filter := bson.D{{Key: "email", Value: email}}
	opts := options.FindOne().SetProjection(projection.toBSON())

	var user User
	err := coll.FindOne(context.Background(), filter, opts).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
	}
	return &user, err
}

func (repo *userRepository) UpdateUser(userID string, user User) error {
	user.ID = ""
	coll := repo.client.Database(repo.database).Collection(userCollection)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf(INVALID_OBJECT_ID)
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	fields := bson.M{"$set": user.projection().toBSON()}

	_, err = coll.UpdateOne(context.Background(), filter, fields)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) DeleteUser(userID string) error {
	coll := repo.client.Database(repo.database).Collection(userCollection)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf(INVALID_OBJECT_ID)
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	_, err = coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

type Projection []ProjectionsFields

func (d Projection) Map() ProjectionMap {
	m := make(ProjectionMap, len(d))
	for _, e := range d {
		m[e.Key] = e.Value
	}
	return m
}

func (d Projection) toBSON() primitive.M {
	m := make(primitive.M, len(d))
	for _, e := range d {
		m[e.Key] = e.Value
	}
	return m
}

type ProjectionsFields struct {
	Key   string
	Value interface{}
}
type ProjectionMap map[string]interface{}
