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

	_, err = s.repository.GetSpecieById(payload.Specie)
	if err != nil {
		return errors.ErrTreeSpecieNotFound
	}
	
	location := fmt.Sprintf("POINT(%f %f)", payload.Longitude, payload.Latitude)

	tree := Tree{
		//RouteId:     payload.RouteId,
		Specie:      payload.Specie,
		State:       payload.State,
		Location:    location,
		Antique:     payload.Antique,
		Height:      payload.Height,
		Diameter:    payload.Diameter,
		PhotoUrl:    payload.PhotoUrl,
		Description: payload.Description,
		CreatedBy:   userId,
	}


	return s.repository.CreateTree(tree)
}

func (s *Service) GetSpecies() ([]TreeSpecie, error) {

	return s.repository.GetSpecies()

}