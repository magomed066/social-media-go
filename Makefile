include .env

MIGRATION_PATH = ./cmd/migrate/migrations
DATABASE_URL = $(patsubst "%",%,$(patsubst '%',%,$(DB_ADDR)))
MIGRATE = migrate -path="$(MIGRATION_PATH)" -database="$(DATABASE_URL)"
SWAG ?= $(shell go env GOPATH)/bin/swag
steps ?= 1
MIGRATION_NAME = $(word 2,$(MAKECMDGOALS))

.PHONY: migration migrate-create migrate-up migrate-down
migration migrate-create:
	@migrate create -seq -ext sql -dir "$(MIGRATION_PATH)" "$(MIGRATION_NAME)"

.PHONY: seed
seed:
	@go run ./cmd/migrate/seed

migrate-up:
	@$(MIGRATE) up

migrate-down:
	@$(MIGRATE) down $(steps)

.PHONY: gen-docs
gen-docs:
	@"$(SWAG)" fmt -d cmd/api,internal
	@"$(SWAG)" init -g main.go -d cmd/api,internal --parseInternal -o docs

# Treat the migration name as an argument, not another target.
ifneq ($(filter migration migrate-create,$(firstword $(MAKECMDGOALS))),)
.PHONY: $(MIGRATION_NAME)
$(MIGRATION_NAME): ; @:
endif
