SERVICES := imports

.PHONY: test coverage local

test:
	@for service in $(SERVICES); do \
		echo "Running tests for $$service..."; \
		cd $$service && go test ./... && cd ..; \
	done

coverage:
	@for service in $(SERVICES); do \
		echo "Generating coverage report for $$service..."; \
		cd $$service && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out && cd ..; \
	done

local:
	docker-compose up -d