# Benchmarks comparing rs/cors vs. jub0bs/fcors

This repo contains benchmarks that compare the performance
of two CORS middleware libraries:

- the more popular [rs/cors](https://github.com/rs/cors), and
- the more modern and user-friendly [jub0bs/fcors](https://github.com/jub0bs/fcors).

## Running the benchmarks

```shell
git clone https://github.com/jub0bs/fcors-benchmarks
cd fcors-benchmarks
go test -run ^$ -bench . -benchmem
```

## Some results

(I've slightly redacted the results below for better readability.)

### Machine details

```text
goos: darwin
goarch: amd64
pkg: github.com/jub0bs/fcors-benchmarks
cpu: Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
```

### Base case: no CORS middleware

```text
without_CORS      217746176        5 ns/op       0 B/op      0 allocs/op
```

### rs/cors (v1.10.1)

```text
single_origin____vs_actual                   8096647      148 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_actual                   7702515      156 ns/op       0 B/op      0 allocs/op
pathological_____vs_actual                   1916860      622 ns/op      32 B/op      2 allocs/op
many_origins_____vs_actual                   1865460      644 ns/op       0 B/op      0 allocs/op
any_origin_______vs_actual                  10262497      115 ns/op       0 B/op      0 allocs/op
single_origin____vs_preflight                3332138      359 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_preflight                3280244      365 ns/op       0 B/op      0 allocs/op
pathological_____vs_preflight                1671585      715 ns/op      32 B/op      2 allocs/op
many_origins_____vs_preflight                1393687      857 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight                3687146      326 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight_one_header     3099513      398 ns/op       0 B/op      0 allocs/op
```

### jub0bs/fcors (v0.7.0)

```text
single_origin____vs_actual                  18442140       64 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_actual                   5684814      209 ns/op       0 B/op      0 allocs/op
pathological_____vs_actual                    406676     2940 ns/op       0 B/op      0 allocs/op
many_origins_____vs_actual                   5345062      225 ns/op       0 B/op      0 allocs/op
any_origin_______vs_actual                  19009729       64 ns/op       0 B/op      0 allocs/op
single_origin____vs_preflight                7207303      172 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_preflight                5061046      235 ns/op       0 B/op      0 allocs/op
pathological_____vs_preflight                 405902     2915 ns/op       0 B/op      0 allocs/op
many_origins_____vs_preflight                4696436      255 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight                7315965      165 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight_one_header     6449601      187 ns/op       0 B/op      0 allocs/op
```
