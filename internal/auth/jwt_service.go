package auth

import (
	"context"
	"time"

	"github.com/PabloPei/SmartSpend-backend/conf"
	"github.com/PabloPei/SmartSpend-backend/internal/errors"
	"github.com/PabloPei/SmartSpend-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type UserJWT struct {
	UserId   string
	Email    string
	UserName string
}

func CreateJWT(user UserJWT, refreshToken bool) (string, error) {

	var expirationTime int64
	var secret []byte

	if refreshToken {
		secret = []byte(conf.ServerConfig.RefreshTokenSecret)
		expiration := time.Duration(conf.ServerConfig.RefreshTokenExpirationInHours) * time.Hour
		expirationTime = time.Now().UTC().Add(expiration).Unix()
	} else {
		secret = []byte(conf.ServerConfig.JWTSecret)
		expiration := time.Duration(conf.ServerConfig.JWTExpirationInSeconds) * time.Second
		expirationTime = time.Now().UTC().Add(expiration).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    user.UserId,
		"email":     user.Email,
		"userName":  user.UserName,
		"expiresAt": expirationTime,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, refreshToken bool) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrSignMethod(token.Header["alg"].(string))
		}

		if refreshToken {
			return []byte(conf.ServerConfig.RefreshTokenSecret), nil
		} else {
			return []byte(conf.ServerConfig.JWTSecret), nil
		}

	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrJWTInvalidToken
	}

	exp, ok := claims["expiresAt"].(float64)
	if !ok {
		return nil, errors.ErrJWTInvalidToken
	}

	expirationTime := int64(exp)
	currentTime := time.Now().UTC().Unix()
	if currentTime > expirationTime {
		return nil, errors.ErrJWTTokenExpired
	}

	return claims, nil
}

func GetUserIDFromContext(ctx context.Context) ([]uint8, error) {

	userId, ok := ctx.Value(models.UserKey).(string)

	userIdUint := []uint8(userId)

	if !ok {
		return nil, errors.ErrJWTInvalidToken
	}

	return userIdUint, nil
}
