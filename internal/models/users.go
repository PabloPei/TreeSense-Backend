package models

import (
	"time"
)

type User struct {
	UserId       []uint8   `json:"userId"`
	UserName     string    `json:"userName"`
	PhotoUrl     string    `json:"photoUrl"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	LanguageCode string    `json:"languageCode"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(User) error
	UploadPhoto(photoUrl string, email string) error
	GetUserById(id []uint8) (*User, error)
}

type UserService interface {
	RegisterUser(payload RegisterUserPayload) error
	LogInUser(user LogInUserPayload) (string, string, error)
	GetUserPublicByEmail(email string) (*UserPublicPayload, error)
	RefreshToken(userId []uint8) (string, error)
	UploadPhoto(payload UploadPhotoPayload, email string) error
}

type ContextKey string

var UserKey ContextKey = "userId"

type RegisterUserPayload struct {
	UserName string `json:"userName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	PhotoUrl string `json:"photoUrl" validate:"omitempty,uri"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LogInUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UploadPhotoPayload struct {
	PhotoUrl string `json:"photoUrl" validate:"required,uri"`
}

type UserPublicPayload struct {
	UserName string  `json:"userName" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	UserId   []uint8 `json:"userId"`
}
