# Blossom Filter

Generic-type Bloom filter implementation in Go.

## Usage

See [examples](./example) for usage scenarios.

## Implementation Details

- **Hash Function**: Uses SipHash for cryptographically strong hashing
- **Optimal Hash Count**: Dynamically calculated as `(filterSize/elements) * ln(2)`
- **Bitset Storage**: Efficient packed byte array with bit-level operations
- **False Positive Rate**: Adjusts automatically based on filter size and element count

## Testing

Run tests with coverage:

```bash
go test ./... -cover
```

Generate HTML coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Current coverage: **100%** for both `pkg/core` and `internal/bitset`

## Benchmarks

```bash
go test -bench=. -benchmem ./pkg/core
```

Typical performance:
- Int operations: ~4 us/op
- String operations: ~900 us/op
- Bitset operations: ~1.6 ns/op



