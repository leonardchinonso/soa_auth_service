package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Client is the client data access object
type Client struct {
	Id          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name    string              `json:"name" binding:"required" bson:"name"`
	Email       string              `json:"email" binding:"required" bson:"email"`
	Address       string              `json:"address" binding:"required" bson:"address"`
	PhoneNumber string              `json:"phone_number" bson:"phone_number"`
	Password    string              `json:"password,omitempty" binding:"required" bson:"password"`
	BusinessType       string              `json:"business_type" binding:"required" bson:"business_type"`
	ApiKey       string              `json:"api_key" binding:"required" bson:"api_key"`
	AccountActive bool              `json:"account_active" binding:"required" bson:"account_active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// NewClient formats the client details and creates a new client
func NewClient(name, email, address, password, businessType, apiKey string) *Client {
	caser := cases.Title(language.English)
	name = caser.String(name)

	return &Client{
		Name: name,
		Email:       email,
		Address: address,
		Password:    password,
		BusinessType: businessType,
		ApiKey: apiKey,
		AccountActive: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
