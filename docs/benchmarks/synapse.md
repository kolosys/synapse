# synapse Benchmarks

Performance benchmarks for the synapse package.

**Import Path:** `github.com/kolosys/synapse`

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
go test -bench=. -benchmem ./github.com/kolosys/synapse
```

To run benchmarks for all packages:

```bash
go test -bench=. -benchmem ./...
```

## External Links

- [Package Overview](../core-concepts/synapse.md)
- [API Reference](../api-reference/synapse.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/synapse)
- [Source Code](https://github.com/kolosys/synapse/tree/main/github.com/kolosys/synapse)
