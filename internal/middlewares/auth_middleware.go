package middlewares

import (
	"context"
	"net/http"

	"github.com/PabloPei/TreeSense-Backend/internal/audit"
	"github.com/PabloPei/TreeSense-Backend/internal/auth"
	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/utils"
)

type ContextKey string

var UserKey ContextKey = "userId"

type Middleware struct {
	permissionService PermissionService
	userService       UserService
	auditService      audit.AuditService
}

func NewAuthMiddleware(permissionService PermissionService, userService UserService, auditService audit.AuditService) *Middleware {
	return &Middleware{permissionService: permissionService, userService: userService, auditService: auditService}
}

func (m *Middleware) RequireAuthAndPermission(permissions []string, useRefreshToken bool) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// TODO: chequear los cambios para prod. (cambio = este if para que maneje OPTIONS desde el navegador)
			if r.Method == http.MethodOptions {
				handler(w, r)
				return
			}

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

			hasPerm, err := m.permissionService.UserHasPermissions(permissions, userID)
			if err != nil || !hasPerm {
				utils.WriteError(w, http.StatusForbidden, errors.ErrUserNotHavePermissions(permissions))
				return
			}

			// Agregamos userID al contexto
			ctx := context.WithValue(r.Context(), UserKey, userIDStr)
			handler(w, r.WithContext(ctx))

			// Registro de actividad si corresponde
			if auditable, ok := ctx.Value("audit").(bool); ok && auditable {
				action := inferActionName(r.Method, r.URL.Path)
				logActivity := audit.ActivityLog{
					UserID: userID,
					Action: action,
				}

				if err := m.auditService.LogActivity(logActivity); err != nil {
					errors.ErrLogActivity(err)
				}
			}
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
