package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
	"github.com/leonardchinonso/auth_service_cmp7174/models/interfaces"
)

type clientRepo struct {
	c *mongo.Collection
}

const clientCollectionName = "clients"

// NewClientRepository returns a token interface with all the model repository methods
func NewClientRepository(db *mongo.Database) interfaces.ClientRepositoryInterface {
	return &clientRepo{
		c: db.Collection(clientCollectionName),
	}
}

// Create creates a new client document in the database
func (ur *clientRepo) Create(ctx context.Context, client *dao.Client) (primitive.ObjectID, error) {
	result, err := ur.c.InsertOne(ctx, client)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// FindByID finds a client by id in the database
func (ur *clientRepo) FindByID(ctx context.Context, client *dao.Client) (bool, error) {
	err := ur.c.FindOne(ctx, bson.M{"_id": client.Id}).Decode(client)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find client: %w", err)
	}
	return true, nil
}

// FindByEmail finds a client by email in the database
func (ur *clientRepo) FindByEmail(ctx context.Context, client *dao.Client) (bool, error) {
	err := ur.c.FindOne(ctx, bson.M{"email": client.Email}).Decode(client)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find client: %w", err)
	}
	return true, nil
}

// Update updates a client in the database
func (ur *clientRepo) Update(ctx context.Context, client *dao.Client) error {
	filter := bson.D{{Key: "_id", Value: client.Id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: client.Name},
		{Key: "email", Value: client.Email},
		{Key: "address", Value: client.Address},
		{Key: "phone_number", Value: client.PhoneNumber},
		{Key: "business_type", Value: client.BusinessType},
		{Key: "updated_at", Value: client.UpdatedAt},
	}}}
	return ur.updateByQuery(ctx, filter, update)
}

// updateByQuery updates a savedPlace by a specified query
func (ur *clientRepo) updateByQuery(ctx context.Context, filter primitive.D, update primitive.D) error {
	_, err := ur.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
