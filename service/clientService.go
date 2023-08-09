package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/leonardchinonso/auth_service_cmp7174/errors"
	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
	"github.com/leonardchinonso/auth_service_cmp7174/models/dto"
	"github.com/leonardchinonso/auth_service_cmp7174/models/interfaces"
)

type clientService struct {
	clientRepository  interfaces.ClientRepositoryInterface
	tokenRepository interfaces.TokenRepositoryInterface
}

// NewClientService returns an interface for the client service methods
func NewClientService(clientRepo interfaces.ClientRepositoryInterface, tokenRepo interfaces.TokenRepositoryInterface) interfaces.ClientServiceInterface {
	return &clientService{
		clientRepository:  clientRepo,
		tokenRepository: tokenRepo,
	}
}

// Signup handles the client creation and logs the client in
func (us *clientService) Signup(ctx context.Context, client *dao.Client, password dto.Password) (primitive.ObjectID, error) {
	// hash the password to hide its real value
	hashedPassword, err := password.Hash()
	if err != nil {
		log.Printf("Error hashing client password: %s. Error: %v\n", password, err.Error())
		return primitive.ObjectID{}, errors.ErrInternalServerError("failed to sign up client", err)
	}

	// update the client password to its hashed value
	client.Password = hashedPassword

	// check that the email is not taken
	clientExists, err := us.clientRepository.FindByEmail(ctx, client)
	if err != nil {
		log.Printf("Error finding client with email: %s. Error: %v\n", client.Email, err.Error())
		return primitive.ObjectID{}, errors.ErrInternalServerError("failed to fetch client details", err)
	}

	// if the email already exists, return an error saying the email is taken
	if clientExists {
		return primitive.ObjectID{}, errors.ErrBadRequest("sorry, email is taken", nil)
	}

	// create a new client with the credentials
	insertedId, err := us.clientRepository.Create(ctx, client)
	if err != nil {
		log.Printf("Error creating client with email: %s. Error: %v\n", client.Email, err.Error())
		return primitive.ObjectID{}, errors.ErrInternalServerError("failed to sign up client", err)
	}

	return insertedId, nil
}

// Login logs the client into the application and returns the authentication tokens
func (us *clientService) Login(ctx context.Context, client *dao.Client, password dto.Password) error {
	// find the client by email and password
	clientExists, err := us.clientRepository.FindByEmail(ctx, client)
	if err != nil { // if an unexpected error occurs
		log.Printf("Error finding client with email: %s. Error: %v\n", client.Email, err.Error())
		return errors.ErrInternalServerError("failed to fetch client details", err)
	}

	// if the client does not exist, then the password and/or email are wrong
	if !clientExists {
		return errors.ErrUnauthorized(errors.ErrInvalidLogin, nil)
	}

	// if the client exists, but the password is not correct
	if !password.IsEqualHash(client.Password) {
		return errors.ErrUnauthorized(errors.ErrInvalidLogin, nil)
	}

	return nil
}

// Logout logs the client out of the application
func (us *clientService) Logout(ctx context.Context, clientId primitive.ObjectID) error {
	token := &dao.Token{
		ClientId: clientId,
	}

	err := us.tokenRepository.Delete(ctx, token.ClientId)
	if err != nil {
		log.Printf("Error trying to delete token with clientId: %v. Error: %v\n", clientId, err.Error())
		return errors.ErrInternalServerError("failed to log client out", err)
	}

	return nil
}

// GetClientByID gets a client by their ID
func (us *clientService) GetClientByID(ctx context.Context, clientId primitive.ObjectID) (*dao.Client, error) {
	// check that the client id is not empty
	if clientId.IsZero() {
		return nil, errors.ErrBadRequest("invalid client id", nil)
	}

	// create a new client object
	client := &dao.Client{Id: clientId}

	// find the client by the objectId
	clientExists, err := us.clientRepository.FindByID(ctx, client)
	if err != nil {
		log.Printf("Error finding client with id: %s. Error: %v\n", client.Id, err.Error())
		return nil, errors.ErrInternalServerError("failed to retrieve client", nil)
	}

	// return if the client does not exist
	if !clientExists {
		return nil, errors.ErrBadRequest("client not found", nil)
	}

	return client, nil
}

func (us *clientService) EditClientProfile(ctx context.Context, client *dao.Client) error {
	// check that the client id is not empty
	if client.Id.IsZero() {
		log.Printf("Error validating client Id: %v\n", client.Id)
		return errors.ErrBadRequest("invalid client id", nil)
	}

	clientChecker := &dao.Client{Email: client.Email}

	// check that the email is not taken
	clientExists, err := us.clientRepository.FindByEmail(ctx, clientChecker)
	if err != nil {
		log.Printf("Error finding client with email: %s. Error: %v\n", client.Email, err.Error())
		return errors.ErrInternalServerError("failed to fetch client details", err)
	}

	// if the email already exists, return an error saying the email is taken
	if clientExists && client.Id != clientChecker.Id {
		return errors.ErrBadRequest("sorry, email is taken", nil)
	}

	// update the client with the new information
	err = us.clientRepository.Update(ctx, client)
	if err != nil {
		log.Printf("Error updating client with id: %v. Error: %v\n", client.Id, err.Error())
		return errors.ErrInternalServerError("failed to update client information", nil)
	}

	return nil
}
