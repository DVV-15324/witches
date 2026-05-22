SHELL := bash
# TOOLS INSTALL
install:
	go install github.com/mailru/easyjson/...@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

DB_PROFILE ?=
# DOCKER INFRA
database-up-docker:
	cd ./internal/core/database && docker compose --profile $(DB_PROFILE) --profile redis --profile migrate up --build -d

database-down-docker:
	cd ./internal/core/database && docker compose down

run:
	easyjson init
	swag init
	go run .

migrate-init:
	@mkdir -p $(PWD)/migrate/migrations
# MIGRATE CONFIG
DB_URL ?=
# MIGRATE COMMANDS (DOCKER WAY)
migrate-up-docker:
	docker run --rm \
		-v $(PWD)/migrate/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" up 1

migrate-down-docker:
	docker run --rm \
		-v $(PWD)/migrate/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" down 1

migrate-version-docker:
	docker run --rm \
		-v $(PWD)/migrate/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" version

migrate-force-docker:
	docker run --rm \
		-v $(PWD)/migrate/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" force $(VERSION)

