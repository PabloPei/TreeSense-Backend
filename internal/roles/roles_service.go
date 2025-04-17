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

/// Roles /// 
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

func (s *Service) GetUserRoles(email string) ([]RoleAssigment, error) {

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return s.repository.GetUserRoles(user.UserId)

}

func (s *Service) GetCurrentUserRoles(userId []uint8)([]RoleAssigment, error) {
	return s.repository.GetUserRoles(userId)
}

func (s *Service) UserHasRole(roleName string, userId []uint8)(bool, error){

	if roleName == "" {
		return true, nil 
	}
	
	role, err := s.repository.GetRoleByName(roleName)

	if err != nil {
		return false, errors.ErrRoleNotFound
	}


	userRoles, err := s.repository.GetUserRoles(userId)

	for _, userRole := range userRoles {
		if userRole.RoleName == role.RoleName {
			return true, nil
		}
	}

	return false, nil

}

/// Assigments /// 

func (s *Service) CreateRoleAssigment(payload CreateUserRoleAssigmentPayload, email string, by []uint8) error {

	role, err := s.repository.GetRoleByName(payload.RoleName)

	if err != nil {
		return errors.ErrRoleNotFound
	}

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return errors.ErrUserNotFound
	}

	userRoles, err := s.repository.GetUserRoles(user.UserId)

	for _, userRole := range userRoles {
		if userRole.RoleName == role.RoleName {
			return errors.ErrRoleAssigmentExist
		}
	}
	return s.repository.CreateRoleAssigment(user.UserId, role.RoleId, by, payload.ValidUntil)
}

func (s *Service) DeleteRoleAssigment(payload DeleteUserRoleAssigmentPayload, email string) error {

	role, err := s.repository.GetRoleByName(payload.RoleName)

	if err != nil {
		return errors.ErrRoleNotFound
	}

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return errors.ErrUserNotFound
	}

	userRoles, err := s.repository.GetUserRoles(user.UserId)

	roleAssigned := false
	for _, userRole := range userRoles {
		if userRole.RoleName == role.RoleName {
			roleAssigned = true
			break
		}
	}
	if !roleAssigned {
		return errors.ErrRoleAssigmentNotExist
	}

	return s.repository.DeleteRoleAssigment(user.UserId, role.RoleId)

}