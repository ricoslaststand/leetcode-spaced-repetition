
# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

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

clean:
	@echo "Cleaning..."
	@rm -f main

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