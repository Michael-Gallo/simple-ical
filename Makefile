.PHONY: test pre-commit lint fmt vet

test-slow:
	go test ./... --race --count 1

test:
	go test --count 1 ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

bench:
	cd benchmarks && go test -bench=BenchmarkAllScenarios -benchmem 

bench-long:
	cd benchmarks && go test -bench=BenchmarkAllScenarios -benchmem  -count 10 > results.txt

bench-comparative:
	cd benchmarks && go test -bench=BenchmarkComparativeAll -benchmem -count 10 > results_comparative.txt


pre-commit: fmt vet test-slow
