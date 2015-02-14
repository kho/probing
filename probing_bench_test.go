package probing

import (
	"flag"
	"math/rand"
	"runtime"
	"testing"
)

var (
	numKeys       = flag.Int("keys", 1e6, "number of unique keys to insert")
	numInQueries  = flag.Int("in", 1e3, "number of in-map queries")
	numOutQueries = flag.Int("out", 1e3, "number of out-of-map queries")
	numReplicas   = flag.Int("replicas", 1, "number of replicas to use in memory benchmark")
	scale         = flag.Float64("scale", 1.5, "scale to compute init bucket size")
	load          = flag.Float64("load", 0, "load factor ( <= 0 or >= 1 means default)")
	seed          = flag.Int64("seed", 0, "seed to random number generator")
)

var (
	keys, inQs, outQs []__Key
	filled            bool
)

func fillData() {
	if filled {
		return
	}
	seen := map[__Key]bool{}
	rng := rand.New(rand.NewSource(*seed))
	for len(keys) < *numKeys {
		k := __Key(rng.Int())
		if k == __KEY_NIL {
			continue
		}
		if !seen[k] {
			keys = append(keys, k)
			seen[k] = true
		}
	}
	for len(inQs) < *numInQueries {
		k := keys[rng.Intn(*numKeys)]
		inQs = append(inQs, k)
	}
	for len(outQs) < *numOutQueries {
		k := __Key(rng.Int())
		if k == __KEY_NIL {
			continue
		}
		if !seen[k] {
			outQs = append(outQs, k)
		}
	}
	filled = true
}

func testMap() *__Map {
	m := New__Map(int(*scale*float64(*numKeys)), *load)
	for _, k := range keys {
		*m.FindOrInsert(k) = __Value(k + 1)
	}
	return m
}

func BenchmarkGoMapSpeed(b *testing.B) {
	fillData()
	m := make(map[__Key]__Value)
	for _, k := range keys {
		m[k] = __Value(k + 1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, k := range inQs {
			_ = m[k]
		}
		for _, k := range outQs {
			_ = m[k]
		}
	}
}

func BenchmarkProbingMapSpeed(b *testing.B) {
	fillData()
	m := testMap()
	for _, k := range keys {
		*m.FindOrInsert(k) = __Value(k + 1)
	}
	// numLookUps = 0
	// numCollisions = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, k := range inQs {
			_ = m.Find(k)
		}
		for _, k := range outQs {
			_ = m.Find(k)
		}
	}
	// b.Logf("%.2f hashes per look-up", float64(numLookUps+numCollisions)/float64(numLookUps))
}

func BenchmarkGoMapMem(b *testing.B) {
	fillData()

	ms := make([]map[__Key]__Value, *numReplicas)
	measureMem("go-map", func() {
		for i := range ms {
			m := make(map[__Key]__Value)
			for _, k := range keys {
				m[k] = __Value(k + 1)
			}
			ms[i] = m
		}
	}, b)

	b.SkipNow()
}

func BenchmarkProbingMapMem(b *testing.B) {
	fillData()

	ms := make([]*__Map, *numReplicas)
	measureMem("probing", func() {
		for i := range ms {
			m := testMap()
			if i == 0 {
				b.Logf("probing: %d keys; %d buckets; load factor %.2f", m.Size(), len(m.buckets), float64(m.Size())/float64(len(m.buckets)))
			}
			ms[i] = m
		}
	}, b)

	b.SkipNow()
}

func measureMem(name string, action func(), b *testing.B) {
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	action()
	runtime.GC()
	runtime.ReadMemStats(&after)
	b.Logf("%s: %.2fMB", name, float64(after.Alloc-before.Alloc)/float64(1<<20))
}
