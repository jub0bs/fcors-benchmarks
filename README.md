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
without_CORS      205514234        6 ns/op       0 B/op      0 allocs/op
```

### rs/cors (v1.10.0)

```text
single_origin____vs_actual                   7139101      168 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_actual                   6878386      173 ns/op       0 B/op      0 allocs/op
pathological_____vs_actual                   1946841      613 ns/op      32 B/op      2 allocs/op
many_origins_____vs_actual                   1828209      654 ns/op       0 B/op      0 allocs/op
any_origin_______vs_actual                   9203700      129 ns/op       0 B/op      0 allocs/op
single_origin____vs_preflight                3326931      360 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_preflight                3280380      369 ns/op       0 B/op      0 allocs/op
pathological_____vs_preflight                1668750      703 ns/op      32 B/op      2 allocs/op
many_origins_____vs_preflight                1414594      850 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight                3726684      334 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight_one_header     3103416      380 ns/op       0 B/op      0 allocs/op
```

### jub0bs/fcors (v0.6.0)

```text
pathological_____vs_actual                  19591662       61 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_actual                   5750486      201 ns/op       0 B/op      0 allocs/op
pathological_____vs_actual                    378091     3223 ns/op       0 B/op      0 allocs/op
many_origins_____vs_actual                   5524437      215 ns/op       0 B/op      0 allocs/op
any_origin_______vs_actual                  19334355       62 ns/op       0 B/op      0 allocs/op
single_origin____vs_preflight                6386853      173 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_preflight                4372368      242 ns/op       0 B/op      0 allocs/op
pathological_____vs_preflight                 377569     3311 ns/op       0 B/op      0 allocs/op
many_origins_____vs_preflight                4686819      254 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight                7188351      167 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight_one_header     6407846      186 ns/op       0 B/op      0 allocs/op
```
