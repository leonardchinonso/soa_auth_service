package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Map map[string]string

const (
	// ATExpiresIn is the global config name for the AT_EXPIRES_IN variable
	ATExpiresIn = "AT_EXPIRES_IN"
	// RTExpiresIn is the global config name for the RT_EXPIRES_IN variable
	RTExpiresIn = "RT_EXPIRES_IN"
	// ATSecretKey is the global config name for the AT_SECRET_KEY variable
	ATSecretKey = "AT_SECRET_KEY"
	// RTSecretKey is the global config name for the RT_SECRET_KEY variable
	RTSecretKey = "RT_SECRET_KEY"

	// DatabaseName is the global config name for the DATABASE_NAME variable
	DatabaseName = "DATABASE_NAME"
	// LocalDatabaseUri  is the global config name for the LOCAL_DATABASE_URI variable
	LocalDatabaseUri = "LOCAL_DATABASE_URI"
	// AtlasDatabaseUri  is the global config name for the ATLAS_DATABASE_URI variable
	AtlasDatabaseUri = "ATLAS_DATABASE_URI"
	// Version is the global config name for the VERSION variable
	Version = "VERSION"
	// Environment is the global config name for the ENVIRONMENT variable
	Environment = "ENVIRONMENT"
)

// getEnv retrieves teh value of a given key from the environment variables set
func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("failed to get value for key: %v", key)
}

// InitConfig loads the config variables into the application and populates the config map with the variables
func InitConfig() (*map[string]string, error) {
	// load values from the dev.env file into the application
	if err := godotenv.Load("./config/dev.env"); err != nil {
		return nil, fmt.Errorf("failed to load config variables: %v", err)
	}

	Map = make(map[string]string)
	var defaultConfig = []string{
		DatabaseName, ATExpiresIn, RTExpiresIn, ATSecretKey, RTSecretKey, Version, LocalDatabaseUri, AtlasDatabaseUri, Environment,
	}

	// iterate the preset config variables, retrieve their values and set them in the config map
	for _, c := range defaultConfig {
		v, err := getEnv(c)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve env variables: %v", err)
		}
		Map[c] = v
	}

	return &Map, nil
}
