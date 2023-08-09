package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
)

// TokenRepositoryInterface defines methods that are applicable to the token repository
type TokenRepositoryInterface interface {
	Upsert(ctx context.Context, token *dao.Token) error
	Delete(ctx context.Context, clientId primitive.ObjectID) error
}

// TokenServiceInterface defines methods that are applicable to the token service
type TokenServiceInterface interface {
	GenerateTokenPair(ctx context.Context, client *dao.Client) (string, string, error)
	ClientFromAccessToken(tokenString string) (*dao.Client, error)
}
