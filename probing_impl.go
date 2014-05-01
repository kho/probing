package probing

import (
	"bytes"
	"encoding/gob"
)

type __Entry struct {
	Key   __Key
	Value __Value
}

type __Map struct {
	buckets               __buckets
	numEntries, threshold int
}

func __NewMap(initNumBuckets int, maxUsed float64) *__Map {
	if initNumBuckets < 2 {
		initNumBuckets = 2
	}
	if maxUsed <= 0 || maxUsed >= 1 {
		maxUsed = 0.8
	}
	// threshold = min(max(1, (initNumBuckets-1) * maxUsed), initNumBuckets-1)
	threshold := int(float64(initNumBuckets-1) * maxUsed)
	if threshold < 1 {
		threshold = 1
	}
	if threshold > initNumBuckets-1 {
		threshold = initNumBuckets - 1
	}
	return &__Map{__initBuckets(initNumBuckets), 0, threshold}
}

func (m *__Map) Size() int {
	return m.numEntries
}

func (m *__Map) Find(k __Key) *__Value {
	return m.buckets.Find(k)
}

func (m *__Map) FindOrInsert(k __Key) *__Value {
	v := m.Find(k)
	if v != nil {
		return v
	}
	return m.insert(k)
}

func (m *__Map) Range() chan __Entry {
	out := make(chan __Entry)
	go func() {
		for _, e := range m.buckets {
			if e.Key != __KEY_NIL {
				out <- e
			}
		}
		close(out)
	}()
	return out
}

func (m *__Map) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(m.buckets); err != nil {
		return
	}
	if err = enc.Encode(m.numEntries); err != nil {
		return
	}
	if err = enc.Encode(m.threshold); err != nil {
		return
	}
	return buf.Bytes(), nil
}

func (m *__Map) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err = dec.Decode(&m.buckets); err != nil {
		return
	}
	if err = dec.Decode(&m.numEntries); err != nil {
		return
	}
	if err = dec.Decode(&m.threshold); err != nil {
		return
	}
	return nil
}

func (m *__Map) insert(k __Key) *__Value {
	if m.numEntries >= m.threshold {
		m.double()
	}
	ei := m.buckets.nextAvailable(k)
	*ei = __Entry{Key: k}
	m.numEntries++
	return &ei.Value
}

func (m *__Map) double() {
	buckets := __initBuckets(len(m.buckets) * 2)
	for _, e := range m.buckets {
		k := e.Key
		if !__equal(k, __KEY_NIL) {
			dst := buckets.nextAvailable(k)
			*dst = e
		}
	}
	m.buckets = buckets
	m.threshold *= 2
}

type __buckets []__Entry

func __initBuckets(n int) __buckets {
	s := make(__buckets, n)
	for i := range s {
		s[i].Key = __KEY_NIL
	}
	return s
}

// var numLookUps, numCollisions int

func (b __buckets) Find(k __Key) (v *__Value) {
	// numLookUps++
	i := b.start(k)
	for {
		// Maybe switch to range to trade 1 bound check for 1 copy?
		ei := &b[i]
		ki := ei.Key
		if __equal(ki, k) {
			return &ei.Value
		}
		if __equal(ki, __KEY_NIL) {
			return nil
		}
		// numCollisions++
		i++
		if i == len(b) {
			i = 0
		}
	}
}

func (b __buckets) start(k __Key) int {
	return int(__hash(k) % uint(len(b)))
}

func (b __buckets) nextAvailable(k __Key) *__Entry {
	i := b.start(k)
	for {
		ei := &b[i]
		if __equal(ei.Key, __KEY_NIL) {
			return ei
		}
		i++
		if i == len(b) {
			i = 0
		}
	}
}
