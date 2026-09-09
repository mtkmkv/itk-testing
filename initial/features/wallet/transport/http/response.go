package feature_wallet_transport_http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type walletResponse struct {
	WalletID uuid.UUID `json:"walletId"`
	Balance  int64    `json:"balance"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{
		Error: message,
	})
}