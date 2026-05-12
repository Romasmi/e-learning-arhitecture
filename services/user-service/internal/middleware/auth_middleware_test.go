package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		headerUserID   string
		pathUserID     string
		expectedStatus int
	}{
		{
			name:           "valid user id",
			headerUserID:   "user-1",
			pathUserID:     "user-1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing header",
			headerUserID:   "",
			pathUserID:     "user-1",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "mismatched user id",
			headerUserID:   "user-1",
			pathUserID:     "user-2",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "no path user id (exempt)",
			headerUserID:   "user-1",
			pathUserID:     "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req, _ := http.NewRequest(http.MethodGet, "/user/"+tt.pathUserID, nil)
			if tt.headerUserID != "" {
				req.Header.Set("X-Auth-User-ID", tt.headerUserID)
			}

			if tt.pathUserID != "" {
				req = mux.SetURLVars(req, map[string]string{"userId": tt.pathUserID})
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}
