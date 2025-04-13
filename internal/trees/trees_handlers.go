package trees

import (
	"net/http"
	"log"

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

	// User routes
	router.HandleFunc("/tree", middleware.RequireAuthAndPermission("FIELD AGENT", false)(h.handleCreateTree)).Methods("POST")

}

func (h *Handler) handleCreateTree(w http.ResponseWriter, r *http.Request) {

	userId, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrJWTInvalidToken)
		return
	}

	log.Println("userId %v", userId)

	var tree createTreePayload
	if err := utils.ParseJSON(r, &tree); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	log.Println("tree %v", tree)
	
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
