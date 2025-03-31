package roles

import (
	"github.com/PabloPei/TreeSense-Backend/internal/errors"
)

type Service struct {
	repository RoleRepository
}

func NewService(repository RoleRepository) *Service {
	return &Service{repository: repository}
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