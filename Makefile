.PHONY: mocks test test-unit test-integration test-coverage lint

mocks:
	mockery

test: test-unit test-integration

test-unit:
	go test ./... -count=1 -v -short

test-integration:
	go test ./... -count=1 -v -tags=integration -run Integration

test-coverage:
	go test ./... -count=1 -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...
