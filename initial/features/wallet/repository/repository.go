package features_wallet_repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/mtkmkv/itk-testing/initial/core/domain"
)

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
	ErrInvalidOperation    = errors.New("invalid operation type")
)

type WalletRepository interface {
	GetByID(
		ctx context.Context,
		walletID uuid.UUID,
	) (domain.Wallet, error)

	ApplyOperation(
		ctx context.Context,
		walletID uuid.UUID,
		operationType domain.OperationType,
		amount *int64,
	) error
}
