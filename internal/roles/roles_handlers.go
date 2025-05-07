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

	/// Roles ///
	router.HandleFunc("", middleware.RequireAuthAndPermission([]string{"MANAGE"}, false)(h.handleCreateRole)).Methods("POST")
	router.HandleFunc("", middleware.RequireAuthAndPermission([]string{"MANAGE"}, false)(h.handleGetRoles)).Methods("GET")
	router.HandleFunc("/all", middleware.RequireAuthAndPermission([]string{"MANAGE"}, false)(h.handleGetAllRoles)).Methods("GET")
	router.HandleFunc("/{email}", middleware.RequireAuthAndPermission([]string{"MANAGE"}, false)(h.handleGetUserRoles)).Methods("GET")

	/// Assigments /// 
	router.HandleFunc("/{email}", middleware.RequireAuthAndPermission([]string{}, false)(h.handleCreateRoleAssigment)).Methods("POST")
	router.HandleFunc("/{email}", middleware.RequireAuthAndPermission([]string{}, false)(h.handleDeleteRoleAssigment)).Methods("DELETE")
}

/// Roles ///
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

func (h *Handler) handleGetAllRoles(w http.ResponseWriter, r *http.Request) {

	roles, err := h.service.GetRoles()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, roles)
}

func (h *Handler) handleGetRoles(w http.ResponseWriter, r *http.Request) {

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

	utils.WriteJSON(w, http.StatusOK, roles)
}

func (h *Handler) handleGetUserRoles(w http.ResponseWriter, r *http.Request) {

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

	utils.WriteJSON(w, http.StatusOK, roles)
}

/// Assigments /// 

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


func (h *Handler) handleDeleteRoleAssigment(w http.ResponseWriter, r *http.Request) {

	var roleAssigment DeleteUserRoleAssigmentPayload
	if err := utils.ParseJSON(r, &roleAssigment); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(roleAssigment); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload(validationErrors.Error()))
		return
	}

	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrUserNotFound)
		return
	}

	
	err := h.service.DeleteRoleAssigment(roleAssigment, email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Role assigned deleted",
	})
}