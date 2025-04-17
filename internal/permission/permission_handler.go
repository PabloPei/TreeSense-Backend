package permission

import (
	"net/http"

	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/internal/middlewares"
	"github.com/PabloPei/TreeSense-Backend/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	service PermissionService
}

func NewHandler(service PermissionService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router, middleware *middlewares.Middleware) {

	router.HandleFunc("", middleware.RequireAuthAndPermission([]string{"MANAGE"}, false)(h.handleGetPermissions)).Methods("GET")
	router.HandleFunc("/{email}", middleware.RequireAuthAndPermission([]string{"MANAGE"}, false)(h.handleGetUserPermissions)).Methods("GET")
}


func (h *Handler) handleGetPermissions(w http.ResponseWriter, r *http.Request) {

	userId, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrJWTInvalidToken)
		return
	}

	permissions, err := h.service.GetCurrentUserPermissions(userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, permissions)
}

func (h *Handler) handleGetUserPermissions(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email, ok := vars["email"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrUserNotFound)
		return
	}

	permissions, err := h.service.GetUserPermissions(email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, permissions)
}
