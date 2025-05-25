package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(uuid string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type JwtService struct {
}

var SECRET_KEY = []byte("secret")

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (s *JwtService) GenerateToken(uuid string) (string, error) {
	if uuid == "" {
		return "", errors.New("uuid and allowedToken cannot be empty")
	}
	expiration := time.Now().Add(48 * time.Hour).Unix()
	claim := jwt.MapClaims{}
	claim["uuid"] = uuid
	claim["exp"] = expiration
	//claim["name"] = name
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (s *JwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				return nil, errors.New("token has expired")
			}
		} else {
			return nil, errors.New("invalid expiration time in token")
		}
	} else {
		return nil, errors.New("invalid token claims")
	}
	return token, nil

}
