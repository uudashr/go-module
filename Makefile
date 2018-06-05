# Linter
.PHONY: lint-prepare
lint-prepare:
	@echo "Installing golangci-lint"
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint:
	@golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...

# Testing
.PHONY: test
test:
	@go test $(TEST_OPTS)
