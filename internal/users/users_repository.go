package users

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

func (s *SQLRepository) CreateUser(user User) error {
	_, err := s.db.Exec(
		"INSERT INTO auth.\"user\" (user_name, email, password) VALUES ($1, $2, $3)",
		user.UserName, user.Email, user.Password,
	)
	if err != nil {
		return errors.ErrCantUploadUser(err.Error())
	}

	return nil
}

func (s *SQLRepository) UploadPhoto(photoUrl string, email string) error {

	_, err := s.db.Exec(
		"UPDATE auth.\"user\" SET photo_url = $1 WHERE email = $2",
		photoUrl, email,
	)

	if err != nil {
		return errors.ErrCantUploadUser(err.Error())
	}

	return nil
}

func (s *SQLRepository) GetUserByEmail(email string) (*User, error) {
	row := s.db.QueryRow("SELECT * FROM auth.\"user\" WHERE email = $1", email)
	return scanRowIntoUser(row)
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
