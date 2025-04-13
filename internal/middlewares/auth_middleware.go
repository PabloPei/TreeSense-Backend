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

func (m *Middleware) RequireAuthAndPermission(permission string, useRefreshToken bool) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := utils.GetTokenFromRequest(r)

			claims, err := auth.ValidateJWT(token, useRefreshToken)
			if err != nil {
				utils.WriteError(w, http.StatusForbidden, errors.ErrJWTInvalidToken)
				return
			}

			userIDStr, ok := claims["userId"].(string)
			if !ok {
				utils.WriteError(w, http.StatusForbidden, errors.ErrJWTInvalidToken)
				return
			}

			userID := []uint8(userIDStr)

			exists, err := m.userService.UserExist(userID)
			if err != nil || !exists {
				utils.WriteError(w, http.StatusForbidden, errors.ErrUserNotFound)
				return
			}

			if permission != "" {
				hasPerm, err := m.roleService.UserHasRole(permission, userID)
				if err != nil || !hasPerm {
					utils.WriteError(w, http.StatusForbidden, errors.ErrUserNotHaveRole(permission))
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserKey, userIDStr)
			handler(w, r.WithContext(ctx))
		}
	}
}

func GetUserIDFromContext(ctx context.Context) ([]uint8, error) {
	userID, ok := ctx.Value(UserKey).(string)
	if !ok {
		return nil, errors.ErrJWTInvalidToken
	}
	return []uint8(userID), nil
}
