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

// AuthHandler handles authentication related requests
type AuthHandler struct {
	clientService  interfaces.ClientServiceInterface
	tokenService interfaces.TokenServiceInterface
}

// InitAuthHandler initializes and sets up the auth handler
func InitAuthHandler(router *gin.Engine, version string, clientService interfaces.ClientServiceInterface, tokenService interfaces.TokenServiceInterface) {
	h := &AuthHandler{
		clientService:  clientService,
		tokenService: tokenService,
	}

	// group routes according to paths
	path := fmt.Sprintf("%s%s", version, "/auth")
	g := router.Group(path)

	// register endpoints
	g.POST("/signup", h.Signup)
	g.POST("/login", h.Login)
	g.POST("/logout", middlewares.AuthorizeClient(h.tokenService), h.Logout)
}

// Signup handles the incoming signup request
func (ah *AuthHandler) Signup(c *gin.Context) {
	var sr dto.SignupRequest

	// fill the signup request from binding the JSON request
	if err := c.ShouldBindJSON(&sr); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// validate the signup request for invalid fields
	if errs := sr.Validate(); len(errs) > 0 {
		resErr := errors.ErrBadRequest("invalid signup request", errors.ErrorToStringSlice(errs))
		c.JSON(resErr.Status, resErr)
		return
	}

	// create ah new client object with the details
	client := dao.NewClient(sr.Name, string(sr.Email), sr.Address, string(sr.Password), string(sr.BusinessType), sr.ApiKey)

	// start the signup process
	clientId, err := ah.clientService.Signup(c, client, sr.Password)
	if err != nil {
		log.Printf("Failed to sign client up. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), err)
		return
	}

	// retrieve the new client object from the database
	client, err = ah.clientService.GetClientByID(c, clientId)
	if err != nil {
		log.Printf("Failed to get client from database. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), err)
		return
	}

	fmt.Printf("Client retrieved: %+v\n", client)

	// create the access and refresh token pairs
	at, rt, err := ah.tokenService.GenerateTokenPair(c, client)
	if err != nil {
		log.Printf("Failed to generate client token pair. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), err)
		return
	}

	// create signup (login) response and return it to the handler's caller
	loginResp := dto.NewLoginResponse(*client, at, rt)
	resp := utils.ResponseStatusCreated("signed up successfully", loginResp)

	c.JSON(resp.Status, resp)
}

// Login handles the incoming login request
func (ah *AuthHandler) Login(c *gin.Context) {
	var lr dto.LoginRequest

	// fill the login request from binding the JSON request
	if err := c.ShouldBindJSON(&lr); err != nil {
		log.Printf("Failed to bind JSON with request. Error: %v\n", err)
		resErr := errors.ErrBadRequest(err.Error(), nil)
		c.JSON(resErr.Status, resErr)
		return
	}

	// validate the login request for invalid fields
	if errs := lr.Validate(); len(errs) > 0 {
		resErr := errors.ErrBadRequest("invalid login request", errs)
		c.JSON(resErr.Status, resErr)
		return
	}

	// create ah new client object with the details
	client := &dao.Client{Email: string(lr.Email), Password: string(lr.Password)}

	// start the login process
	err := ah.clientService.Login(c, client, lr.Password)
	if err != nil {
		log.Printf("Failed to login client. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), err)
		return
	}

	// create the access and refresh token pairs
	at, rt, err := ah.tokenService.GenerateTokenPair(c, client)
	if err != nil {
		log.Printf("Failed to generate client token pair. Error: %v\n", err.Error())
		c.JSON(errors.Status(err), err)
		return
	}

	// create ah login response and return it to the handler's caller
	loginResp := dto.NewLoginResponse(*client, at, rt)
	resp := utils.ResponseStatusCreated("logged in successfully", loginResp)

	c.JSON(resp.Status, resp)
}

// Logout handles the incoming logout request
func (ah *AuthHandler) Logout(c *gin.Context) {
	// retrieve the logged-in client from the authenticated request
	client, ok := ClientFromRequest(c)
	if !ok {
		log.Printf("Failed to retrieve client from authenticated request")
		resErr := errors.ErrUnauthorized("you are not logged in", nil)
		c.JSON(resErr.Status, gin.H{"errors": resErr})
		return
	}

	// attempt to log the client out
	err := ah.clientService.Logout(c, client.Id)
	if err != nil {
		c.JSON(errors.Status(err), err)
		return
	}

	resp := utils.ResponseStatusOK("logged out successfully", nil)
	c.JSON(resp.Status, resp)
}
