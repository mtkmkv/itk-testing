package features_wallet_repository_postgres

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mtkmkv/itk-testing/initial/core/domain"
	features_wallet_repository "github.com/mtkmkv/itk-testing/initial/features/wallet/repository"
)

func TestApplyOperationConcurrency(t *testing.T) {
	db, err := sql.Open(
		"postgres",
		"postgres://user:123@localhost:15432/itk_testing_db?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	walletID := uuid.New()

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO wallet.wallets (id, balance) VALUES ($1, 0)`,
		walletID,
	)
	if err != nil {
		t.Fatal(err)
	}

	defer db.ExecContext(
		ctx,
		`DELETE FROM wallet.wallets WHERE id = $1`,
		walletID,
	)

	repository := NewRepository(db)

	const requests = 1000

	var wg sync.WaitGroup
	wg.Add(requests)

	errCh := make(chan error, requests)

	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()

			amount := int64(1)

			if err := repository.ApplyOperation(
				ctx,
				walletID,
				domain.OperationTypeDeposit,
				&amount,
			); err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}

	var balance int64

	err = db.QueryRowContext(
		ctx,
		`SELECT balance FROM wallet.wallets WHERE id = $1`,
		walletID,
	).Scan(&balance)
	if err != nil {
		t.Fatal(err)
	}

	if balance != requests {
		t.Fatalf(
			"expected balance %d, got %d",
			requests,
			balance,
		)
	}

	var operations int

	err = db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM wallet.operations WHERE wallet_id = $1`,
		walletID,
	).Scan(&operations)
	if err != nil {
		t.Fatal(err)
	}

	if operations != requests {
		t.Fatalf(
			"expected operations %d, got %d",
			requests,
			operations,
		)
	}
}

func TestApplyOperationWithdrawConcurrency(t *testing.T) {
	db, err := sql.Open(
		"postgres",
		"postgres://user:123@localhost:15432/itk_testing_db?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	walletID := uuid.New()

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO wallet.wallets (id, balance) VALUES ($1, $2)`,
		walletID,
		1000,
	)
	if err != nil {
		t.Fatal(err)
	}

	defer db.ExecContext(
		ctx,
		`DELETE FROM wallet.wallets WHERE id = $1`,
		walletID,
	)

	repository := NewRepository(db)

	const requests = 1000

	var wg sync.WaitGroup
	wg.Add(requests)

	errCh := make(chan error, requests)

	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()

			amount := int64(1)

			if err := repository.ApplyOperation(
				ctx,
				walletID,
				domain.OperationTypeWithdraw,
				&amount,
			); err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}

	var balance int64

	err = db.QueryRowContext(
		ctx,
		`SELECT balance FROM wallet.wallets WHERE id = $1`,
		walletID,
	).Scan(&balance)
	if err != nil {
		t.Fatal(err)
	}

	if balance != 0 {
		t.Fatalf(
			"expected balance %d, got %d",
			0,
			balance,
		)
	}

	var operations int

	err = db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM wallet.operations WHERE wallet_id = $1`,
		walletID,
	).Scan(&operations)
	if err != nil {
		t.Fatal(err)
	}

	if operations != requests {
		t.Fatalf(
			"expected operations %d, got %d",
			requests,
			operations,
		)
	}
}

func TestApplyOperationWithdrawInsufficientBalance(t *testing.T) {
	db, err := sql.Open(
		"postgres",
		"postgres://user:123@localhost:15432/itk_testing_db?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	walletID := uuid.New()

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO wallet.wallets (id, balance) VALUES ($1, $2)`,
		walletID,
		500,
	)
	if err != nil {
		t.Fatal(err)
	}

	defer db.ExecContext(
		ctx,
		`DELETE FROM wallet.wallets WHERE id = $1`,
		walletID,
	)

	repository := NewRepository(db)

	const requests = 1000

	var wg sync.WaitGroup
	wg.Add(requests)

	errCh := make(chan error, requests)

	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()

			amount := int64(1)

			if err := repository.ApplyOperation(
				ctx,
				walletID,
				domain.OperationTypeWithdraw,
				&amount,
			); err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if errors.Is(err, features_wallet_repository.ErrInsufficientBalance) {
			continue
		}

		t.Error(err)
	}

	var successfulOperations int

	err = db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM wallet.operations WHERE wallet_id = $1`,
		walletID,
	).Scan(&successfulOperations)
	if err != nil {
		t.Fatal(err)
	}

	if successfulOperations != 500 {
		t.Fatalf(
			"expected operations %d, got %d",
			500,
			successfulOperations,
		)
	}

	var balance int64

	err = db.QueryRowContext(
		ctx,
		`SELECT balance FROM wallet.wallets WHERE id = $1`,
		walletID,
	).Scan(&balance)
	if err != nil {
		t.Fatal(err)
	}

	if balance != 0 {
		t.Fatalf(
			"expected balance %d, got %d",
			0,
			balance,
		)
	}
}