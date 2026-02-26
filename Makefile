
# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Run integration tests only (requires Docker)
test-integration:
	@echo "Running integration tests..."
	@go test ./tests/integration/... -v -count=1

lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Linting..."; \
		golangci-lint run ./...; \
	else \
		read -p "golangci-lint is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
			echo "Installed golangci-lint"; \
			golangci-lint run ./...; \
		else \
			echo "You chose not to install golangci-lint. Exiting..."; \
			exit 1; \
		fi; \
	fi

fix-lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Fixing lint issues..."; \
		golangci-lint run --fix ./...; \
	else \
		read -p "golangci-lint is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
			echo "Installed golangci-lint"; \
			golangci-lint run --fix ./...; \
		else \
			echo "You chose not to install golangci-lint. Exiting..."; \
			exit 1; \
		fi; \
	fi

clean:
	@echo "Cleaning..."
	@rm -f main

swag:
	@if command -v swag > /dev/null; then \
		echo "Generating Swagger docs..."; \
		swag init \
			--generalInfo applications/api/main.go \
			--dir . \
			--output docs \
			--parseDependency \
			--parseInternal; \
	else \
		read -p "swag is not installed. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/swaggo/swag/cmd/swag@latest; \
			swag init \
				--generalInfo applications/api/main.go \
				--dir . \
				--output docs \
				--parseDependency \
				--parseInternal; \
		else \
			echo "You chose not to install swag. Exiting..."; \
			exit 1; \
		fi; \
	fi

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                echo "Installed air"; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

lint-frontend:
	@npx -y react-doctor@latest frontend

migrate-db-up:
	@echo "Running migrations (up)..."
	@docker build -f Dockerfile.migrations -t leetcode-spaced-repetition-migrations .
	@docker run --rm --env-file .env --network leetcode-spaced-repetition_app-network leetcode-spaced-repetition-migrations -direction up

migrate-db-down:
	@echo "Running migrations (down)..."
	@docker build -f Dockerfile.migrations -t leetcode-spaced-repetition-migrations .
	@docker run --rm --env-file .env --network leetcode-spaced-repetition_app-network leetcode-spaced-repetition-migrations -direction down

seed-db:
	@echo "Seeding database with leetcode problems"
	@docker build -t data-seeding -f Dockerfile.data-seeding .
	@docker run --rm --name db-seeding --network leetcode-spaced-repetition_app-network data-seeding:latest

logs-app:
	@docker logs leetcode-buddy-app-frontend

logs-app-tail:
	@docker logs -f leetcode-buddy-app-frontend

logs-api:
	@docker logs leetcode-buddy-app-api

logs-api-tail:
	@docker logs -f leetcode-buddy-app-api
