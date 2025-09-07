.PHONY: test pre-commit lint fmt vet

test:
	@git stash push -m "temp stash for test" --keep-index
	go test ./... --race --count 1
	@git stash pop

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

bench:
	go test -bench=. ./... -benchmem

pre-commit: fmt vet test
