# Benchmarks comparing rs/cors and jub0bs/fcors

This repo contains benchmarks that compare the performance
of two CORS middleware libraries:

- the more popular [rs/cors](https://github.com/rs/cors) (v1.10.1), and
- the more modern and user-friendly [jub0bs/fcors](https://github.com/jub0bs/fcors) (v0.7.0).

## Running the benchmarks

Run the following commands in your shell.

```shell
git clone https://github.com/jub0bs/fcors-benchmarks
cd fcors-benchmarks
go test -run ^$ -bench . -benchmem -benchtime 1000000x
```

## Some results

I've slightly redacted the results below for better readability.
In particular, I've added a red dot next to cases where jub0bs/fcors
fares worse than rs/cors, and a green dot otherwise.

```text
goos: darwin
goarch: amd64
pkg: github.com/jub0bs/fcors-benchmarks
cpu: Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
identity_____________________vs_actual-8      4 ns/op     0 B/op   0 allocs/op
rs_cors_______________single_vs_actual-8    455 ns/op   352 B/op   1 allocs/op
jub0bs_fcors__________single_vs_actual-8    185 ns/op   352 B/op   1 allocs/op  🟢
rs_cors_____________multiple_vs_actual-8    265 ns/op   352 B/op   1 allocs/op
jub0bs_fcors________multiple_vs_actual-8    323 ns/op   352 B/op   1 allocs/op  🔴
rs_cors_________pathological_vs_actual-8    696 ns/op   384 B/op   3 allocs/op 
jub0bs_fcors____pathological_vs_actual-8   3135 ns/op   352 B/op   1 allocs/op  🔴
rs_cors_________________many_vs_actual-8    765 ns/op   352 B/op   1 allocs/op 
jub0bs_fcors____________many_vs_actual-8    338 ns/op   352 B/op   1 allocs/op  🟢
rs_cors__________________any_vs_actual-8    241 ns/op   352 B/op   1 allocs/op 
jub0bs_fcors_____________any_vs_actual-8    201 ns/op   352 B/op   1 allocs/op  🟢
rs_cors____________single_vs_preflight-8   1254 ns/op   800 B/op   4 allocs/op 
jub0bs_fcors_______single_vs_preflight-8    763 ns/op   784 B/op   4 allocs/op  🟢
rs_cors__________multiple_vs_preflight-8   1034 ns/op   800 B/op   4 allocs/op 
jub0bs_fcors_____multiple_vs_preflight-8    821 ns/op   784 B/op   4 allocs/op  🟢
rs_cors______pathological_vs_preflight-8   1175 ns/op   800 B/op   6 allocs/op 
jub0bs_fcors_pathological_vs_preflight-8   3535 ns/op   768 B/op   4 allocs/op  🔴
rs_cors______________many_vs_preflight-8   1588 ns/op   800 B/op   4 allocs/op 
jub0bs_fcors_________many_vs_preflight-8    857 ns/op   784 B/op   4 allocs/op  🟢
rs_cors_______________any_vs_preflight-8   1041 ns/op   800 B/op   4 allocs/op 
jub0bs_fcors__________any_vs_preflight-8    814 ns/op   784 B/op   4 allocs/op  🟢
rs_cors_______any_1header_vs_preflight-8   1131 ns/op   816 B/op   4 allocs/op 
jub0bs_fcors__any_1header_vs_preflight-8    865 ns/op   800 B/op   4 allocs/op  🟢

```
