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

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

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

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userId)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

// TODO agregar validacion del access token que este vencido pero pertenezca al usuario
func WithRefreshTokenAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

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

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userId)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func GetUserIDFromContext(ctx context.Context) ([]uint8, error) {

	userId, ok := ctx.Value(UserKey).(string)

	userIdUint := []uint8(userId)

	if !ok {
		return nil, errors.ErrJWTInvalidToken
	}

	return userIdUint, nil
}


/* TODO middle ware para verificar los permisos
func RequirePermission(handlerFunc http.HandlerFunc, permission string) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		userID := 1   // Simulación (debe venir del contexto o JWT)
		groupID := 10 // Simulación (debe venir de la ruta o sesión)

		// Verificar si el usuario tiene el permiso necesario
		hasPermission, err := userHasSpecificPermission(db, userID, groupID, permission)
		if err != nil {
			http.Error(w, "Error verificando permisos", http.StatusInternalServerError)
			return
		}

		if !hasPermission {
			http.Error(w, "No tienes permisos para esta acción", http.StatusForbidden)
			return
		}

		// Si tiene permisos, ejecutar el handler
		next(w, r)
	}
}
*/
