package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/leonardchinonso/auth_service_cmp7174/errors"
	"github.com/leonardchinonso/auth_service_cmp7174/middlewares"
	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
	"github.com/leonardchinonso/auth_service_cmp7174/models/dto"
	"github.com/leonardchinonso/auth_service_cmp7174/models/interfaces"
	"github.com/leonardchinonso/auth_service_cmp7174/utils"
)

// ClientHandler represents the router handler object for the client requests
type ClientHandler struct {
	clientService  interfaces.ClientServiceInterface
	tokenService interfaces.TokenServiceInterface
}

// InitClientHandler initializes the client handler
func InitClientHandler(router *gin.Engine, version string, clientService interfaces.ClientServiceInterface, tokenService interfaces.TokenServiceInterface) {
	h := &ClientHandler{
		clientService:  clientService,
		tokenService: tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/client")
	g := router.Group(path)

	g.PUT("/update-profile", middlewares.AuthorizeClient(h.tokenService), h.UpdateProfile)
}

// UpdateProfile handles the request to update client details
func (h *ClientHandler) UpdateProfile(c *gin.Context) {
	// retrieve the logged-in client from the authenticated request
	cl, ok := ClientFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve client from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

	var epr dto.EditProfileRequest
	// fill the edit profile request from binding the JSON request
	if err := c.ShouldBindJSON(&epr); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// validate the login request for invalid fields
	if errs := epr.Validate(); len(errs) > 0 {
		log.Printf("Failed to validate request. Errors: %+v", errs)
		resErr := errors.ErrBadRequest("invalid request", errors.ErrorToStringSlice(errs))
		c.JSON(resErr.Status, resErr)
		return
	}

	// create a new client and set their details
	client := dao.NewClient(epr.Name, string(epr.Email), epr.Address, "", string(epr.BusinessType), "")
	client.Id = cl.Id
	client.PhoneNumber = epr.PhoneNumber

	// start the signup process
	err := h.clientService.EditClientProfile(c, client)
	if err != nil {
		log.Printf("Failed to sign client up. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), err)
		return
	}

	resp := utils.ResponseStatusOK("profile edited successfully", client)
	c.JSON(resp.Status, resp)
}
