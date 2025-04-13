package roles

import (
	"net/http"

	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/internal/middlewares"
	"github.com/PabloPei/TreeSense-Backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	service RoleService
}

func NewHandler(service RoleService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router, middleware *middlewares.Middleware) {

	router.HandleFunc("/role", middleware.RequireAuthAndPermission("ADMIN", false)(h.handleCreateRole)).Methods("POST")
	router.HandleFunc("/role", middleware.RequireAuthAndPermission("", false)(h.handleGetRole)).Methods("GET")
	router.HandleFunc("/role/all", middleware.RequireAuthAndPermission("ADMIN", false)(h.handleGetAllRole)).Methods("GET")
	router.HandleFunc("/role/{email}/assign", middleware.RequireAuthAndPermission("", false)(h.handleCreateRoleAssigment)).Methods("POST")
	router.HandleFunc("/role/{email}", middleware.RequireAuthAndPermission("ADMIN", false)(h.handleGetUserRole)).Methods("GET")

}

func (h *Handler) handleCreateRole(w http.ResponseWriter, r *http.Request) {

	var role CreateRolePayload
	if err := utils.ParseJSON(r, &role); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(role); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload(validationErrors.Error()))
		return
	}

	err := h.service.CreateRole(role)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Role created successfully",
	})
}

func (h *Handler) handleGetAllRole(w http.ResponseWriter, r *http.Request) {

	roles, err := h.service.GetRoles()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, roles)
}

func (h *Handler) handleGetRole(w http.ResponseWriter, r *http.Request) {

	userId, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrJWTInvalidToken)
		return
	}

	roles, err := h.service.GetCurrentUserRoles(userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, roles)
}

func (h *Handler) handleGetUserRole(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email, ok := vars["email"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrUserNotFound)
		return
	}

	roles, err := h.service.GetUserRoles(email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, roles)
}


//TODO verificar que no exista el rol assigment ya, para eso crear roles by user etc
func (h *Handler) handleCreateRoleAssigment(w http.ResponseWriter, r *http.Request) {

	var roleAssigment CreateUserRoleAssigmentPayload
	if err := utils.ParseJSON(r, &roleAssigment); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(roleAssigment); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload(validationErrors.Error()))
		return
	}

	userId, err := middlewares.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrJWTInvalidToken)
		return
	}

	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrUserNotFound)
		return
	}

	
	err = h.service.CreateRoleAssigment(roleAssigment, email, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Role successfully assigned",
	})
}