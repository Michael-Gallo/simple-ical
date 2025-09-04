.PHONY: test pre-commit lint fmt vet

test:
	go test ./... --race --count 1

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

pre-commit: fmt vet test
