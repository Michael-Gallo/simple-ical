.PHONY: test pre-commit

test:
	go test ./... --race

pre-commit: test
