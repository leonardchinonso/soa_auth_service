package injection

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leonardchinonso/auth_service_cmp7174/models/interfaces"
	"github.com/leonardchinonso/auth_service_cmp7174/repository"
)

// ServicesConfig is the custom type for starting up services
type ServicesConfig struct {
	ClientRepo             interfaces.ClientRepositoryInterface
	TokenRepo            interfaces.TokenRepositoryInterface
}

// injectRepositories initializes the dependencies and creates them as a config for services injection
func injectRepositories(db *mongo.Database) *ServicesConfig {
	return &ServicesConfig{
		ClientRepo:             repository.NewClientRepository(db),
		TokenRepo:            repository.NewTokenRepository(db),
	}
}
