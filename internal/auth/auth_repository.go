package auth

//TODO agregar interfaces de auth y middleware de permisos

import (
	"database/sql"
)

// Postgres SQL Repository
type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

/* Estafuncion va a ser la que cree una entrada en user role, donde se asigna un rol al usuario sobre un grupo
func (s *SQLRepository) CreateRoleAssigment(user models.User) error {
	_, err := s.db.Exec(
		"INSERT INTO auth.\"user\" (user_name, email, password) VALUES ($1, $2, $3)",
		user.UserName, user.Email, user.Password,
	)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	return nil
}
*/

