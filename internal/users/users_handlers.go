package users

import (
	"net/http"

	"github.com/PabloPei/SmartSpend-backend/internal/errors"
	"github.com/PabloPei/SmartSpend-backend/internal/middlewares"
	"github.com/PabloPei/SmartSpend-backend/internal/models"
	"github.com/PabloPei/SmartSpend-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	service models.UserService
}

func NewHandler(service models.UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	// User routes
	router.HandleFunc("/user/register", h.handleUserRegister).Methods("POST")
	router.HandleFunc("/user/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/user/refresh-token", middlewares.WithRefreshTokenAuth(h.handleRefreshToken)).Methods("POST")
	router.HandleFunc("/user/photo/{email}", middlewares.WithJWTAuth(h.handleUserPhoto)).Methods("POST", "PUT")
	//logout se aplica desde el frontend

	// Admin Routes
	router.HandleFunc("/user/{email}", middlewares.WithJWTAuth(h.handleGetUser)).Methods("GET")
}

func (h *Handler) handleUserRegister(w http.ResponseWriter, r *http.Request) {

	var user models.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload(validationErrors.Error()))
		return
	}

	err := h.service.RegisterUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

	var user models.LogInUserPayload

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload(validationErrors.Error()))
		return
	}

	token, refreshToken, err := h.service.LogInUser(user)
	if err == errors.ErrInvalidCredentials || err == errors.ErrUserNotFound {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"accessToken": token, "refreshToken": refreshToken})
}

func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(models.UserKey).([]uint8)

	newAccessToken, err := h.service.RefreshToken(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"accessToken": newAccessToken})

}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrUserNotFound)
		return
	}

	userPublic, err := h.service.GetUserPublicByEmail(email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, userPublic)
}

func (h *Handler) handleUserPhoto(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email, ok := vars["email"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload("missing email"))
		return
	}

	var payload models.UploadPhotoPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidaPayload(validationErrors.Error()))
		return
	}

	err := h.service.UploadPhoto(payload, email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Photo uploaded successfully",
	})
}
