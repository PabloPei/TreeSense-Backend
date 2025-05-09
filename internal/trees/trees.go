package trees

import (
	"time"
)

type Tree struct {
	TreeId      []uint8   `json:"treeId"`
	RouteId     []uint8   `json:"routeId"`
	Specie      string    `json:"specie"`
	State       string    `json:"state"`
	Location    string    `json:"location"`
	Antique     int       `json:"antique"`
	Height      float64   `json:"height"`
	Diameter    float64   `json:"diameter"`
	PhotoUrl    string    `json:"photoUrl"`
	Description string    `json:"description"`
	CreatedBy   []uint8   `json:"createdBy"`
	UpdatedBy   []uint8   `json:"updatedBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TreeSpecie struct {
	TreeSpecieId string `json:"treeSpecieId"`
	Description string `json:"description"`
}

type TreeState struct {	
	TreeStateId []uint8 `json:"treeStateId"`
	Description string `json:"description"`
}

type TreeRepository interface {
	GetTreeById(treeId []uint8) (*Tree, error)
	GetTreeStateById(stateId string) (*TreeState, error)
	GetSpecieById(specieId string) (*TreeSpecie, error)
	CreateTree(tree Tree) error
	GetSpecies() ([]TreeSpecie, error)
	GetTreesByUserId(id []uint8) ([]Tree, error)
}

type TreeService interface {
	CreateTree(tree createTreePayload, userId []uint8) error
	GetSpecies() ([]TreeSpecie, error)
	GetTreesByUser(userId []uint8) ([]Tree, error)
}


type createTreePayload struct {	
	//RouteId []uint8 `json:"routeId" validate:"required"`
	Specie string `json:"specie" validate:"required"`
	State string `json:"state" validate:"required"`	
	Latitude       float64 `json:"latitude" validate:"required"`
    Longitude   float64 `json:"longitude" validate:"required"`
	Antique int `json:"antique" validate:"required"`
	Height float64 `json:"height" validate:"required"`
	Diameter float64 `json:"diameter" validate:"required"`
	PhotoUrl string `json:"photoUrl" validate:"required,uri"`
	Description string `json:"description" validate:"required"`
}