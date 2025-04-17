package audit

import (
    "database/sql"
    "fmt"
)

type SQLRepository struct {
	db *sql.DB
}


func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) LogActivity(userID []uint8, action string) error {
	_, err := r.db.Exec(`
		INSERT INTO audit."activity_log" (user_id, action_name)
		VALUES ($1, $2)
	`, userID, action)
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}
	return nil
}

