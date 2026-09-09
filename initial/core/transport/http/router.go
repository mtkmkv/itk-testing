package core_transport_http

import (
	"net/http"

	walletHTTP "github.com/mtkmkv/itk-testing/initial/features/wallet/transport/http"
)

func NewRouter(walletHandler *walletHTTP.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/wallet", walletHandler.Operate)
	mux.HandleFunc("GET /api/v1/wallets/{walletID}", walletHandler.GetBalance)

	return mux
}