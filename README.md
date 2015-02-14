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
$ for s in 1.5 2; do for k in 1000 1000000 10000000; do go test -bench=. -gcflags='-l -l -l -l -l' -keys=$k -scale=$s; done; done
PASS
BenchmarkGoMapSpeed	   20000	     72026 ns/op
BenchmarkProbingMapSpeed	   20000	     78436 ns/op
BenchmarkGoMapMem	       0	         0 ns/op
--- BENCH: BenchmarkGoMapMem
	probing_bench_test.go:143: go-map: 0.04MB
BenchmarkProbingMapMem	       0	         0 ns/op
--- BENCH: BenchmarkProbingMapMem
	probing_bench_test.go:127: probing: 1000 keys; 1500 buckets; load factor 0.67
	probing_bench_test.go:143: probing: 0.02MB
ok  	github.com/kho/probing	4.561s
PASS
BenchmarkGoMapSpeed	   10000	    115136 ns/op
BenchmarkProbingMapSpeed	   10000	    110323 ns/op
BenchmarkGoMapMem	       0	         0 ns/op
--- BENCH: BenchmarkGoMapMem
	probing_bench_test.go:143: go-map: 36.58MB
BenchmarkProbingMapMem	       0	         0 ns/op
--- BENCH: BenchmarkProbingMapMem
	probing_bench_test.go:127: probing: 1000000 keys; 1500000 buckets; load factor 0.67
	probing_bench_test.go:143: probing: 22.89MB
ok  	github.com/kho/probing	5.609s
PASS
BenchmarkGoMapSpeed	   10000	    137659 ns/op
BenchmarkProbingMapSpeed	   10000	    131698 ns/op
BenchmarkGoMapMem	       0	         0 ns/op
--- BENCH: BenchmarkGoMapMem
	probing_bench_test.go:143: go-map: 303.56MB
BenchmarkProbingMapMem	       0	         0 ns/op
--- BENCH: BenchmarkProbingMapMem
	probing_bench_test.go:127: probing: 10000000 keys; 15000000 buckets; load factor 0.67
	probing_bench_test.go:143: probing: 228.88MB
ok  	github.com/kho/probing	38.949s
PASS
BenchmarkGoMapSpeed	   20000	     73367 ns/op
BenchmarkProbingMapSpeed	   20000	     64460 ns/op
BenchmarkGoMapMem	       0	         0 ns/op
--- BENCH: BenchmarkGoMapMem
	probing_bench_test.go:143: go-map: 0.04MB
BenchmarkProbingMapMem	       0	         0 ns/op
--- BENCH: BenchmarkProbingMapMem
	probing_bench_test.go:127: probing: 1000 keys; 2000 buckets; load factor 0.50
	probing_bench_test.go:143: probing: 0.03MB
ok  	github.com/kho/probing	4.167s
PASS
BenchmarkGoMapSpeed	   10000	    115724 ns/op
BenchmarkProbingMapSpeed	   20000	     89159 ns/op
BenchmarkGoMapMem	       0	         0 ns/op
--- BENCH: BenchmarkGoMapMem
	probing_bench_test.go:143: go-map: 36.59MB
BenchmarkProbingMapMem	       0	         0 ns/op
--- BENCH: BenchmarkProbingMapMem
	probing_bench_test.go:127: probing: 1000000 keys; 2000000 buckets; load factor 0.50
	probing_bench_test.go:143: probing: 30.52MB
ok  	github.com/kho/probing	7.349s
PASS
BenchmarkGoMapSpeed	   10000	    137724 ns/op
BenchmarkProbingMapSpeed	   10000	    100151 ns/op
BenchmarkGoMapMem	       0	         0 ns/op
--- BENCH: BenchmarkGoMapMem
	probing_bench_test.go:143: go-map: 303.51MB
BenchmarkProbingMapMem	       0	         0 ns/op
--- BENCH: BenchmarkProbingMapMem
	probing_bench_test.go:127: probing: 10000000 keys; 20000000 buckets; load factor 0.50
	probing_bench_test.go:143: probing: 305.18MB
ok  	github.com/kho/probing	37.895s
```
