package permission

import (
	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/internal/users"
)

type Service struct {
	repository PermissionRepository
	userRepository users.UserRepository
}

func NewService(repository PermissionRepository, userRepository users.UserRepository) *Service {
	return &Service{repository: repository, userRepository: userRepository}
}

/// Permissions /// 
func (s *Service) GetUserPermissions(email string) ([]PermissionAssignment, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return s.repository.GetUserPermissions(user.UserId)
}

func (s *Service) GetCurrentUserPermissions(userId []uint8) ([]PermissionAssignment, error) {
	return s.repository.GetUserPermissions(userId)
}

func (s *Service) UserHasPermissions(permissionNames []string, userId []uint8) (bool, error) {

	if len(permissionNames) == 0 {
		return true, nil 
	}


	userPermissions, err := s.repository.GetUserPermissions(userId)
	if err != nil {
		return false, err
	}

	userPermissionMap := make(map[string]bool)
	for _, userPermission := range userPermissions {
		userPermissionMap[userPermission.PermissionName] = true
	}

	for _, permissionName := range permissionNames {
		if !userPermissionMap[permissionName] {
			return false, nil
		}
	}

	return true, nil
}
