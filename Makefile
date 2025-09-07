.PHONY: test pre-commit lint fmt vet

test-slow:
	go test ./... --race --count 1

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

bench:
	go test -bench=. ./... -benchmem

pre-commit: fmt vet test-slow
