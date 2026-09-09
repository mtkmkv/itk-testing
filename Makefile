include config.env
export

export PROJECT_ROOT := $(shell pwd)

env-up:
	@docker compose up -d postgres

env-down:
	@docker compose down postgres

env-cleanup:
	@read -p "Are you sure you want to remove the postgres data? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down postgres; \
		rm -rf out/pgdata; \
		echo "Postgres data successfully removed."; \
	else \
		echo "Cleanup canceled."; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Please provide a migration name via 'seq' variable."; \
		echo "Example: make migrate-create seq=init_database"; \
		exit 1; \
	fi; \
	MSYS_NO_PATHCONV=1 docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make --no-print-directory migrate-action action=up

migrate-down:
	@make --no-print-directory migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Please provide an action (up or down) via 'action' variable."; \
		echo "Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	MSYS_NO_PATHCONV=1 docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable \
		"$(action)"

run:
	@go run ./cmd/server