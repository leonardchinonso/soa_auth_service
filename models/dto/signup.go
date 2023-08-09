package dto

import (
	"fmt"

	"github.com/leonardchinonso/auth_service_cmp7174/utils"
)

// SignupRequest holds the data for signup information
type SignupRequest struct {
	Name       string   `json:"name"`
	Email           Email    `json:"email"`
	Address string `json:"address"`
	Password        Password `json:"password"`
	ConfirmPassword Password `json:"confirm_password"`
	BusinessType BusinessType `json:"business_type"`
	ApiKey string `json:"api_key"`
}

// Validate validates an incoming signup request
func (sr *SignupRequest) Validate() []error {
	var errs []error

	utils.ShouldBePresentString(sr.Name, "name", &errs)
	utils.ShouldBePresentString(string(sr.Email), "email", &errs)
	utils.ShouldBePresentString(sr.Address, "address", &errs)
	utils.ShouldBePresentString(string(sr.Password), "password", &errs)
	utils.ShouldBePresentString(string(sr.ConfirmPassword), "confirmed password", &errs)
	utils.ShouldBePresentString(string(sr.BusinessType), "business type", &errs)
	utils.ShouldBePresentString(sr.ApiKey, "api key", &errs)

	// validate the email
	if err := sr.Email.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("email is invalid"))
	}

	// validate the business type
	if len(sr.BusinessType) > 0 {
		if err := sr.BusinessType.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("business type is invalid"))
		}
	}

	// validate the password
	if err := sr.Password.Validate(); err != nil {
		errs = append(errs, err)
	} else if ok := sr.Password.IsEqualValue(sr.ConfirmPassword); !ok {
		errs = append(errs, fmt.Errorf("passwords do not match"))
	}

	return errs
}
