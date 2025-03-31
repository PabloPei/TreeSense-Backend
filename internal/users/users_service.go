package users

import (
	"github.com/PabloPei/TreeSense-Backend/internal/auth"
	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/internal/roles"
)

type Service struct {
	repository UserRepository
	roleRepository roles.RoleRepository
}

func NewService(repository UserRepository, roleRepository roles.RoleRepository) *Service {
	return &Service{repository: repository, roleRepository: roleRepository}
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
	}, nil
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

	return s.repository.UploadPhoto(payload.PhotoUrl, email)
}

//TODO cambiar a role service envez de repository
func (s *Service) CreateRoleAssigment(payload CreateUserRoleAssigmentPayload, email string, by []uint8) error {

	role, err := s.roleRepository.GetRoleByName(payload.RoleName)

	if err != nil {
		return errors.ErrRoleNotFound
	}

	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return errors.ErrUserNotFound
	}


	return s.repository.CreateRoleAssigment(user.UserId, role.RoleId, by, payload.ValidUntil)
}

// Aux Functions

func createJWTPayload(user User) auth.UserJWT {

	var userJWT auth.UserJWT

	userJWT.UserId = string(user.UserId)
	userJWT.Email = user.Email
	userJWT.UserName = user.UserName

	return userJWT

}
