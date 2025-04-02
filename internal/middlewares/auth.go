package middlewares

type RoleService interface{
	UserHasRole(string,[]uint8)( bool, error)
}

type UserService interface{
	UserExist(userId []uint8) (bool, error)
}

