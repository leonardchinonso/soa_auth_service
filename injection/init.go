package injection

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/leonardchinonso/auth_service_cmp7174/datasource"
)

// Inject injects all the repos and services necessary
func Inject(ds *datasource.DataSource) (*gin.Engine, error) {
	log.Printf("Injecting Data Sources...\n")

	// load repositories
	servCfg := injectRepositories(ds.Database)

	// load services
	handCfg, err := injectServices(ds.Cfg, servCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to inject services: %v", err)
	}

	// load router
	router := gin.Default()

	// load handlers
	injectHandlers(router, ds.Cfg, handCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to inject handlers: %v", err)
	}

	return router, nil
}