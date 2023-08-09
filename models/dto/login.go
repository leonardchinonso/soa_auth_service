package dto

import (
	"fmt"

	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
	"github.com/leonardchinonso/auth_service_cmp7174/utils"
)

// LoginRequest holds the data for the login information
type LoginRequest struct {
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

// Validate validates an incoming login request
func (lr *LoginRequest) Validate() []error {
	var errs []error

	utils.ShouldBePresentString(string(lr.Email), "email", &errs)
	utils.ShouldBePresentString(string(lr.Password), "password", &errs)

	if err := lr.Email.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("email is invalid"))
	}

	return errs
}

// LoginResponse holds the data for login response
type LoginResponse struct {
	Client         dao.Client `json:"client"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

// NewLoginResponse returns a new LoginResponse
func NewLoginResponse(client dao.Client, accessToken, refreshToken string) *LoginResponse {
	return &LoginResponse{
		Client: dao.Client{
			Id:          client.Id,
			Name: client.Name,
			Email:       client.Email,
			Address: client.Address,
			BusinessType: client.BusinessType,
			PhoneNumber: client.PhoneNumber,
			AccountActive: client.AccountActive,
			ApiKey: client.ApiKey, // TODO: ask if should return to client
			CreatedAt: client.CreatedAt,
			UpdatedAt: client.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
