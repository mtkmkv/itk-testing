package feature_wallet_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mtkmkv/itk-testing/initial/core/domain"
	repository "github.com/mtkmkv/itk-testing/initial/features/wallet/repository"
)

type mockRepository struct {
	getByIDFunc          func(context.Context, uuid.UUID) (domain.Wallet, error)
	applyOperationFunc   func(context.Context, uuid.UUID, domain.OperationType, *int64) error
}

func (m *mockRepository) GetByID(
	ctx context.Context,
	walletID uuid.UUID,
) (domain.Wallet, error) {
	return m.getByIDFunc(ctx, walletID)
}

func (m *mockRepository) ApplyOperation(
	ctx context.Context,
	walletID uuid.UUID,
	operationType domain.OperationType,
	amount *int64,
) error {
	return m.applyOperationFunc(ctx, walletID, operationType, amount)
}

func TestService_GetByID(t *testing.T) {
	walletID := uuid.New()
	createdAt := time.Now()

	expectedWallet := domain.Wallet{
		ID:        walletID,
		Balance:   1000,
		CreatedAt: createdAt,
	}

	repo := &mockRepository{
		getByIDFunc: func(
			ctx context.Context,
			id uuid.UUID,
		) (domain.Wallet, error) {
			if id != walletID {
				t.Fatalf("expected wallet ID %s, got %s", walletID, id)
			}

			return expectedWallet, nil
		},
	}

	service := NewService(repo)

	wallet, err := service.GetByID(context.Background(), walletID)
	if err != nil {
		t.Fatalf("get wallet: %v", err)
	}

	if wallet != expectedWallet {
		t.Fatalf("expected wallet %+v, got %+v", expectedWallet, wallet)
	}
}

func TestService_GetByID_InvalidWalletID(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(
			ctx context.Context,
			id uuid.UUID,
		) (domain.Wallet, error) {
			t.Fatal("repository should not be called")
			return domain.Wallet{}, nil
		},
	}

	service := NewService(repo)

	_, err := service.GetByID(
		context.Background(),
		uuid.Nil,
	)

	if !errors.Is(err, ErrInvalidWalletID) {
		t.Fatalf(
			"expected ErrInvalidWalletID, got %v",
			err,
		)
	}
}

func TestService_GetByID_RepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")
	walletID := uuid.New()

	repo := &mockRepository{
		getByIDFunc: func(
			ctx context.Context,
			id uuid.UUID,
		) (domain.Wallet, error) {
			return domain.Wallet{}, expectedErr
		},
	}

	service := NewService(repo)

	_, err := service.GetByID(
		context.Background(),
		walletID,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}

func TestService_ApplyOperation(t *testing.T) {
	walletID := uuid.New()
	amount := int64(500)

	repo := &mockRepository{
		applyOperationFunc: func(
			ctx context.Context,
			id uuid.UUID,
			operationType domain.OperationType,
			gotAmount *int64,
		) error {
			if id != walletID {
				t.Fatalf("expected wallet ID %s, got %s", walletID, id)
			}

			if operationType != domain.OperationTypeDeposit {
				t.Fatalf(
					"expected operation type %s, got %s",
					domain.OperationTypeDeposit,
					operationType,
				)
			}

			if gotAmount == nil || *gotAmount != amount {
				t.Fatalf(
					"expected amount %d, got %v",
					amount,
					gotAmount,
				)
			}

			return nil
		},
	}

	service := NewService(repo)

	err := service.ApplyOperation(
		context.Background(),
		walletID,
		domain.OperationTypeDeposit,
		&amount,
	)

	if err != nil {
		t.Fatalf("apply operation: %v", err)
	}
}

func TestService_ApplyOperation_InvalidWalletID(t *testing.T) {
	amount := int64(100)

	repo := &mockRepository{
		applyOperationFunc: func(
			ctx context.Context,
			id uuid.UUID,
			operationType domain.OperationType,
			amount *int64,
		) error {
			t.Fatal("repository should not be called")
			return nil
		},
	}

	service := NewService(repo)

	err := service.ApplyOperation(
		context.Background(),
		uuid.Nil,
		domain.OperationTypeDeposit,
		&amount,
	)

	if !errors.Is(err, ErrInvalidWalletID) {
		t.Fatalf(
			"expected ErrInvalidWalletID, got %v",
			err,
		)
	}
}

func TestService_ApplyOperation_InvalidOperation(t *testing.T) {
	walletID := uuid.New()
	amount := int64(100)

	repo := &mockRepository{
		applyOperationFunc: func(
			ctx context.Context,
			id uuid.UUID,
			operationType domain.OperationType,
			amount *int64,
		) error {
			t.Fatal("repository should not be called")
			return nil
		},
	}

	service := NewService(repo)

	err := service.ApplyOperation(
		context.Background(),
		walletID,
		domain.OperationType("INVALID"),
		&amount,
	)

	if !errors.Is(err, repository.ErrInvalidOperation) {
		t.Fatalf(
			"expected ErrInvalidOperation, got %v",
			err,
		)
	}
}

func TestService_ApplyOperation_InvalidAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount *int64
	}{
		{
			name:   "nil amount",
			amount: nil,
		},
		{
			name: "zero amount",
			amount: func() *int64 {
				value := int64(0)
				return &value
			}(),
		},
		{
			name: "negative amount",
			amount: func() *int64 {
				value := int64(-100)
				return &value
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				applyOperationFunc: func(
					ctx context.Context,
					id uuid.UUID,
					operationType domain.OperationType,
					amount *int64,
				) error {
					t.Fatal("repository should not be called")
					return nil
				},
			}

			service := NewService(repo)

			err := service.ApplyOperation(
				context.Background(),
				uuid.New(),
				domain.OperationTypeDeposit,
				tt.amount,
			)

			if !errors.Is(err, ErrInvalidAmount) {
				t.Fatalf(
					"expected ErrInvalidAmount, got %v",
					err,
				)
			}
		})
	}
}

func TestService_ApplyOperation_RepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")
	walletID := uuid.New()
	amount := int64(100)

	repo := &mockRepository{
		applyOperationFunc: func(
			ctx context.Context,
			id uuid.UUID,
			operationType domain.OperationType,
			gotAmount *int64,
		) error {
			return expectedErr
		},
	}

	service := NewService(repo)

	err := service.ApplyOperation(
		context.Background(),
		walletID,
		domain.OperationTypeDeposit,
		&amount,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}