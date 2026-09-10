package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	coreHTTP "github.com/mtkmkv/itk-testing/initial/core/transport/http"
	httpserver "github.com/mtkmkv/itk-testing/initial/core/transport/http/server"
	walletRepository "github.com/mtkmkv/itk-testing/initial/features/wallet/repository/postgres"
	walletService "github.com/mtkmkv/itk-testing/initial/features/wallet/service"
	walletHTTP "github.com/mtkmkv/itk-testing/initial/features/wallet/transport/http"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	dsn := "postgres://" +
		os.Getenv("POSTGRES_USER") + ":" +
		os.Getenv("POSTGRES_PASSWORD") + "@" +
		os.Getenv("POSTGRES_HOST") + ":" +
		os.Getenv("POSTGRES_PORT") + "/" +
		os.Getenv("POSTGRES_DB") +
		"?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	repository := walletRepository.NewRepository(db)
	service := walletService.NewService(repository)
	handler := walletHTTP.NewHandler(service)
	router := coreHTTP.NewRouter(handler)

	config, err := httpserver.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	server := httpserver.NewServer(config, router)

	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}