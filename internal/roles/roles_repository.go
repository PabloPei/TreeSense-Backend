package roles

import (
	"database/sql"
	"time"

	"github.com/PabloPei/TreeSense-Backend/internal/errors"
)

// Postgres SQL Repository
type SQLRepository struct {
	db *sql.DB
}

type scannable interface {
	Scan(dest ...interface{}) error
}


func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (s *SQLRepository) CreateRole(role Role) error {
	_, err := s.db.Exec(
		"INSERT INTO auth.\"role\" (role_name, description) VALUES ($1, $2)",
		role.RoleName, role.RoleDescription,
	)
	if err != nil {
		return errors.ErrCantUploadRole(err.Error())
	}

	return nil
}


func (s *SQLRepository) GetRoles() ([]Role, error) {

	rows, err := s.db.Query("SELECT role_id, role_name, description, created_at, updated_at FROM auth.\"role\"")
	if err != nil {
		return nil, errors.ErrReadingRole(err.Error())
	}

	defer rows.Close()

	var roles []Role

	for rows.Next() {
		role, err := scanRowIntoRole(rows)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.ErrRoleNotFound
			}
			return nil, errors.ErrRoleScan(err.Error())
		}
		roles = append(roles, *role)

	}

	if err != nil {
		return nil, errors.ErrRoleScan(err.Error())
	}

	return roles, nil
}

func (s *SQLRepository) GetUserRoles(userId []uint8)([]RoleAssigment, error){

	rows, err := s.db.Query("SELECT r.role_id, r.role_name, r.description, ur.valid_until, ur.created_by FROM auth.user_role ur JOIN auth.\"role\" r ON ur.role_id = r.role_id WHERE ur.user_id = $1", userId)

	if err != nil {
		return nil, errors.ErrReadingRole(err.Error())
	}

	defer rows.Close()

	var roles []RoleAssigment

	for rows.Next() {
		role, err := scanRowIntoRoleAssigment(rows)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.ErrRoleNotFound
			}
			return nil, errors.ErrRoleScan(err.Error())
		}
		roles = append(roles, *role)

	}

	if err != nil {
		return nil, errors.ErrRoleScan(err.Error())
	}

	return roles, nil
}

func (s *SQLRepository) GetRoleByName(roleName string) (*Role, error) {

	row := s.db.QueryRow("SELECT role_id, role_name, description, created_at, updated_at FROM auth.\"role\" WHERE role_name = $1", roleName)
	

	role, err := scanRowIntoRole(row)
	if err != nil {
		return nil, errors.ErrRoleScan(err.Error())
	}

	return role, nil
}

func (s *SQLRepository) CreateRoleAssigment(userId []uint8, roleId []uint8, by []uint8, valid_until time.Time) error {
	_, err := s.db.Exec(
		"INSERT INTO auth.\"user_role\" (user_id, role_id, created_by, updated_by, valid_until) VALUES ($1, $2, $3, $3, $4)",
		userId, roleId, by, valid_until,
	)
	if err != nil {
		return errors.ErrCantUploadRole(err.Error())
	}

	return nil
}

func scanRowIntoRole(row scannable) (*Role, error) {

	role := new(Role)
	err := row.Scan(
		&role.RoleId,
		&role.RoleName,
		&role.RoleDescription,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrRoleNotFound
		}
		return nil, errors.ErrRoleScan(err.Error())
	}
	return role, nil
}

func scanRowIntoRoleAssigment(row scannable) (*RoleAssigment, error) {

	roleAssigment := new(RoleAssigment)
	err := row.Scan(
		&roleAssigment.RoleId,
		&roleAssigment.RoleName,
		&roleAssigment.RoleDescription,
		&roleAssigment.ValidUntil,
		&roleAssigment.AssignedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrRoleNotFound
		}
		return nil, errors.ErrRoleScan(err.Error())
	}
	return roleAssigment, nil
}