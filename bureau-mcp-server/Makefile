include .env
export

BG_BLUE = \033[44m
BG_GREEN = \033[42m
BG_RED = \033[41m
BG_RESET = \033[0m
FG_WHITE = \033[97m
BG_GREY = \033[100m

# Connection string for the migrate CLI, derived from DATABASE_URL so the app
# and the migrations share a single source of truth.
#   * migrate needs the postgres:// scheme (not postgresql://)
#   * sslmode=disable is forced for local/dockerised Postgres
#   * outside production the migrate CLI runs on the host, so the docker service
#     address (postgres:5432) is rewritten to the exposed 127.0.0.1:5433
ifeq ($(ENV),production)
	MIGRATE_URL := $(shell echo '$(DATABASE_URL)' | sed 's|^postgresql://|postgres://|')?sslmode=disable
else
	MIGRATE_URL := $(shell echo '$(DATABASE_URL)' | sed -e 's|^postgresql://|postgres://|' -e 's|@postgres:5432/|@127.0.0.1:5433/|')?sslmode=disable
endif

documentation:
	@echo -e "$(BG_BLUE) $(FG_WHITE) Generating API documentation... $(BG_RESET)"
	@swag init --parseDependency --parseInternal -g cmd/api/main.go
	@echo -e "$(BG_GREEN) API documentation generated. $(BG_RESET)"

run:
	@docker compose up -d --build && \
	make migration-up

run-db:
	docker compose up -d pgadmin && \
	docker compose up -d postgres

migration-status:
	@echo -e "$(BG_BLUE) $(FG_WHITE) Checking database migration status... $(BG_RESET)"
	@printf "Migration list:\n" && ls internal/infra/thirdparty/database/migrations | grep up.sql || echo "No migration files found."
	@printf "Current version: " && migrate -path internal/infra/thirdparty/database/migrations -database "$(MIGRATE_URL)" version || echo "No migrations applied yet."
	@echo -e "$(BG_GREEN) Migration status checked. $(BG_RESET)"

migration-up:
	@echo -e "$(BG_BLUE) $(FG_WHITE) Running database migrations up... $(BG_RESET)"
	@migrate -path internal/infra/thirdparty/database/migrations -database "$(MIGRATE_URL)" up
	@echo -e "$(BG_GREEN) Migrations completed. $(BG_RESET)"

migration-up-single:
	@echo -e "$(BG_BLUE) $(FG_WHITE) Running single database migration... $(BG_RESET)"
	@migrate -path internal/infra/thirdparty/database/migrations -database "$(MIGRATE_URL)" up 1
	@echo -e "$(BG_GREEN) Single migration completed. $(BG_RESET)"

migration-down:
	@echo -e "$(BG_BLUE) $(FG_WHITE) Running database migrations down... $(BG_RESET)"
	@migrate -path internal/infra/thirdparty/database/migrations -database "$(MIGRATE_URL)" down
	@echo -e "$(BG_GREEN) Migrations rolled back. $(BG_RESET)"

migration-down-single:
	@echo -e "$(BG_BLUE) $(FG_WHITE) Running database migrations to previous version... $(BG_RESET)"
	@migrate -path internal/infra/thirdparty/database/migrations -database "$(MIGRATE_URL)" down 1
	@echo -e "$(BG_GREEN) Migration rolled back. $(BG_RESET)"

migration-down-fixed: migration-fix migration-down

migration-new:
	@read -p "Enter migration name (use underscores instead of spaces): " name; \
	echo -e "$(BG_BLUE) $(FG_WHITE)Creating $$name migration file... $(BG_RESET)"; \
	file=$$(migrate create -ext sql -dir internal/infra/thirdparty/database/migrations -seq $$name 2>&1 | tail -n1); \
	echo -e "$(BG_GREEN) File created. $(BG_RESET)"; \
	echo -e "$(BG_GREY) $(FG_WHITE)You can find the new migration file at: $$file $(BG_RESET)"

migration-fix:
	@read -p "Enter version to force: " version; \
	read -p "Are you sure you want to force the migration to version $$version? This can be dangerous. [y/N]: " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo -e "$(BG_BLUE) $(FG_WHITE) Forcing migration to version $$version... $(BG_RESET)"; \
		migrate -path internal/infra/thirdparty/database/migrations -database "$(MIGRATE_URL)" force $$version; \
		echo -e "$(BG_GREEN) Migration version forced.$(BG_RESET)"; \
	else \
		echo -e "$(BG_RED) $(FG_WHITE) Migration force aborted. $(BG_RESET)"; \
	fi