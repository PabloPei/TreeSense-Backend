package permission

import (
	"database/sql"
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


/// Permissions ///

func (s *SQLRepository) GetUserPermissions(userId []uint8) ([]PermissionAssignment, error) {

    query := `
    SELECT DISTINCT ON (p.permission_name)
        p.permission_name,
        p.description,
        r.role_name
    FROM 
        auth.user_role ur
    JOIN 
        auth.role r ON ur.role_id = r.role_id
    JOIN 
        auth.role_permission rp ON r.role_name = rp.role_name
    JOIN 
        auth.permission p ON rp.permission_name = p.permission_name
    WHERE 
        ur.user_id = $1
    ORDER BY 
        p.permission_name, r.role_id;  -- Agregar orden para DISTINCT ON
    `

	rows, err := s.db.Query(query, userId)
	if err != nil {
		return nil, errors.ErrReadingPermission(err.Error())
	}
	defer rows.Close()

	var permissions []PermissionAssignment

	for rows.Next() {
		permission, err := scanRowIntoPermissionsAssigment(rows)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, *permission)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ErrPermissionScan(err.Error())
	}

	return permissions, nil
}

func (s *SQLRepository) GetPermissionByName(name string) (*PermissionAssignment, error) {
	query := `
		SELECT 
			p.permission_name,
			p.description,
			NULL AS role_id,
			NULL AS role_name
		FROM 
			auth.permission p
		WHERE 
			p.permission_name = $1
	`

	row := s.db.QueryRow(query, name)

	perm, err := scanRowIntoPermissionsAssigment(row)
	if err != nil {
		return nil, err
	}

	return perm, nil
}


/// Aux Function ///
func scanRowIntoPermissionsAssigment(row scannable) (*PermissionAssignment, error) {
	permission := new(PermissionAssignment)
	err := row.Scan(
		&permission.PermissionName,
		&permission.Description,
		&permission.RoleName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrPermissionNotFound
		}
		return nil, errors.ErrPermissionScan(err.Error())
	}
	return permission, nil
}


