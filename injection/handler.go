package injection

import (
	"github.com/gin-gonic/gin"

	"github.com/leonardchinonso/auth_service_cmp7174/config"
	"github.com/leonardchinonso/auth_service_cmp7174/handler"
)

// injectHandlers initializes the dependencies and creates them as a config
func injectHandlers(router *gin.Engine, cfg *map[string]string, handlerCfg *HandlerConfig) {
	// get the current version number for correct routing
	version := (*cfg)[config.Version]

	// initialize the handlers
	handler.InitAuthHandler(router, version, handlerCfg.ClientService, handlerCfg.TokenService)
	handler.InitClientHandler(router, version, handlerCfg.ClientService, handlerCfg.TokenService)
}
