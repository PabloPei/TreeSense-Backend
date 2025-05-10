package users

import (
	"github.com/PabloPei/TreeSense-Backend/internal/auth"
	"github.com/PabloPei/TreeSense-Backend/internal/errors"
)

type Service struct {
	repository UserRepository
}

func NewService(repository UserRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) RegisterUser(payload RegisterUserPayload) error {

	_, err := s.repository.GetUserByEmail(payload.Email)
	if err == nil {
		return errors.ErrUserAlreadyExist(payload.Email)
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return errors.ErrHashingPassword(err)
	}

	user := User{
		UserName: payload.UserName,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	return s.repository.CreateUser(user)
}

func (s *Service) LogInUser(user LogInUserPayload) (string, string, error) {

	u, err := s.repository.GetUserByEmail(user.Email)

	if err != nil {
		return "", "", errors.ErrInvalidCredentials
	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		return "", "", errors.ErrInvalidCredentials
	}

	userJWT := createJWTPayload(*u)

	token, err := auth.CreateJWT(userJWT, false)
	if err != nil {
		return "", "", errors.ErrJWTCreation
	}

	refreshToken, err := auth.CreateJWT(userJWT, true)
	if err != nil {
		return "", "", errors.ErrJWTCreation
	}

	return token, refreshToken, nil
}

func (s *Service) GetUserPublicByEmail(email string) (*UserPublicPayload, error) {

	u, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return &UserPublicPayload{
		UserId:   u.UserId,
		Email:    u.Email,
		UserName: u.UserName,
		LanguageCode: u.LanguageCode,
		Photo: u.Photo,
	}, nil
}


func (s *Service) GetUserPublicById(userId []uint8) (*UserPublicPayload, error) {

	u, err := s.repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return &UserPublicPayload{
		UserId:   u.UserId,
		Email:    u.Email,
		UserName: u.UserName,
		LanguageCode: u.LanguageCode,
		Photo: u.Photo,
	}, nil
}

func (s *Service) UserExist(userId []uint8) (bool, error) {

	_, err := s.repository.GetUserById(userId)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *Service) RefreshToken(userId []uint8) (string, error) {

	user, err := s.repository.GetUserById(userId)

	if err != nil {
		return "", err
	}

	userJWT := createJWTPayload(*user)

	accessToken, err := auth.CreateJWT(userJWT, false)

	if err != nil {
		return "", errors.ErrJWTCreation
	}

	return accessToken, nil
}

func (s *Service) UploadPhoto(payload UploadPhotoPayload, email string) error {

	_, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return errors.ErrUploadPhoto
	}

	return s.repository.UploadPhoto(payload.Photo, email)
}

// Aux Functions

func createJWTPayload(user User) auth.UserJWT {

	var userJWT auth.UserJWT

	userJWT.UserId = string(user.UserId)
	userJWT.Email = user.Email
	userJWT.UserName = user.UserName

	return userJWT

}
