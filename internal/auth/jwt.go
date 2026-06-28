package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userID, username, email string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type JWT struct {
	secretKey []byte
	duration time.Duration
}

func NewJWTService(secretKey string, duration time.Duration) JWTService {

	return &JWT{
		secretKey: []byte(secretKey),
		duration: duration,
	}
}

func (j *JWT) GenerateToken(userID, username, email string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "spotsync",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}