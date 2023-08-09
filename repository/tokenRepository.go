package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
	"github.com/leonardchinonso/auth_service_cmp7174/models/interfaces"
)

type tokenRepo struct {
	c *mongo.Collection
}

const tokenCollectionName = "tokens"

// NewTokenRepository returns a token interface with all the model repository methods
func NewTokenRepository(db *mongo.Database) interfaces.TokenRepositoryInterface {
	return &tokenRepo{
		c: db.Collection(tokenCollectionName),
	}
}

// Upsert updates a token by the clientId if it exists in the database
// it inserts a new document if it does not exist
func (tr *tokenRepo) Upsert(ctx context.Context, token *dao.Token) error {
	filter := bson.D{{"client_id", token.ClientId}}
	update := bson.D{{"$set", bson.D{{"client_id", token.ClientId}, {"refresh_token", token.RefreshToken}, {"access_token", token.AccessToken}}}}
	opts := options.Update().SetUpsert(true)
	_, err := tr.c.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a token from the token collection
func (tr *tokenRepo) Delete(ctx context.Context, clientId primitive.ObjectID) error {
	filter := bson.D{{"client_id", clientId}}
	_, err := tr.c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
