package roles

import (
	"time"
)

type Role struct {
	RoleId       	 []uint8   `json:"roleId"`
	RoleName     	 string    `json:"roleName"`
	RoleDescription  string    `json:"roleDescription"`
	CreatedAt    	 time.Time `json:"createdAt"`
	UpdatedAt    	 time.Time `json:"updatedAt"`
}

type RoleRepository interface {
	CreateRole(Role) error
	GetRoles() ([]Role, error) 
	GetRoleByName(roleName string) (*Role, error)
}

type RoleService interface {
	CreateRole(payload CreateRolePayload) error
	GetRoles() ([]Role, error) 
}

type CreateRolePayload struct {
	RoleName           string `json:"roleName" validate:"required"`
	RoleDescription    string `json:"roleDescription" validate:"required"`
}
