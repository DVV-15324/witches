.PHONY: all \
	test \
	race \
	test-db-connect \
	test-handler-core \
	test-handler \
	test-utils \
	test-response \
	test-logger \
	test-easyjson \
	test-templates \
	test-cmd-utils \
	test-cmd \
	test-integration \
	coverage \
	coverage-summary \
	clean

# ============================================================
# Default
# ============================================================

all: race
	@echo "All tests completed!"

# ============================================================
# Individual tests
# ============================================================

test-db-connect:
	go test -race -covermode=atomic -coverprofile=coverage_db_connect.out ./pkg/core/database/connect/sql/...
	go tool cover -html=coverage_db_connect.out -o coverage_db_connect.html
	go tool cover -func=coverage_db_connect.out | findstr total

test-handler-core:
	go test -race -covermode=atomic -coverprofile=coverage_handler_core.out ./pkg/core/handler/...
	go tool cover -html=coverage_handler_core.out -o coverage_handler_core.html
	go tool cover -func=coverage_handler_core.out | findstr total

test-handler:
	go test -race -covermode=atomic -coverprofile=coverage_handler.out ./test/handler/...
	go tool cover -html=coverage_handler.out -o coverage_handler.html
	go tool cover -func=coverage_handler.out | findstr total

test-utils:
	go test -race -covermode=atomic -coverprofile=coverage_utils.out ./pkg/core/utils
	go tool cover -html=coverage_utils.out -o coverage_utils.html
	go tool cover -func=coverage_utils.out | findstr total

test-response:
	go test -race -covermode=atomic -coverprofile=coverage_response.out ./pkg/core/response/...
	go tool cover -html=coverage_response.out -o coverage_response.html
	go tool cover -func=coverage_response.out | findstr total

test-logger:
	go test -race -covermode=atomic -coverprofile=coverage_logger.out ./pkg/core/response/logger
	go tool cover -html=coverage_logger.out -o coverage_logger.html
	go tool cover -func=coverage_logger.out | findstr total

test-easyjson:
	go test -race -covermode=atomic -coverprofile=coverage_easyjson.out ./pkg/core/easyjson
	go tool cover -html=coverage_easyjson.out -o coverage_easyjson.html
	go tool cover -func=coverage_easyjson.out | findstr total

test-templates:
	go test -race -covermode=atomic -coverprofile=coverage_templates.out ./pkg/core/templates/...
	go tool cover -html=coverage_templates.out -o coverage_templates.html
	go tool cover -func=coverage_templates.out | findstr total

test-cmd-utils:
	go test -race -covermode=atomic -coverprofile=coverage_cmdutils.out ./cmd/utils
	go tool cover -html=coverage_cmdutils.out -o coverage_cmdutils.html
	go tool cover -func=coverage_cmdutils.out | findstr total

test-cmd:
	go test -race -covermode=atomic -coverprofile=coverage_cmd.out ./cmd/root_test.go
	go tool cover -html=coverage_cmd.out -o coverage_cmd.html
	go tool cover -func=coverage_cmd.out | findstr total

test-cmd-database:
	go test -race -covermode=atomic -coverprofile=coverage_cmd_database.out ./cmd/database
	go tool cover -html=coverage_cmd_database.out -o coverage_cmd_database.html
	go tool cover -func=coverage_cmd_database.out | findstr total

test-integration:
	go test -race -covermode=atomic -coverprofile=coverage_integration.out -tags=integration ./pkg/core/utils
	go tool cover -html=coverage_integration.out -o coverage_integration.html
	go tool cover -func=coverage_integration.out | findstr total

# ============================================================
# Run all tests
# ============================================================

test: clean
	@echo "========================================"
	@echo "Running all tests..."
	@echo "========================================"

	go test -race -covermode=atomic -coverprofile=coverage_db_connect.out ./pkg/core/database/connect/sql/...
	go test -race -covermode=atomic -coverprofile=coverage_handler_core.out ./pkg/core/handler/...
	go test -race -covermode=atomic -coverprofile=coverage_handler.out ./test/handler/...
	go test -race -covermode=atomic -coverprofile=coverage_easyjson.out ./pkg/core/easyjson
	go test -race -covermode=atomic -coverprofile=coverage_logger.out ./pkg/core/response/logger
	go test -race -covermode=atomic -coverprofile=coverage_response.out ./pkg/core/response/...
	go test -race -covermode=atomic -coverprofile=coverage_utils.out ./pkg/core/utils
	go test -race -covermode=atomic -coverprofile=coverage_templates.out ./pkg/core/templates/...
	go test -race -covermode=atomic -coverprofile=coverage_cmd.out ./cmd/...
	go test -race -covermode=atomic -coverprofile=coverage_cmdutils.out ./cmd/utils
	go test -race -covermode=atomic -coverprofile=coverage_cmd_database.out ./cmd/database
	go test -race -covermode=atomic -coverprofile=coverage_integration.out -tags=integration ./pkg/core/utils

	@echo ""
	@echo "========================================"
	@echo "All tests passed!"
	@echo "========================================"

# ============================================================
# Race
# ============================================================

race: test
	@echo ""
	@echo "========================================"
	@echo "Race tests completed successfully!"
	@echo "========================================"

# ============================================================
# Coverage
# ============================================================

coverage:
	@echo "Generating coverage reports..."

	@if exist coverage_db_connect.out ( \
		go tool cover -html=coverage_db_connect.out -o coverage_db_connect.html \
	) else (echo "No db connect coverage")

	@if exist coverage_handler_core.out ( \
		go tool cover -html=coverage_handler_core.out -o coverage_handler_core.html \
	) else (echo "No handler core coverage")

	@if exist coverage_handler.out ( \
		go tool cover -html=coverage_handler.out -o coverage_handler.html \
	) else (echo "No handler coverage")

	@if exist coverage_easyjson.out ( \
		go tool cover -html=coverage_easyjson.out -o coverage_easyjson.html \
	) else (echo "No easyjson coverage")

	@if exist coverage_logger.out ( \
		go tool cover -html=coverage_logger.out -o coverage_logger.html \
	) else (echo "No logger coverage")

	@if exist coverage_response.out ( \
		go tool cover -html=coverage_response.out -o coverage_response.html \
	) else (echo "No response coverage")

	@if exist coverage_utils.out ( \
		go tool cover -html=coverage_utils.out -o coverage_utils.html \
	) else (echo "No utils coverage")

	@if exist coverage_templates.out ( \
		go tool cover -html=coverage_templates.out -o coverage_templates.html \
	) else (echo "No templates coverage")

	@if exist coverage_cmd.out ( \
		go tool cover -html=coverage_cmd.out -o coverage_cmd.html \
	) else (echo "No cmd coverage")

	@if exist coverage_cmdutils.out ( \
		go tool cover -html=coverage_cmdutils.out -o coverage_cmdutils.html \
	) else (echo "No cmd-utils coverage")

	@if exist coverage_cmd_database.out ( \
		go tool cover -html=coverage_cmd_database.out -o coverage_cmd_database.html \
	) else (echo "No cmd-database coverage")

	@if exist coverage_migrate.out ( \
		go tool cover -html=coverage_migrate.out -o coverage_migrate.html \
	) else (echo "No migrate coverage")

	@if exist coverage_integration.out ( \
		go tool cover -html=coverage_integration.out -o coverage_integration.html \
	) else (echo "No integration coverage")

	@echo ""
	@echo "Coverage HTML reports generated."

# ============================================================
# Coverage Summary
# ============================================================

coverage-summary:
	@echo ""
	@echo "========================================"
	@echo "         COVERAGE SUMMARY"
	@echo "========================================"

	@echo ""
	@echo "[DB Connect]"
	@if exist coverage_db_connect.out (go tool cover -func=coverage_db_connect.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[Handler Core]"
	@if exist coverage_handler_core.out (go tool cover -func=coverage_handler_core.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[Handler]"
	@if exist coverage_handler.out (go tool cover -func=coverage_handler.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[EasyJSON]"
	@if exist coverage_easyjson.out (go tool cover -func=coverage_easyjson.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[Logger]"
	@if exist coverage_logger.out (go tool cover -func=coverage_logger.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[Response]"
	@if exist coverage_response.out (go tool cover -func=coverage_response.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[Utils]"
	@if exist coverage_utils.out (go tool cover -func=coverage_utils.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[Templates]"
	@if exist coverage_templates.out (go tool cover -func=coverage_templates.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[CMD]"
	@if exist coverage_cmd.out (go tool cover -func=coverage_cmd.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[CMD Utils]"
	@if exist coverage_cmdutils.out (go tool cover -func=coverage_cmdutils.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[CMD Database]"
	@if exist coverage_cmd_database.out (go tool cover -func=coverage_cmd_database.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[Migrate]"
	@if exist coverage_migrate.out (go tool cover -func=coverage_migrate.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "[Integration]"
	@if exist coverage_integration.out (go tool cover -func=coverage_integration.out | findstr total) else (echo "No coverage file")

	@echo ""
	@echo "========================================"

# ============================================================
# Clean
# ============================================================

clean:
	@echo "Cleaning coverage files..."
	-del coverage*.out 2>nul
	-del coverage*.html 2>nul
	-del coverage_all.out 2>nul

	@echo "Clean completed."