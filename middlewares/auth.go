package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/leonardchinonso/auth_service_cmp7174/errors"
	"github.com/leonardchinonso/auth_service_cmp7174/models/interfaces"
)

// AuthorizeClient reads the token from the request header and gets the logged-in client
func AuthorizeClient(ts interfaces.TokenServiceInterface) gin.HandlerFunc {
	// return a function to handle the middleware
	return func(c *gin.Context) {
		tokenArr := c.Request.Header["Token"]
		if len(tokenArr) == 0 || tokenArr[0] == "" {
			resErr := errors.ErrUnauthorized("empty token value", nil)
			c.JSON(resErr.Status, resErr)
			c.Abort()
			return
		}

		// split the authorization token by the Bearer token
		splitTokenStr := strings.Split(tokenArr[0], "Bearer ")
		if len(splitTokenStr) != 2 {
			resErr := errors.ErrUnauthorized("must provide Authorization header with format `Bearer {token}`", nil)
			c.JSON(resErr.Status, resErr)
			c.Abort()
			return
		}

		// get the client from the access token
		client, err := ts.ClientFromAccessToken(splitTokenStr[1])
		if err != nil {
			resErr := errors.ErrUnauthorized("sorry, you're not authorized for this request", nil)
			c.JSON(resErr.Status, resErr)
			c.Abort()
			return
		}

		c.Set("client", client)

		c.Next()
	}
}
