package roles

import (
	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/internal/users"
)

type Service struct {
	repository RoleRepository
	userRepository users.UserRepository
}

func NewService(repository RoleRepository, userRepository users.UserRepository) *Service {
	return &Service{repository: repository, userRepository: userRepository}
}

func (s *Service) CreateRole(payload CreateRolePayload) error {

	_, err := s.repository.GetRoleByName(payload.RoleName)

	if err == nil {
		return errors.ErrRoleAlreadyExist(payload.RoleName)
	}

	role := Role{
		RoleName: 			payload.RoleName,
		RoleDescription:    payload.RoleDescription,
	}

	return s.repository.CreateRole(role)
}


func (s *Service) GetRoles() ([]Role, error) {

	return s.repository.GetRoles()

}

//TODO cambiar a role service envez de repository
func (s *Service) CreateRoleAssigment(payload CreateUserRoleAssigmentPayload, email string, by []uint8) error {

	role, err := s.repository.GetRoleByName(payload.RoleName)

	if err != nil {
		return errors.ErrRoleNotFound
	}

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return errors.ErrUserNotFound
	}

	return s.repository.CreateRoleAssigment(user.UserId, role.RoleId, by, payload.ValidUntil)
}