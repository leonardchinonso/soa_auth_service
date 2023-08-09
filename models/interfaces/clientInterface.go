package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
	"github.com/leonardchinonso/auth_service_cmp7174/models/dto"
)

// ClientRepositoryInterface defines methods that are associated with the client repository
type ClientRepositoryInterface interface {
	Create(ctx context.Context, client *dao.Client) (primitive.ObjectID, error)
	FindByID(ctx context.Context, client *dao.Client) (bool, error)
	FindByEmail(ctx context.Context, client *dao.Client) (bool, error)
	Update(ctx context.Context, client *dao.Client) error
}

// ClientServiceInterface defines methods that are associated with the client repository
type ClientServiceInterface interface {
	Signup(ctx context.Context, client *dao.Client, password dto.Password) (primitive.ObjectID, error)
	Login(ctx context.Context, client *dao.Client, password dto.Password) error
	Logout(ctx context.Context, clientId primitive.ObjectID) error
	GetClientByID(ctx context.Context, clientId primitive.ObjectID) (*dao.Client, error)
	EditClientProfile(ctx context.Context, client *dao.Client) error
}
