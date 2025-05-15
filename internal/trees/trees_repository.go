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

type scannable interface {
	Scan(dest ...interface{}) error
}

func (s *SQLRepository) CreateTree(tree Tree) error {
	_, err := s.db.Exec(
		"INSERT INTO treesense.\"tree\" (species, state, age, height, diameter, photo_url, description, location, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, ST_GeomFromText($8, 4326), $9)",
		tree.Species, tree.State, tree.Age, tree.Height, tree.Diameter, tree.PhotoUrl, tree.Description, tree.Location, tree.CreatedBy,
	)
	if err != nil {
		return errors.ErrCantUploadTree(err.Error())
	}

	return nil
}

func (s *SQLRepository) GetTreeStateById(stateId string) (*TreeState, error) {
	row := s.db.QueryRow("SELECT * FROM treesense.\"tree_state\" where tree_state_id = $1", stateId)
	return scanRowIntoTreeState(row)
}

func (s *SQLRepository) GetSpeciesById(stateId string) (*TreeSpecies, error) {
	row := s.db.QueryRow("SELECT * FROM treesense.\"tree_species\" where tree_species_id = $1", stateId)
	return scanRowIntoTreeSpecies(row)
}

func (s *SQLRepository) GetSpecies() ([]TreeSpecies, error) {

	rows, err := s.db.Query("SELECT * FROM treesense.\"tree_species\"")
	if err != nil {
		return nil, errors.ErrReadingSpecies(err.Error())
	}
	defer rows.Close()

	var treeSpecies []TreeSpecies

	for rows.Next() {
		species, err := scanRowIntoTreeSpecies(rows)
		if err != nil {
			return nil, err
		}
		treeSpecies = append(treeSpecies, *species)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrTreeSpeciesScan(err.Error())
	}

	return treeSpecies, nil
}

func (s *SQLRepository) GetTreesByUserId(id []uint8) ([]Tree, error) {
	rows, err := s.db.Query("SELECT * FROM treesense.\"tree\" WHERE created_by = $1", id)

	if err != nil {
		return nil, errors.ErrTreeScan(err.Error())
	}

	defer rows.Close()

	var trees []Tree

	for rows.Next() {
		tree, err := scanRowIntoTree(rows)
		if err != nil {
			return nil, err
		}
		trees = append(trees, *tree)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrTreeScan(err.Error())
	}

	return trees, nil
}

func (s *SQLRepository) GetTreeById(id []uint8) (*Tree, error) {
	row := s.db.QueryRow("SELECT * FROM treesense.\"tree\" WHERE tree_id = $1", id)
	return scanRowIntoTree(row)
}

func scanRowIntoTree(row scannable) (*Tree, error) {

	tree := new(Tree)
	err := row.Scan(
		&tree.TreeId,
		&tree.RouteId,
		&tree.Species,
		&tree.State,
		&tree.Location,
		&tree.Age,
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

func scanRowIntoTreeSpecies(row scannable) (*TreeSpecies, error) {

	treeSpecies := new(TreeSpecies)
	err := row.Scan(
		&treeSpecies.TreeSpeciesId,
		&treeSpecies.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrTreeSpeciesNotFound
		}
		return nil, errors.ErrTreeScan(err.Error())
	}
	return treeSpecies, nil
}
