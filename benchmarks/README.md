# iCalendar Parser Benchmarks

This directory contains comparative benchmarks against other Go iCalendar parsers.

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Install gocal:
```bash
go get github.com/apognu/gocal
```

## Running Benchmarks

### Run All Benchmarks
```bash
go test -bench=. -benchmem
```

### Run Specific Parser Benchmarks
```bash
# Your parser only
go test -bench=BenchmarkAllSimpleIcal -benchmem

# gocal only
go test -bench=BenchmarkAllGocal -benchmem

# Comparative benchmarks
go test -bench=BenchmarkComparative -benchmem
```

### Memory Usage Comparison
```bash
go test -bench=BenchmarkMemoryUsage -benchmem
```

### Run with CPU Profiling
```bash
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### Run with Memory Profiling
```bash
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
```

## Benchmark Results

The benchmarks compare:
- **Parsing Speed**: Time per operation
- **Memory Usage**: Bytes allocated per operation
- **Allocation Count**: Number of allocations per operation

## Adding More Parsers

To add another parser:

1. Add the dependency to `go.mod`
2. Create a new benchmark file (e.g., `other_parser_benchmark.go`)
3. Implement benchmark functions following the same pattern
4. Add to `BenchmarkComparative` function

## Test Data

Test data is loaded from `../parse/test_data/` directory. To add new test cases:

1. Add the `.ical` file to the test data directory
2. Update the `testCases` slice in `LoadTestData()` function

