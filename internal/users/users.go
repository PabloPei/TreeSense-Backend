package users

import (
	"time"
)

type User struct {
	UserId       []uint8   `json:"userId"`
	UserName     string    `json:"userName"`
	Photo     string    `json:"photo"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	LanguageCode string    `json:"languageCode"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(User) error
	UploadPhoto(photo string, email string) error
	GetUserById(id []uint8) (*User, error)
}

type UserService interface {
	RegisterUser(payload RegisterUserPayload) error
	LogInUser(user LogInUserPayload) (string, string, error)
	GetUserPublicByEmail(email string) (*UserPublicPayload, error)
	RefreshToken(userId []uint8) (string, error)
	UploadPhoto(payload UploadPhotoPayload, email string) error
	UserExist(userId []uint8) (bool, error)
	GetUserPublicById(userId []uint8) (*UserPublicPayload, error)
}

type RegisterUserPayload struct {
	UserName string `json:"userName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Photo string `json:"photo" validate:"omitempty,base64"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LogInUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UploadPhotoPayload struct {
	Photo string `json:"photo" validate:"omitempty,base64"`
}

type UserPublicPayload struct {
	UserName string  `json:"userName" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	UserId   []uint8 `json:"userId"`
	Photo string `json:"photo"`
	LanguageCode string    `json:"languageCode"`
}
