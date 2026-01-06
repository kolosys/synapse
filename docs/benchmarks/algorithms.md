# algorithms Benchmarks

Performance benchmarks for the algorithms package.

**Import Path:** `github.com/kolosys/synapse/algorithms`

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
go test -bench=. -benchmem ./algorithms
```

To run benchmarks for all packages:

```bash
go test -bench=. -benchmem ./...
```

## External Links

- [Package Overview](../core-concepts/algorithms.md)
- [API Reference](../api-reference/algorithms.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/synapse/algorithms)
- [Source Code](https://github.com/kolosys/synapse/tree/main/algorithms)
