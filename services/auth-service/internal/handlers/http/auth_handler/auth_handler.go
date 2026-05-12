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
	Email    string `json:"email"`
	Password string `json:"password"`
	PortalID string `json:"portal_id"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

	result, err := h.authService.Login(r.Context(), req.Email, req.Password, req.PortalID, ip)
	if err != nil {
		http_utils.JsonUnauthorized(w)
		return
	}

	http_utils.SuccessJsonResponse(w, loginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
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
	PortalID string `json:"portal_id"`
	Role     string `json:"role"`
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.ErrorInvalidRequestBody(w, err)
		return
	}

	var userID uuid.UUID
	var err error
	if req.UserID != "" {
		userID, err = uuid.Parse(req.UserID)
		if err != nil {
			http_utils.JsonError(w, http.StatusBadRequest, err)
			return
		}
	}

	resID, err := h.authService.Register(r.Context(), userID, req.Login, req.Password, req.PortalID, req.Role)
	if err != nil {
		http_utils.JsonError(w, http.StatusInternalServerError, err)
		return
	}

	http_utils.SuccessJsonResponse(w, map[string]string{
		"status":  "registered",
		"user_id": resID.String(),
	})
}
