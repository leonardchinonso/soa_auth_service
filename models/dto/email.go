package dto

import "net/mail"

// Email is a custom type for managing emails
type Email string

// Validate checks that an email meets the email requirements
func (e Email) Validate() error {
	if _, err := mail.ParseAddress(string(e)); err != nil {
		return err
	}
	return nil
}
