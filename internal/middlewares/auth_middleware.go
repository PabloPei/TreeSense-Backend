package middlewares

import (
	"context"
	"net/http"

	"github.com/PabloPei/TreeSense-Backend/internal/auth"
	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/utils"
)

type ContextKey string

var UserKey ContextKey = "userId"

type Middleware struct {
	roleService RoleService
	userService UserService
}

func NewAuthMiddleware(roleService RoleService, userService UserService) *Middleware {
	return &Middleware{roleService: roleService, userService: userService}
}

func (m *Middleware)WithAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := utils.GetTokenFromRequest(r)

		claims, err := auth.ValidateJWT(tokenString, false)
		if err != nil {
			utils.WriteError(w, http.StatusForbidden, errors.ErrJWTInvalidToken)
			return
		}

		userId, ok := claims["userId"].(string)

		if !ok {
			utils.WriteError(w, http.StatusForbidden, errors.ErrJWTInvalidToken)
			return
		}

		userExist, err := m.userService.UserExist([]uint8(userId))

		if !userExist {
			utils.WriteError(w, http.StatusForbidden, errors.ErrUserNotFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userId)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

// TODO Ver si implementar como randomstring almacenado en la base, permite hacer logout
func (m *Middleware)WithRefreshTokenAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := utils.GetTokenFromRequest(r)

		claims, err := auth.ValidateJWT(tokenString, true)
		if err != nil {
			utils.WriteError(w, http.StatusForbidden, errors.ErrJWTInvalidToken)
			return
		}

		userId, ok := claims["userId"].(string)
		if !ok {
			utils.WriteError(w, http.StatusForbidden, errors.ErrJWTInvalidToken)
			return
		}

		userExist, err := m.userService.UserExist([]uint8(userId))

		if !userExist {
			utils.WriteError(w, http.StatusForbidden, errors.ErrUserNotFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userId)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func (m *Middleware)WithPerm(rol string) func(http.HandlerFunc) http.HandlerFunc {

	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			userID, err := GetUserIDFromContext(r.Context())
			if err != nil {
				utils.WriteError(w, http.StatusForbidden, errors.ErrUserNotFound)
				return
			}

			hasPermission, err := m.roleService.UserHasRole(rol, userID)
			if err != nil {
				utils.WriteError(w, http.StatusForbidden, err)
				return
			}

			if !hasPermission {
				utils.WriteError(w, http.StatusForbidden, errors.ErrUserNotHaveRole(rol))
				return
			}

			handlerFunc(w, r)
		}
	}
}

func (s *Middleware) WithAuthAndPerm(permission string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return s.WithAuth(s.WithPerm(permission)(handlerFunc))
}

func GetUserIDFromContext(ctx context.Context) ([]uint8, error) {

	userId, ok := ctx.Value(UserKey).(string)

	userIdUint := []uint8(userId)

	if !ok {
		return nil, errors.ErrJWTInvalidToken
	}

	return userIdUint, nil
}

