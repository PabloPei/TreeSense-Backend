package models

import (
	"time"
)

type Group struct {
	GroupId     []uint8   `json:"groupId"`
	GroupName   string    `json:"groupName"`
	Description string    `json:"description"`
	PhotoUrl    string    `json:"photoUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedBy   []uint8   `json:"createdBy"`
	UpdatedBy   []uint8   `json:"updatedBy"`
}

type GroupRepository interface {
	CreateGroup(Group) error
	GetGroupById(groupId []uint8) (*Group, error)
	GetGroupByName(name string) (*Group, error)
	GetUserGroupByName(user []uint8, name string) (*Group, error)
	UploadPhoto(photoUrl string, groupId []uint8) error
	GetUserGroups(user []uint8) ([]*Group, error)
}

type GroupService interface {
	CreateGroup(payload CreateGroupPayload, userId []uint8) error
	GetGroupById(groupId []uint8) (*Group, error)
	GetUserGroups(userId []uint8) ([]*Group, error)
}

type CreateGroupPayload struct {
	GroupName   string `json:"groupName" validate:"required"`
	Description string `json:"description" validate:"required"`
	PhotoUrl    string `json:"photoUrl" validate:"omitempty,uri"`
}
