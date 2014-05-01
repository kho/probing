# Linear probing hash tables in Go

This is a template for linear probing hash tables. Because this is a
probing hash table, one can only insert or modify entries (i.e. no
deletion). And there must be one invalid key `__KEY_NIL` that is used
to mark unused buckets. The advantages (vs. Go's own map type) are:

1. All entries are stored in a single slice;
2. Memory footprint is generally smaller;
3. Look-up can be often faster or not much slower.

To customize,

1. Define your own key, value types, etc. in `probing_params.go`;
2. Replace `__` in `probing_params.go` and `probing_impl.go` with name
   of your choice;
3. Copy these two files to your project and you are good to go!

## Benchmark results (look-ups only)

Summary: the probing hash table is,

1. Slightly slower with much smaller memory footprint;
2. Significantly faster with comparable memory usage;
3. Faster on larger maps.

```
$ for s in 1.5 2; do for k in 1000000 10000000; do go test -bench=. -gcflags='-l -l -l -l -l' -keys=$k -scale=$s; done; done
PASS
BenchmarkGoMap	   20000	     97371 ns/op
BenchmarkProbingMap	   10000	    109566 ns/op
BenchmarkMapMem	       0	         0 ns/op
--- BENCH: BenchmarkMapMem
	probing_bench_test.go:111: go-map: 36.74MB
	probing_bench_test.go:129: probing: 1000000 keys; 1500000 buckets; load factor 0.67
	probing_bench_test.go:119: probing: 22.88MB
ok  	github.com/kho/probing	7.009s
PASS
BenchmarkGoMap	   10000	    120579 ns/op
BenchmarkProbingMap	   10000	    124704 ns/op
BenchmarkMapMem	       0	         0 ns/op
--- BENCH: BenchmarkMapMem
	probing_bench_test.go:111: go-map: 304.63MB
	probing_bench_test.go:129: probing: 10000000 keys; 15000000 buckets; load factor 0.67
	probing_bench_test.go:119: probing: 228.51MB
ok  	github.com/kho/probing	32.110s
PASS
BenchmarkGoMap	   20000	     98238 ns/op
BenchmarkProbingMap	   20000	     86263 ns/op
BenchmarkMapMem	       0	         0 ns/op
--- BENCH: BenchmarkMapMem
	probing_bench_test.go:111: go-map: 36.73MB
	probing_bench_test.go:129: probing: 1000000 keys; 2000000 buckets; load factor 0.50
	probing_bench_test.go:119: probing: 30.51MB
ok  	github.com/kho/probing	8.743s
PASS
BenchmarkGoMap	   10000	    121225 ns/op
BenchmarkProbingMap	   20000	    100018 ns/op
BenchmarkMapMem	       0	         0 ns/op
--- BENCH: BenchmarkMapMem
	probing_bench_test.go:111: go-map: 304.67MB
	probing_bench_test.go:129: probing: 10000000 keys; 20000000 buckets; load factor 0.50
	probing_bench_test.go:119: probing: 304.80MB
ok  	github.com/kho/probing	36.224s
```
