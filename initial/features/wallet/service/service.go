package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/mtkmkv/itk-testing/initial/core/domain"
	"github.com/mtkmkv/itk-testing/initial/features/wallet/repository"
)

var (
	ErrInvalidWalletID = errors.New("invalid wallet id")
	ErrInvalidAmount   = errors.New("amount must be greater than zero")
)

type Service struct {
	repository repository.WalletRepository
}

func NewService(repository repository.WalletRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByID(
	ctx context.Context,
	walletID uuid.UUID,
) (domain.Wallet, error) {
	if walletID == uuid.Nil {
		return domain.Wallet{}, ErrInvalidWalletID
	}

	return s.repository.GetByID(ctx, walletID)
}

func (s *Service) ApplyOperation(
	ctx context.Context,
	walletID uuid.UUID,
	operationType domain.OperationType,
	amount *int64,
) error {
	if walletID == uuid.Nil {
		return ErrInvalidWalletID
	}

	if !operationType.IsValid() {
		return repository.ErrInvalidOperation
	}

	if amount == nil || *amount <= 0 {
		return ErrInvalidAmount
	}

	return s.repository.ApplyOperation(
		ctx,
		walletID,
		operationType,
		amount,
	)
}
