build-backend:
	@go build -o bin/backend

run-backend: build-backend
	@./bin/backend

MOCKGEN=mockgen
OUTPUT_DIR=mock
# Find all directories with Go files
DIRS=$(shell find . -type d -not -path "./$(OUTPUT_DIR)*" -not -path "./vendor*" -not -path "./api*")

# Generate mocks
# this will skip generating mocks for the gen and frontend directories
generate-mocks: clean-mocks
	@for dir in $(DIRS); do \
		if [[ "$$dir" != *"/gen"* && "$$dir" != *"/frontend"* && "$$dir" != *"/$(OUTPUT_DIR)"* ]] && grep -q "interface" "$$dir"/*.go 2>/dev/null; then \
			echo "Generating mocks in $$dir"; \
			for file in "$$dir"/*.go; do \
				if ! grep -q "type.*interface" "$$file"; then \
					continue; \
				fi; \
				mkdir -p "$$dir/$(OUTPUT_DIR)"; \
				$(MOCKGEN) -source "$$file" -destination "$$dir/$(OUTPUT_DIR)/$$(basename $$file .go)_mock.go"; \
			done \
		fi \
	done

# Clean generated mocks
clean-mocks:
	find . -type d -name "$(OUTPUT_DIR)" | xargs rm -rf

tests:
	@go test ./... --coverprofile=coverage.out && go tool cover -html=coverage.out

migrate-up:
	@cd ./migration && go run migrate.go up

migrate-down:
	@cd .//migration && go run migrate.go down