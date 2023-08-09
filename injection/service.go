package injection

import (
	"github.com/leonardchinonso/auth_service_cmp7174/models/interfaces"
	"github.com/leonardchinonso/auth_service_cmp7174/service"
)

// HandlerConfig holds the configuration values for initializing the handlers
type HandlerConfig struct {
	ClientService             interfaces.ClientServiceInterface
	TokenService            interfaces.TokenServiceInterface
}

// injectServices initializes the dependencies and creates them as a config for handler injection
func injectServices(cfg *map[string]string, servCfg *ServicesConfig) (*HandlerConfig, error) {
	// initialize the client service with the needed config
	clientService := service.NewClientService(servCfg.ClientRepo, servCfg.TokenRepo)

	// initialize the token service with the needed config
	tokenService, err := service.NewTokenService(cfg, servCfg.TokenRepo)
	if err != nil {
		return nil, err
	}

	return &HandlerConfig{
		ClientService:             clientService,
		TokenService:            tokenService,
	}, nil
}
