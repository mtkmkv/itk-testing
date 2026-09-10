package feature_wallet_transport_http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/mtkmkv/itk-testing/initial/core/domain"
	repository "github.com/mtkmkv/itk-testing/initial/features/wallet/repository"
	service "github.com/mtkmkv/itk-testing/initial/features/wallet/service"
)

type mockRepository struct {
	getByIDFunc        func(context.Context, uuid.UUID) (domain.Wallet, error)
	applyOperationFunc func(context.Context, uuid.UUID, domain.OperationType, *int64) error
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
	return m.applyOperationFunc(
		ctx,
		walletID,
		operationType,
		amount,
	)
}

func TestHandler_Operate(t *testing.T) {
	walletID := uuid.New()

	repo := &mockRepository{
		applyOperationFunc: func(
			ctx context.Context,
			id uuid.UUID,
			operationType domain.OperationType,
			amount *int64,
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

			if amount == nil || *amount != 1000 {
				t.Fatalf("expected amount 1000, got %v", amount)
			}

			return nil
		},
	}

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	body := `{
		"walletId": "` + walletID.String() + `",
		"operationType": "DEPOSIT",
		"amount": 1000
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Operate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestHandler_Operate_InvalidJSON(t *testing.T) {
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

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		strings.NewReader(`invalid json`),
	)

	rec := httptest.NewRecorder()

	handler.Operate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}

	var response errorResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Error != "invalid request body" {
		t.Fatalf(
			"expected error %q, got %q",
			"invalid request body",
			response.Error,
		)
	}
}

func TestHandler_Operate_InvalidWalletID(t *testing.T) {
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

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	body := `{
		"walletId": "invalid",
		"operationType": "DEPOSIT",
		"amount": 1000
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Operate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestHandler_Operate_InvalidOperation(t *testing.T) {
	walletID := uuid.New()

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

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	body := `{
		"walletId": "` + walletID.String() + `",
		"operationType": "INVALID",
		"amount": 1000
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Operate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestHandler_Operate_InvalidAmount(t *testing.T) {
	walletID := uuid.New()

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

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	body := `{
		"walletId": "` + walletID.String() + `",
		"operationType": "DEPOSIT",
		"amount": 0
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Operate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestHandler_Operate_WalletNotFound(t *testing.T) {
	walletID := uuid.New()

	repo := &mockRepository{
		applyOperationFunc: func(
			ctx context.Context,
			id uuid.UUID,
			operationType domain.OperationType,
			amount *int64,
		) error {
			return repository.ErrWalletNotFound
		},
	}

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	body := `{
		"walletId": "` + walletID.String() + `",
		"operationType": "DEPOSIT",
		"amount": 1000
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Operate(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestHandler_Operate_InsufficientBalance(t *testing.T) {
	walletID := uuid.New()

	repo := &mockRepository{
		applyOperationFunc: func(
			ctx context.Context,
			id uuid.UUID,
			operationType domain.OperationType,
			amount *int64,
		) error {
			return repository.ErrInsufficientBalance
		},
	}

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	body := `{
		"walletId": "` + walletID.String() + `",
		"operationType": "WITHDRAW",
		"amount": 1000
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Operate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestHandler_Operate_InternalError(t *testing.T) {
	walletID := uuid.New()

	repo := &mockRepository{
		applyOperationFunc: func(
			ctx context.Context,
			id uuid.UUID,
			operationType domain.OperationType,
			amount *int64,
		) error {
			return errors.New("database error")
		},
	}

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	body := `{
		"walletId": "` + walletID.String() + `",
		"operationType": "DEPOSIT",
		"amount": 1000
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Operate(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestHandler_GetBalance(t *testing.T) {
	walletID := uuid.New()

	repo := &mockRepository{
		getByIDFunc: func(
			ctx context.Context,
			id uuid.UUID,
		) (domain.Wallet, error) {
			return domain.Wallet{
				ID:        walletID,
				Balance:   1500,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/wallets/"+walletID.String(),
		nil,
	)

	req.SetPathValue("walletID", walletID.String())

	rec := httptest.NewRecorder()

	handler.GetBalance(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response walletResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.WalletID != walletID {
		t.Fatalf(
			"expected wallet ID %s, got %s",
			walletID,
			response.WalletID,
		)
	}

	if response.Balance != 1500 {
		t.Fatalf(
			"expected balance 1500, got %d",
			response.Balance,
		)
	}
}

func TestHandler_GetBalance_InvalidWalletID(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(
			ctx context.Context,
			id uuid.UUID,
		) (domain.Wallet, error) {
			t.Fatal("repository should not be called")
			return domain.Wallet{}, nil
		},
	}

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/wallets/invalid",
		nil,
	)

	req.SetPathValue("walletID", "invalid")

	rec := httptest.NewRecorder()

	handler.GetBalance(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestHandler_GetBalance_WalletNotFound(t *testing.T) {
	walletID := uuid.New()

	repo := &mockRepository{
		getByIDFunc: func(
			ctx context.Context,
			id uuid.UUID,
		) (domain.Wallet, error) {
			return domain.Wallet{}, repository.ErrWalletNotFound
		},
	}

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/wallets/"+walletID.String(),
		nil,
	)

	req.SetPathValue("walletID", walletID.String())

	rec := httptest.NewRecorder()

	handler.GetBalance(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestHandler_GetBalance_InternalError(t *testing.T) {
	walletID := uuid.New()

	repo := &mockRepository{
		getByIDFunc: func(
			ctx context.Context,
			id uuid.UUID,
		) (domain.Wallet, error) {
			return domain.Wallet{}, errors.New("database error")
		},
	}

	svc := service.NewService(repo)
	handler := NewHandler(svc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/wallets/"+walletID.String(),
		nil,
	)

	req.SetPathValue("walletID", walletID.String())

	rec := httptest.NewRecorder()

	handler.GetBalance(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}