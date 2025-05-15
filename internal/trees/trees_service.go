package trees

// TODO: El tree service tiene que trer el rout service para verificar que existan y que corresponda al usuario
import (
	"fmt"

	"github.com/PabloPei/TreeSense-Backend/internal/errors"
)

type Service struct {
	repository TreeRepository
}

func NewService(repository TreeRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateTree(payload createTreePayload, userId []uint8) error {

	//TODO validar ruta
	_, err := s.repository.GetTreeStateById(payload.State)
	if err != nil {
		return errors.ErrTreeStateNotFound
	}

	_, err = s.repository.GetSpeciesById(payload.Species)
	if err != nil {
		return errors.ErrTreeSpeciesNotFound
	}

	location := fmt.Sprintf("POINT(%f %f)", payload.Longitude, payload.Latitude)

	tree := Tree{
		//RouteId:     payload.RouteId,
		Species:     payload.Species,
		State:       payload.State,
		Location:    location,
		Age:         payload.Age,
		Height:      payload.Height,
		Diameter:    payload.Diameter,
		PhotoUrl:    payload.PhotoUrl,
		Description: payload.Description,
		CreatedBy:   userId,
	}

	return s.repository.CreateTree(tree)
}

func (s *Service) GetSpecies() ([]TreeSpecies, error) {

	return s.repository.GetSpecies()

}

func (s *Service) GetTreesByUser(userId []uint8) ([]Tree, error) {

	return s.repository.GetTreesByUserId(userId)

}
