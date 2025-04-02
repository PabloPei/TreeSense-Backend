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

type RoleAssigment struct {
	RoleId       	 []uint8   `json:"roleId"`
	RoleName     	 string    `json:"roleName"`
	RoleDescription  string    `json:"roleDescription"`
	ValidUntil    	 time.Time `json:"validUntil"`
	AssignedBy    	 []uint8   `json:"assignedBy"`
}

type RoleRepository interface {
	CreateRole(Role) error
	GetRoles() ([]Role, error) 
	GetRoleByName(roleName string) (*Role, error)
	CreateRoleAssigment(userId []uint8, roleId []uint8, by []uint8, valid_until time.Time) error
	GetUserRoles(userId []uint8)([]RoleAssigment, error)
}

type RoleService interface {
	CreateRole(payload CreateRolePayload) error
	GetRoles() ([]Role, error) 
	CreateRoleAssigment(payload CreateUserRoleAssigmentPayload, email string, by []uint8) error
	GetUserRoles(email string) ([]RoleAssigment, error)
	GetCurrentUserRoles(userId []uint8)([]RoleAssigment, error) 
	UserHasRole(roleName string, userId []uint8)(bool, error)
}

type CreateRolePayload struct {
	RoleName           string `json:"roleName" validate:"required"`
	RoleDescription    string `json:"roleDescription" validate:"required"`
}

type CreateUserRoleAssigmentPayload struct {
	RoleName           string `json:"roleName" validate:"required"`
	ValidUntil         time.Time `json:"validUntil" validate:"required"`
}