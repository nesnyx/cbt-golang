package utils

import (
	"fmt"
	"github.com/gorilla/securecookie"
)

var (
	hashKey  = []byte("3a82ff008b1e93676f09b735d20f7d9ba7209593f530a5fa03998e9a422c8485") // 32 atau 64 byte
	blockKey = []byte("3f2cd7075135703729c0a581170e0ec9")                                 // 16, 24, atau 32 byte
	s        = securecookie.New(hashKey, blockKey)
)

func SetEncryptCookies(token string) (string, error) {
	encoded, err := s.Encode("jwt_cookie", token)
	fmt.Println(encoded)
	if err != nil {
		return "", err
	}
	return encoded, nil
}

func GetEncryptCookies(token string) (string, error) {
	var jwtToken string
	err := s.Decode("jwt_cookie", token, &jwtToken)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
