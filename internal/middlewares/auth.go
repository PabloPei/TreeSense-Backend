package middlewares

type PermissionService interface{
	UserHasPermissions(permissionNames []string, userId []uint8) (bool, error)
}

type UserService interface{
	UserExist(userId []uint8) (bool, error)
}

