package http_utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JsonResponse(w http.ResponseWriter, statusCode int, output any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		slog.Error("error while encoding response", "error", err)
	}
}

func ErrorInvalidRequestBody(w http.ResponseWriter, err error) {
	JsonError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %s", err.Error()))
}

func SuccessJsonResponse(w http.ResponseWriter, output any) {
	JsonResponse(w, http.StatusOK, output)
}

func JsonError(w http.ResponseWriter, statusCode int, err error) {
	JsonResponse(w, statusCode, &ErrorResponse{Error: err.Error()})
}

func JsonErrorNotFound(w http.ResponseWriter) {
	JsonError(w, http.StatusNotFound, fmt.Errorf("not found"))
}

func JsonInternalServerError(w http.ResponseWriter) {
	JsonError(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
}

func JsonUnauthorized(w http.ResponseWriter) {
	JsonError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
}
