package users

import (
	"net/http"

	"github.com/PabloPei/TreeSense-Backend/internal/errors"
	"github.com/PabloPei/TreeSense-Backend/internal/middlewares"
	"github.com/PabloPei/TreeSense-Backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	// User routes
	router.HandleFunc("/user/register", h.handleUserRegister).Methods("POST")
	router.HandleFunc("/user/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/user/refresh-token", middlewares.WithRefreshTokenAuth(h.handleRefreshToken)).Methods("POST")
	router.HandleFunc("/user/photo/{email}", middlewares.WithJWTAuth(h.handleUserPhoto)).Methods("POST", "PUT")
	router.HandleFunc("/user/{email}/role", middlewares.WithJWTAuth(h.handleCreateRoleAssigment)).Methods("POST")
	//logout se aplica desde el frontend

	// Admin Routes
	router.HandleFunc("/user/{email}", middlewares.WithJWTAuth(h.handleGetUser)).Methods("GET")
}

func (h *Handler) handleUserRegister(w http.ResponseWriter, r *http.Request) {

	var user RegisterUserPayload
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

	var user LogInUserPayload

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

// TODO Arreglar esta ruta 
func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {

	userId, err := middlewares.GetUserIDFromContext(r.Context())

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

	var payload UploadPhotoPayload
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