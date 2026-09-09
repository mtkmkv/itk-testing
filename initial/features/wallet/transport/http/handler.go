package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/mtkmkv/itk-testing/initial/features/wallet/repository"
	"github.com/mtkmkv/itk-testing/initial/features/wallet/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Operate(w http.ResponseWriter, r *http.Request) {
	var request operateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	walletID, err := uuid.Parse(request.WalletID)
	if err != nil {
		writeError(w, http.StatusBadRequest, service.ErrInvalidWalletID.Error())
		return
	}

	amount := request.Amount

	err = h.service.ApplyOperation(
		r.Context(),
		walletID,
		request.OperationType,
		&amount,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidWalletID):
			writeError(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, service.ErrInvalidAmount):
			writeError(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, repository.ErrInvalidOperation):
			writeError(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, repository.ErrWalletNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, repository.ErrInsufficientBalance):
			writeError(w, http.StatusBadRequest, err.Error())

		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	walletID, err := uuid.Parse(r.PathValue("walletID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, service.ErrInvalidWalletID.Error())
		return
	}

	wallet, err := h.service.GetByID(r.Context(), walletID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidWalletID):
			writeError(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, repository.ErrWalletNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	writeJSON(w, http.StatusOK, walletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
	})
}