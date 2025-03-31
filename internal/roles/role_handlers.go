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

func (h *Handler) RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/role", middlewares.WithJWTAuth(h.handleCreateRole)).Methods("POST")
	router.HandleFunc("/role", middlewares.WithJWTAuth(h.handleGetRole)).Methods("GET")

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

func (h *Handler) handleGetRole(w http.ResponseWriter, r *http.Request) {

	roles, err := h.service.GetRoles()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, roles)
}
