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
	ErrRoleAssigmentExist = errors.New("role assigment already exist")
	ErrUserNotFound       = errors.New("user not found")
	ErrTreeNotFound       = errors.New("tree not found")
	ErrTreeSpecieNotFound = errors.New("tree specie not found")
	ErrTreeStateNotFound  = errors.New("tree state not found")
	ErrRoleNotFound       = errors.New("role not found")
	ErrCantUploadRole     = func(err string) error {
		return fmt.Errorf("can´t upload role info: %v", err)
	}
	ErrCantUploadTree     = func(err string) error {
		return fmt.Errorf("can´t create tree: %v", err)
	}
	ErrCantUploadUser     = func(err string) error {
		return fmt.Errorf("can´t upload user info: %v", err)
	}
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
		return fmt.Errorf("error scanning user: %v", err)
	}
	ErrTreeScan = func(err string) error {
		return fmt.Errorf("error scanning tree: %v", err)
	}
	ErrReadingRole = func(err string) error {
		return fmt.Errorf("error reading role: %v", err)
	}
	ErrRoleScan = func(err string) error {
		return fmt.Errorf("error scaning role: %v", err)
	}
	ErrUserNotHaveRole = func(rol string) error {
		return fmt.Errorf("error user does not have role %s", rol)
	}
	ErrRoleAlreadyExist = func(role string) error {
		return fmt.Errorf("role with name %s already exists", role)
	}
)
