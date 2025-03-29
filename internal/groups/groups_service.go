package groups

import (
	"log"

	"github.com/PabloPei/SmartSpend-backend/internal/errors"
	"github.com/PabloPei/SmartSpend-backend/internal/models"
)

type Service struct {
	repository models.GroupRepository
}

func NewService(repository models.GroupRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateGroup(payload models.CreateGroupPayload, userId []uint8) error {

	group := models.Group{
		GroupName:   payload.GroupName,
		Description: payload.Description,
		PhotoUrl:    payload.PhotoUrl,
		CreatedBy:   userId,
		UpdatedBy:   userId,
	}

	if err := s.repository.CreateGroup(group); err != nil {
		return errors.ErrCreateGroup(err.Error())
	}

	//asignar al userId permisos de admin sobre el grupo
	// if err := s.

	return nil
}

func (s *Service) GetUserGroups(userId []uint8) ([]*models.Group, error) {

	g, err := s.repository.GetUserGroups(userId)
	if err != nil {
		return nil, err
	}
	return g, nil

}

func (s *Service) GetGroupById(groupId []uint8) (*models.Group, error) {

	log.Printf("borrar2 %v", groupId)
	g, err := s.repository.GetGroupById(groupId)
	if err != nil {
		return nil, err
	}
	return g, nil

}
