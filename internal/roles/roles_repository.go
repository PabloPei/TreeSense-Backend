package roles

import (
	"database/sql"
    "log"

	"github.com/PabloPei/TreeSense-Backend/internal/errors"
)

// Postgres SQL Repository
type SQLRepository struct {
	db *sql.DB
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

	roles, err := scanRowsIntoRoles(rows)
	if err != nil {
		return nil, errors.ErrRoleScan(err.Error())
	}

	return roles, nil
}

func (s *SQLRepository) GetRoleByName(roleName string) (*Role, error) {

	log.Println("roleName %s", roleName )
	row := s.db.QueryRow("SELECT role_id, role_name, description, created_at, updated_at FROM auth.\"role\" WHERE role_name = $1", roleName)
	

	role, err := scanRowIntoRole(row)
	if err != nil {
		return nil, errors.ErrRoleScan(err.Error())
	}

	return role, nil
}

//TODO ver como no repetir codigo aca (interfaz?)
func scanRowsIntoRoles(rows *sql.Rows) ([]Role, error) {

	var roles []Role

	for rows.Next() {
		role := new(Role)
		err := rows.Scan(
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
		roles = append(roles, *role)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrRoleScan(err.Error())
	}

	return roles, nil
}
func scanRowIntoRole(row *sql.Row) (*Role, error) {

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


/*
func (s *SQLRepository) CreateRoleAssigment(user User) error {
	_, err := s.db.Exec(
		"INSERT INTO auth.\"role\" (user_name, email, password) VALUES ($1, $2, $3)",
		user.UserName, user.Email, user.Password,
	)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	return nil
}


func (s *SQLRepository) GetUserById(id []uint8) (*User, error) {
	row := s.db.QueryRow("SELECT * FROM auth.\"user\" WHERE user_id = $1", id)
	return scanRowIntoUser(row)
}

func scanRowIntoUser(row *sql.Row) (*User, error) {

	user := new(User)
	err := row.Scan(
		&user.UserId,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.PhotoUrl,
		&user.LanguageCode,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrUserScan(err.Error())
	}
	return user, nil
}
*/