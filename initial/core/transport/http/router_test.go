package core_transport_http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mtkmkv/itk-testing/initial/core/domain"
	feature_wallet_service "github.com/mtkmkv/itk-testing/initial/features/wallet/service"
	walletHTTP "github.com/mtkmkv/itk-testing/initial/features/wallet/transport/http"
)

type mockRepository struct{}

func (m *mockRepository) GetByID(
	_ context.Context,
	_ uuid.UUID,
) (domain.Wallet, error) {
	return domain.Wallet{}, nil
}

func (m *mockRepository) ApplyOperation(
	_ context.Context,
	_ uuid.UUID,
	_ domain.OperationType,
	_ *int64,
) error {
	return nil
}

func TestNewRouter(t *testing.T) {
	repo := &mockRepository{}
	svc := feature_wallet_service.NewService(repo)
	handler := walletHTTP.NewHandler(svc)

	router := NewRouter(handler)

	walletID := uuid.New()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "post wallet",
			method:     http.MethodPost,
			path:       "/api/v1/wallet",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "get wallet",
			method:     http.MethodGet,
			path:       "/api/v1/wallets/" + walletID.String(),
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong method for post",
			method:     http.MethodGet,
			path:       "/api/v1/wallet",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "wrong method for get",
			method:     http.MethodPost,
			path:       "/api/v1/wallets/" + walletID.String(),
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "unknown route",
			method:     http.MethodGet,
			path:       "/api/v1/unknown",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}
		})
	}
}