package dto

import "fmt"

// BusinessType is a custom enum type
type BusinessType string

var (
	SoftwareEngineering BusinessType = "Software Engineering"
	Construction BusinessType = "Construction"
	InformationTechnology BusinessType = "Information Technology"
)

var BusinessTypes []BusinessType = []BusinessType{
	SoftwareEngineering, Construction, InformationTechnology,
}

// Validate checks if a business type is valid
func (bt BusinessType) Validate() error {
	for _, b := range BusinessTypes {
		if b == bt {
			return nil
		}
	}
	return fmt.Errorf("invalid business type: %v", bt)
}