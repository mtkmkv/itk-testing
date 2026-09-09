package features_wallet_repository_postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/mtkmkv/itk-testing/initial/core/domain"
	repository "github.com/mtkmkv/itk-testing/initial/features/wallet/repository"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(
	ctx context.Context,
	walletID uuid.UUID,
) (domain.Wallet, error) {
	const query = `
		SELECT id, balance, created_at
		FROM wallet.wallets
		WHERE id = $1
	`

	var wallet domain.Wallet

	err := r.db.QueryRowContext(
		ctx,
		query,
		walletID,
	).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Wallet{}, repository.ErrWalletNotFound
		}

		return domain.Wallet{}, fmt.Errorf("get wallet: %w", err)
	}

	return wallet, nil
}

func (r *Repository) ApplyOperation(
	ctx context.Context,
	walletID uuid.UUID,
	operationType domain.OperationType,
	amount *int64,
) error {
	if amount == nil || *amount <= 0 {
		return repository.ErrInvalidAmount
	}

	if !operationType.IsValid() {
		return repository.ErrInvalidOperation
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback()

	var query string

	switch operationType {
	case domain.OperationTypeDeposit:
		query = `
			UPDATE wallet.wallets
			SET balance = balance + $1
			WHERE id = $2
		`

	case domain.OperationTypeWithdraw:
		query = `
			UPDATE wallet.wallets
			SET balance = balance - $1
			WHERE id = $2
			AND balance >= $1
		`
	}

	result, err := tx.ExecContext(
		ctx,
		query,
		*amount,
		walletID,
	)
	if err != nil {
		return fmt.Errorf("update wallet balance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if rows == 0 {
		if operationType == domain.OperationTypeWithdraw {
			var exists bool

			err := tx.QueryRowContext(
				ctx,
				`SELECT EXISTS (
					SELECT 1
					FROM wallet.wallets
					WHERE id = $1
				)`,
				walletID,
			).Scan(&exists)

			if err != nil {
				return fmt.Errorf("check wallet: %w", err)
			}

			if exists {
				return repository.ErrInsufficientBalance
			}
		}

		return repository.ErrWalletNotFound
	}

	const operationQuery = `
		INSERT INTO wallet.operations (
			wallet_id,
			operation_type,
			amount
		)
		VALUES ($1, $2, $3)
	`

	_, err = tx.ExecContext(
		ctx,
		operationQuery,
		walletID,
		operationType,
		*amount,
	)
	if err != nil {
		return fmt.Errorf("create operation: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
