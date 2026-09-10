include config.env
export

export PROJECT_ROOT := $(shell pwd)

up:
	@docker compose --env-file config.env up --build

down:
	@docker compose --env-file config.env down

restart:
	@docker compose --env-file config.env down
	@docker compose --env-file config.env up --build