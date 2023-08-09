package dto

import (
	"fmt"

	"github.com/leonardchinonso/auth_service_cmp7174/utils"
)

// EditProfileRequest holds the data for the edit profile information
type EditProfileRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Email       Email  `json:"email"`
	PhoneNumber string `json:"phone_number"`
	BusinessType BusinessType `json:"business_type"`
}

// Validate validates an incoming edit profile request
func (epr *EditProfileRequest) Validate() []error {
	var errs []error

	utils.ShouldBePresentString(epr.Name, "name", &errs)
	utils.ShouldBePresentString(epr.Address, "address", &errs)
	utils.ShouldBePresentString(string(epr.Email), "email", &errs)
	utils.ShouldBePresentString(string(epr.BusinessType), "business type", &errs)

	// validate the email
	if err := epr.Email.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("email is invalid"))
	}

	// validate the business type
	if len(epr.BusinessType) > 0 {
		if err := epr.BusinessType.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("business type is invalid"))
		}
	}

	return errs
}
