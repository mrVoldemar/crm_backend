include .env

LOCAL_BIN:=$(CURDIR)/bin
GOOSE=./bin/goose

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	$(GOOSE) -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(GOOSE) -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(GOOSE) -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

migration-create:
	@echo "Enter migration name (lowercase, words separated by underscore): "
	@read -p "> " MIGRATE_NAME && \
	if [ -z "$$MIGRATE_NAME" ]; then \
		echo "Error: Migration name cannot be empty"; \
		exit 1; \
	fi && \
	$(GOOSE) -dir ./migrations create "$$MIGRATE_NAME" sql

