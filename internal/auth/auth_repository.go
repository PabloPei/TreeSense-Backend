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

/* funcion que va a verificar los permisos de un usuuario sobre un grupo
func (s *SQLRepository) userHasPermission(userID []uint8, groupID []uint8, permission string) (bool, error) {

	query := `
	SELECT COUNT(*)
	FROM user u
	JOIN user_role ur ON u.user_id = ur.user_id
	JOIN role_permission rp ON ur.role_id = rp.role_id
	JOIN permission p ON rp.permission_id = p.permission_id
	WHERE u.user_id = $1 AND ur.group_id = $2 AND p.description = $3;`

	var count int
	err := s.db.QueryRow(query, userID, groupID, permission).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
*/
