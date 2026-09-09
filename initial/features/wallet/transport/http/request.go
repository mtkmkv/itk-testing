package http

import "github.com/mtkmkv/itk-testing/initial/core/domain"

type operateRequest struct {
	WalletID      string               `json:"walletId"`
	OperationType domain.OperationType `json:"operationType"`
	Amount        int64                `json:"amount"`
}
