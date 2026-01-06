# eviction Benchmarks

Performance benchmarks for the eviction package.

**Import Path:** `github.com/kolosys/synapse/eviction`

## No Benchmarks Available

No benchmark results are available for this package. To add benchmarks:

1. Create a `*_test.go` file in the package directory
2. Add benchmark functions following the pattern:
   ```go
   func BenchmarkFunctionName(b *testing.B) {
       for i := 0; i < b.N; i++ {
           // Your code here
       }
   }
   ```
3. Run `proton benchmark` to generate benchmark results

## Running Benchmarks

To run benchmarks for this package:

```bash
go test -bench=. -benchmem ./eviction
```

To run benchmarks for all packages:

```bash
go test -bench=. -benchmem ./...
```

## External Links

- [Package Overview](../core-concepts/eviction.md)
- [API Reference](../api-reference/eviction.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/synapse/eviction)
- [Source Code](https://github.com/kolosys/synapse/tree/main/eviction)
