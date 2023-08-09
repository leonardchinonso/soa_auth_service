package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
)

// ClientFromRequest gets a client set by the authentication middleware
func ClientFromRequest(c *gin.Context) (*dao.Client, bool) {
	u, ok := c.Get("client")
	if !ok {
		return nil, false
	}

	// convert the type to ah client type
	client := u.(*dao.Client)
	return client, true
}
