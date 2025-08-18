package services

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	secretKey  string
	clientCode string
	expiresIn  time.Duration
}

func NewAuthService(secretKey, clientCode string, expiresIn time.Duration) *AuthService {
	return &AuthService{
		secretKey:  secretKey,
		clientCode: clientCode,
		expiresIn:  expiresIn,
	}
}

func (s *AuthService) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"client_code": s.clientCode,
		"exp":         time.Now().Add(s.expiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if code, ok := claims["client_code"].(string); ok && code == s.clientCode {
			return true, nil
		}
	}

	return false, nil
}
