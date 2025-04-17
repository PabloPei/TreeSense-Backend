package permission

type PermissionAssignment struct {
	RoleName       string   `json:"roleName"`
	PermissionName string   `json:"permissionName"`
	Description    string   `json:"description"`
}

type PermissionRepository interface {
	GetUserPermissions(userId []uint8) ([]PermissionAssignment, error)
	GetPermissionByName(name string) (*PermissionAssignment, error)
}

type PermissionService interface {
	GetUserPermissions(email string) ([]PermissionAssignment, error)
	GetCurrentUserPermissions(userId []uint8) ([]PermissionAssignment, error)
}
