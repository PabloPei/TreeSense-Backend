package trees

import (
	"database/sql"

	"github.com/PabloPei/TreeSense-Backend/internal/errors"
)

// Postgres SQL Repository
type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (s *SQLRepository) CreateTree(tree Tree) error {
	_, err := s.db.Exec(
		"INSERT INTO treesense.\"tree\" (specie, state, antique, height, diameter, photo_url, description, location, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, ST_GeomFromText($8, 4326), $9)",
		tree.Specie, tree.State, tree.Antique, tree.Height, tree.Diameter, tree.PhotoUrl, tree.Description, tree.Location, tree.CreatedBy,
	)
	if err != nil {
		return errors.ErrCantUploadTree(err.Error())
	}

	return nil
}

func (s *SQLRepository) GetTreeStateById(stateId string) (*TreeState, error) {
	row :=	s.db.QueryRow("SELECT * FROM treesense.\"tree_state\" where tree_state_id = $1", stateId)
	return scanRowIntoTreeState(row)
}

func (s *SQLRepository) GetSpecieById(stateId string) (*TreeSpecie, error) {
	row :=	s.db.QueryRow("SELECT * FROM treesense.\"tree_specie\" where tree_specie_id = $1", stateId)
	return scanRowIntoTreeSpecie(row)
}

func (s *SQLRepository) GetTreeById(id []uint8) (*Tree, error) {
	row := s.db.QueryRow("SELECT * FROM treesense.\"tree\" WHERE tree_id = $1", id)
	return scanRowIntoTree(row)
}

func scanRowIntoTree(row *sql.Row) (*Tree, error) {

	tree := new(Tree)
	err := row.Scan(
		&tree.TreeId,
		&tree.RouteId,
		&tree.Specie,	
		&tree.State,
		&tree.Location,
		&tree.Antique,
		&tree.Height,
		&tree.Diameter,
		&tree.PhotoUrl,
		&tree.Description,
		&tree.CreatedBy,
		&tree.UpdatedBy,
		&tree.CreatedAt,
		&tree.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrTreeNotFound
		}
		return nil, errors.ErrTreeScan(err.Error())
	}
	return tree, nil
}

func scanRowIntoTreeState(row *sql.Row) (*TreeState, error) {

	treeState := new(TreeState)
	err := row.Scan(
		&treeState.TreeStateId,
		&treeState.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrTreeStateNotFound
		}
		return nil, errors.ErrTreeScan(err.Error())
	}
	return treeState, nil
}

func scanRowIntoTreeSpecie(row *sql.Row) (*TreeSpecie, error) {

	treeSpecie := new(TreeSpecie)
	err := row.Scan(
		&treeSpecie.TreeSpecieId,
		&treeSpecie.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrTreeSpecieNotFound
		}
		return nil, errors.ErrTreeScan(err.Error())
	}
	return treeSpecie, nil
}
