package auth_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/services"
	"github.com/Romasmi/e-learning-arhitecture/auth-service/internal/utils/http_utils"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}

	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	token, err := h.authService.Login(r.Context(), req.Login, req.Password, ip)
	if err != nil {
		http_utils.JsonUnauthorized(w)
		return
	}

	http_utils.SuccessJsonResponse(w, loginResponse{Token: token})
}

func (h *AuthHandler) ValidateHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http_utils.JsonUnauthorized(w)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http_utils.JsonUnauthorized(w)
		return
	}

	claims, err := h.authService.Validate(r.Context(), parts[1])
	if err != nil {
		http_utils.JsonUnauthorized(w)
		return
	}

	w.Header().Set("X-Auth-User-ID", claims.UserID)
	w.Header().Set("X-Auth-User-Login", claims.Login)
	w.WriteHeader(http.StatusOK)
}

type registerRequest struct {
	UserID   string `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http_utils.JsonError(w, http.StatusBadRequest, err)
		return
	}

	err = h.authService.Register(r.Context(), userID, req.Login, req.Password)
	if err != nil {
		http_utils.JsonError(w, http.StatusInternalServerError, err)
		return
	}

	http_utils.SuccessJsonResponse(w, map[string]string{"status": "registered"})
}
