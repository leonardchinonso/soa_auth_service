package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/leonardchinonso/auth_service_cmp7174/utils"

	"github.com/golang-jwt/jwt"

	"github.com/leonardchinonso/auth_service_cmp7174/config"
	"github.com/leonardchinonso/auth_service_cmp7174/models/dao"
	"github.com/leonardchinonso/auth_service_cmp7174/models/interfaces"
)

type tokenService struct {
	tokenRepository interfaces.TokenRepositoryInterface
	atSecret        string
	rtSecret        string
	atExpiresIn     int64
	rtExpiresIn     int64
}

// NewTokenService returns an interface for the token service methods
func NewTokenService(cfg *map[string]string, tokenRepo interfaces.TokenRepositoryInterface) (interfaces.TokenServiceInterface, error) {
	atExpiresIn, err := strconv.Atoi((*cfg)[config.ATExpiresIn])
	if err != nil {
		return nil, err
	}

	rtExpiresIn, err := strconv.Atoi((*cfg)[config.RTExpiresIn])
	if err != nil {
		return nil, err
	}

	return &tokenService{
		tokenRepository: tokenRepo,
		atSecret:        (*cfg)[config.ATSecretKey],
		rtSecret:        (*cfg)[config.RTSecretKey],
		atExpiresIn:     int64(atExpiresIn),
		rtExpiresIn:     int64(rtExpiresIn),
	}, nil
}

// GenerateTokenPair generates an access token and a refresh token for the specified client
func (ts *tokenService) GenerateTokenPair(ctx context.Context, client *dao.Client) (string, string, error) {
	at, err := generateAccessToken(client)
	if err != nil {
		log.Printf("Error generating access token for uid: %v. Error: %v\n", client.Id, err.Error())
		return "", "", err
	}

	rt, err := generateRefreshToken(client)
	if err != nil {
		log.Printf("Error generating refresh token for uid: %v. Error: %v\n", client.Id, err.Error())
		return "", "", err
	}

	token := &dao.Token{
		ClientId:       client.Id,
		AccessToken:  at,
		RefreshToken: rt,
		CreatedAt:    utils.CurrentPrimitiveTime(),
	}

	if err = ts.tokenRepository.Upsert(ctx, token); err != nil {
		log.Printf("Error upserting token in database for uid: %v. Error: %v\n", client.Id, err.Error())
		return "", "", err
	}

	return at, rt, nil
}

// ClientFromAccessToken gets a client from their access token
func (ts *tokenService) ClientFromAccessToken(tokenString string) (*dao.Client, error) {
	claims, err := verifyAccessToken(tokenString, config.Map[config.ATSecretKey])

	if err != nil {
		log.Printf("Unable to validate or parse access token. Error: %v\n", err)
		return nil, fmt.Errorf("cannot authenticate client: %v", err)
	}

	return claims.Client, nil
}

type tokenCustomClaims struct {
	Client *dao.Client `json:"client"`
	jwt.StandardClaims
}

// generateToken generates a new jwt
func generateToken(client *dao.Client, jwtSecretKey string, expiresIn int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExpiresIn := unixTime + expiresIn

	// create a claims object
	claims := tokenCustomClaims{
		Client: client,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiresIn,
			IssuedAt:  unixTime,
		},
	}

	// create a jwt token object and set the expiry time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the token string
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		log.Printf("Error generating single token for clientId: %v. Error: %v\n", client.Id, err.Error())
		return "", err
	}

	return tokenString, nil
}

// generateAccessToken generates a new jwt for the access token
func generateAccessToken(client *dao.Client) (string, error) {
	// get the access token secret key for signing the token
	atExpiresIn, err := strconv.Atoi(config.Map[config.ATExpiresIn])
	if err != nil {
		return "", err
	}

	// get the access token secret key
	atSecretKey := config.Map[config.ATSecretKey]

	return generateToken(client, atSecretKey, int64(atExpiresIn))
}

// generateRefreshToken generates a new jwt for the refresh token
func generateRefreshToken(client *dao.Client) (string, error) {
	// get the refresh token secret key for signing the token
	rtExpiresIn, err := strconv.Atoi(config.Map[config.RTExpiresIn])
	if err != nil {
		return "", err
	}

	// get the refresh token secret key
	rtSecretKey := config.Map[config.RTSecretKey]

	return generateToken(client, rtSecretKey, int64(rtExpiresIn))
}

// verifyAccessToken verifies that an access token is correct
func verifyAccessToken(tokenString, atSecretKey string) (*tokenCustomClaims, error) {
	claims := &tokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(atSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*tokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}
