.PHONY: test pre-commit lint fmt vet

test:
	go test ./... --race

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

pre-commit: fmt vet test
