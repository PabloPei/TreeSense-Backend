package groups

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/PabloPei/SmartSpend-backend/internal/errors"
	"github.com/PabloPei/SmartSpend-backend/internal/models"
)

// Postgres SQL Repository
type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (s *SQLRepository) CreateGroup(group models.Group) error {

	_, err := s.db.Exec(
		"INSERT INTO public.\"group\" (group_name, description, created_by, updated_by) VALUES ($1, $2, $3, $4)",
		group.GroupName, group.Description, group.CreatedBy, group.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("error al crear el grupo: %w", err)
	}

	return nil
}

func (s *SQLRepository) GetGroupById(groupId []uint8) (*models.Group, error) {

	groupIdstr := string(groupId)
	row := s.db.QueryRow("SELECT * FROM public.\"group\" WHERE group_id = $1", groupIdstr)
	log.Printf("borrar %v", row)
	return scanRowIntoUser(row)

}

func (s *SQLRepository) GetGroupByName(name string) (*models.Group, error) {

	row := s.db.QueryRow("SELECT * FROM public.\"group\" WHERE group_name = $1", name)
	return scanRowIntoUser(row)

}

func (s *SQLRepository) GetUserGroupByName(user []uint8, name string) (*models.Group, error) {

	row := s.db.QueryRow("SELECT g.* FROM public.\"group\" g INNER JOIN auth.\"user_role\" ur ON g.group_id = ur.group_id WHERE g.group_name = $1 AND ur.user_id = $2", name, user)
	return scanRowIntoUser(row)

}

func (s *SQLRepository) GetUserGroups(user []uint8) ([]*models.Group, error) {

	rows, err := s.db.Query(`
        SELECT g.* 
        FROM public."group" g
        INNER JOIN auth."user_role" ur ON g.group_id = ur.group_id
        WHERE ur.user_id = $1`, user)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los grupos del usuario: %w", err)
	}
	defer rows.Close()

	var groups []*models.Group

	for rows.Next() {
		group := new(models.Group)
		err := rows.Scan(
			&group.GroupId,
			&group.GroupName,
			&group.Description,
			&group.PhotoUrl,
			&group.CreatedAt,
			&group.CreatedBy,
			&group.UpdatedAt,
			&group.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		// Agregar el grupo al slice
		groups = append(groups, group)
	}

	return groups, nil

}

func (s *SQLRepository) UploadPhoto(photoUrl string, groupId []uint8) error {

	_, err := s.db.Exec(
		"UPDATE  public.\"group\" SET photo_url = $1 WHERE group_id = $2",
		photoUrl, groupId,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar la foto: %w", err)
	}

	return nil
}

func scanRowIntoUser(row *sql.Row) (*models.Group, error) {

	group := new(models.Group)
	err := row.Scan(
		&group.GroupId,
		&group.GroupName,
		&group.Description,
		&group.PhotoUrl,
		&group.CreatedAt,
		&group.CreatedBy,
		&group.UpdatedAt,
		&group.UpdatedBy,
	)

	if err != nil {
		log.Printf("borrar error")
		if err == sql.ErrNoRows {
			return nil, errors.ErrGroupNotFound
		}
		log.Printf("error %v", err)
		return nil, errors.ErrUserScan(err.Error())
	}
	return group, nil
}
