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
without_CORS                               205514234        6 ns/op       0 B/op      0 allocs/op
```

### rs/cors

```text
single_origin____vs_actual                   3194002      375 ns/op      48 B/op      3 allocs/op
multiple_origins_vs_actual                   3152779      382 ns/op      48 B/op      3 allocs/op
many_origins_____vs_actual                   1246867     1026 ns/op      48 B/op      3 allocs/op
any_origin_______vs_actual                   3499288      342 ns/op      48 B/op      3 allocs/op
single_origin____vs_preflight                1223632      983 ns/op     160 B/op      6 allocs/op
multiple_origins_vs_preflight                1214544      986 ns/op     160 B/op      6 allocs/op
many_origins_____vs_preflight                 765120     1574 ns/op     160 B/op      6 allocs/op
any_origin_______vs_preflight                1264927      951 ns/op     160 B/op      6 allocs/op
any_origin_______vs_preflight_one_header      907146     1302 ns/op     208 B/op     10 allocs/op
```

### jub0bs/fcors

```text
single_origin____vs_actual                  21439983       53 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_actual                   6499868      184 ns/op       0 B/op      0 allocs/op
many_origins_____vs_actual                   5030031      216 ns/op       0 B/op      0 allocs/op
any_origin_______vs_actual                  21453963       53 ns/op       0 B/op      0 allocs/op
single_origin____vs_preflight                7923229      152 ns/op       0 B/op      0 allocs/op
multiple_origins_vs_preflight                4998284      222 ns/op       0 B/op      0 allocs/op
many_origins_vs_preflight                    4601485      259 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight                8044881      149 ns/op       0 B/op      0 allocs/op
any_origin_______vs_preflight_one_header     7042791      171 ns/op       0 B/op      0 allocs/op
```
