package app

import (
	"errors"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (tokenString string, err error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	claims := &JWTClaim{
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
			ID:        user.ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return err
	}
	if claims.ExpiresAt.Time.Unix() < time.Now().Unix() {
		err = errors.New("token expired")
		return err
	}
	return
}

func GetData(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	if claims.ExpiresAt.Time.Unix() < time.Now().Unix() {
		err = errors.New("token expired")
		return nil, err
	}

	return claims, nil
}
