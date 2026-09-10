package features_wallet_repository_postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/lib/pq"
	"github.com/google/uuid"

	"github.com/mtkmkv/itk-testing/initial/core/domain"
	repository "github.com/mtkmkv/itk-testing/initial/features/wallet/repository"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open(
		"postgres",
		"postgres://user:123@localhost:15432/itk_testing_db?sslmode=disable",
	)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("ping database: %v", err)
	}

	return db
}

func createWallet(t *testing.T, db *sql.DB, balance int64) uuid.UUID {
	t.Helper()

	walletID := uuid.New()

	_, err := db.Exec(
		`INSERT INTO wallet.wallets (id, balance) VALUES ($1, $2)`,
		walletID,
		balance,
	)
	if err != nil {
		t.Fatalf("create wallet: %v", err)
	}

	return walletID
}

func deleteWallet(t *testing.T, db *sql.DB, walletID uuid.UUID) {
	t.Helper()

	_, err := db.Exec(
		`DELETE FROM wallet.wallets WHERE id = $1`,
		walletID,
	)
	if err != nil {
		t.Fatalf("delete wallet: %v", err)
	}
}

func TestRepository_GetByID(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	walletID := createWallet(t, db, 1000)
	defer deleteWallet(t, db, walletID)

	wallet, err := repo.GetByID(context.Background(), walletID)
	if err != nil {
		t.Fatalf("get wallet: %v", err)
	}

	if wallet.ID != walletID {
		t.Fatalf("expected wallet ID %s, got %s", walletID, wallet.ID)
	}

	if wallet.Balance != 1000 {
		t.Fatalf("expected balance 1000, got %d", wallet.Balance)
	}
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	_, err := repo.GetByID(context.Background(), uuid.New())
	if !errors.Is(err, repository.ErrWalletNotFound) {
		t.Fatalf("expected ErrWalletNotFound, got %v", err)
	}
}

func TestRepository_ApplyOperation_Deposit(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	walletID := createWallet(t, db, 1000)
	defer deleteWallet(t, db, walletID)

	amount := int64(500)

	err := repo.ApplyOperation(
		context.Background(),
		walletID,
		domain.OperationTypeDeposit,
		&amount,
	)
	if err != nil {
		t.Fatalf("apply deposit: %v", err)
	}

	wallet, err := repo.GetByID(context.Background(), walletID)
	if err != nil {
		t.Fatalf("get wallet: %v", err)
	}

	if wallet.Balance != 1500 {
		t.Fatalf("expected balance 1500, got %d", wallet.Balance)
	}
}

func TestRepository_ApplyOperation_Withdraw(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	walletID := createWallet(t, db, 1000)
	defer deleteWallet(t, db, walletID)

	amount := int64(300)

	err := repo.ApplyOperation(
		context.Background(),
		walletID,
		domain.OperationTypeWithdraw,
		&amount,
	)
	if err != nil {
		t.Fatalf("apply withdraw: %v", err)
	}

	wallet, err := repo.GetByID(context.Background(), walletID)
	if err != nil {
		t.Fatalf("get wallet: %v", err)
	}

	if wallet.Balance != 700 {
		t.Fatalf("expected balance 700, got %d", wallet.Balance)
	}
}

func TestRepository_ApplyOperation_InsufficientBalance(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	walletID := createWallet(t, db, 1000)
	defer deleteWallet(t, db, walletID)

	amount := int64(1500)

	err := repo.ApplyOperation(
		context.Background(),
		walletID,
		domain.OperationTypeWithdraw,
		&amount,
	)

	if !errors.Is(err, repository.ErrInsufficientBalance) {
		t.Fatalf(
			"expected ErrInsufficientBalance, got %v",
			err,
		)
	}
}

func TestRepository_ApplyOperation_WalletNotFound(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	amount := int64(100)

	err := repo.ApplyOperation(
		context.Background(),
		uuid.New(),
		domain.OperationTypeDeposit,
		&amount,
	)

	if !errors.Is(err, repository.ErrWalletNotFound) {
		t.Fatalf(
			"expected ErrWalletNotFound, got %v",
			err,
		)
	}
}

func TestRepository_ApplyOperation_InvalidAmount(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	walletID := createWallet(t, db, 1000)
	defer deleteWallet(t, db, walletID)

	amount := int64(0)

	err := repo.ApplyOperation(
		context.Background(),
		walletID,
		domain.OperationTypeDeposit,
		&amount,
	)

	if !errors.Is(err, repository.ErrInvalidAmount) {
		t.Fatalf(
			"expected ErrInvalidAmount, got %v",
			err,
		)
	}
}

func TestRepository_ApplyOperation_InvalidOperation(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	walletID := createWallet(t, db, 1000)
	defer deleteWallet(t, db, walletID)

	amount := int64(100)

	err := repo.ApplyOperation(
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