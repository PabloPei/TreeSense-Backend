package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrJWTCreation        = errors.New("unable to create JWT token")
	ErrJWTInvalidToken    = errors.New("error authenticating user: Token not valid")
	ErrJWTTokenExpired    = errors.New("error authenticating user: JWT token expired")
	ErrUploadPhoto        = errors.New("unable to upload photo")
	ErrUserNotFound       = errors.New("user not found")
	ErrGroupNotFound      = errors.New("group not found")
	ErrPermissionDenied   = func(permission string) error {
		return fmt.Errorf("user do not have %v permissions", permission)
	}
	ErrInvalidaPayload = func(err string) error {
		return fmt.Errorf("invalid payload: %v", err)
	}
	ErrUserAlreadyExist = func(email string) error {
		return fmt.Errorf("user with email %s already exists", email)
	}
	ErrHashingPassword = func(hashError error) error {
		return fmt.Errorf("failed to hash password: %v", hashError)
	}
	ErrSignMethod = func(alg string) error {
		return fmt.Errorf("unexpected signing method: %v", alg)
	}
	ErrUserScan = func(err string) error {
		return fmt.Errorf("error scaning user: %v", err)
	}
	ErrCreateGroup = func(err string) error {
		return fmt.Errorf("group can´t be created: %v", err)
	}
)
