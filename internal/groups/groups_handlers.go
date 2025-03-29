package groups

import (
	"net/http"

	"github.com/PabloPei/SmartSpend-backend/internal/auth"
	"github.com/PabloPei/SmartSpend-backend/internal/errors"
	"github.com/PabloPei/SmartSpend-backend/internal/middlewares"
	"github.com/PabloPei/SmartSpend-backend/internal/models"
	"github.com/PabloPei/SmartSpend-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	service models.GroupService
}

func NewHandler(service models.GroupService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	// User routes
	router.HandleFunc("/group/create", middlewares.WithJWTAuth(h.handleGroupCreate)).Methods("POST")
	router.HandleFunc("/group/all", middlewares.WithJWTAuth(h.handleGetGroups)).Methods("GET")

	// Admin routes
	router.HandleFunc("/group/{groupId}", middlewares.WithJWTAuth(h.handleGetGroup)).Methods("GET")
}

func (h *Handler) handleGroupCreate(w http.ResponseWriter, r *http.Request) {

	userId, err := auth.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrJWTInvalidToken)
		return
	}

	var group models.CreateGroupPayload
	if err := utils.ParseJSON(r, &group); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(group); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload(validationErrors.Error()))
		return
	}

	err = h.service.CreateGroup(group, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Group created successfully",
	})
}
func (h *Handler) handleGetGroups(w http.ResponseWriter, r *http.Request) {

	userId, err := auth.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.ErrJWTInvalidToken)
		return
	}

	userPublic, err := h.service.GetUserGroups(userId)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, userPublic)
}

func (h *Handler) handleGetGroup(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	groupId := []uint8(vars["groupId"])
	if groupId == nil {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrUserNotFound)
		return
	}

	userPublic, err := h.service.GetGroupById(groupId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, userPublic)
}
