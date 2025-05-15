package trees

import (
	"net/http"

	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/internal/middlewares"
	"github.com/PabloPei/TreeSense-Backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	service TreeService
}

func NewHandler(service TreeService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router, middleware *middlewares.Middleware) {

	router.HandleFunc("", middleware.RequireAuthAndPermission([]string{"SURVEY"}, false)(h.handleCreateTree)).Methods("POST")
	router.HandleFunc("", middleware.RequireAuthAndPermission([]string{"SURVEY"}, false)(h.handleGetTreesCurrentUser)).Methods("GET")
	router.HandleFunc("/species", middleware.RequireAuthAndPermission([]string{"SURVEY"}, false)(h.handleGetSpecies)).Methods("GET")

}

func (h *Handler) handleCreateTree(w http.ResponseWriter, r *http.Request) {

	userId, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrJWTInvalidToken)
		return
	}

	var tree createTreePayload
	if err := utils.ParseJSON(r, &tree); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(tree); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload(validationErrors.Error()))
		return
	}

	err = h.service.CreateTree(tree, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Tree created successfully",
	})
}

func (h *Handler) handleGetTreesCurrentUser(w http.ResponseWriter, r *http.Request) {

	userId, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrJWTInvalidToken)
		return
	}

	trees, err := h.service.GetTreesByUser(userId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"trees": trees})
}

func (h *Handler) handleGetSpecies(w http.ResponseWriter, r *http.Request) {

	species, err := h.service.GetSpecies()

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, species)
}
