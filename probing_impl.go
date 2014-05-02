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
	buckets               __Buckets
	numEntries, threshold int
}

func New__Map(initNumBuckets int, maxUsed float64) *__Map {
	if initNumBuckets == 0 {
		initNumBuckets = 4
	} else if initNumBuckets < 2 {
		initNumBuckets = 2
	}
	if maxUsed <= 0 || maxUsed >= 1 {
		maxUsed = 0.8
	}
	// threshold = min(max(1, initNumBuckets * maxUsed), initNumBuckets-1)
	threshold := int(float64(initNumBuckets) * maxUsed)
	if threshold < 1 {
		threshold = 1
	}
	if threshold > initNumBuckets-1 {
		threshold = initNumBuckets - 1
	}
	return &__Map{__InitBuckets(initNumBuckets), 0, threshold}
}

func (m *__Map) Size() int {
	return m.numEntries
}

func (m *__Map) Find(k __Key) *__Value {
	return m.buckets.Find(k)
}

func (m *__Map) FindOrInsert(k __Key) *__Value {
	e := m.buckets.FindEntry(k)
	if e.Key != __KEY_NIL {
		return &e.Value
	}
	// Need to insert.
	if m.numEntries >= m.threshold {
		m.Resize(len(m.buckets) * 2)
		e = m.buckets.nextAvailable(k)
	}
	*e = __Entry{Key: k}
	m.numEntries++
	return &e.Value
}

func (m *__Map) Resize(numBuckets int) {
	if numBuckets < m.numEntries+1 {
		numBuckets = m.numEntries + 1
	}
	buckets := __InitBuckets(numBuckets)
	for _, e := range m.buckets {
		k := e.Key
		if !__Equal(k, __KEY_NIL) {
			dst := buckets.nextAvailable(k)
			*dst = e
		}
	}
	oldNumBuckets := len(m.buckets)
	m.buckets = buckets
	m.threshold = m.threshold * numBuckets / oldNumBuckets
	if m.threshold < m.numEntries {
		m.threshold = m.numEntries
	}
}

func (m *__Map) Range() chan __Entry {
	return m.buckets.Range()
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

type __Buckets []__Entry

func __InitBuckets(n int) __Buckets {
	s := make(__Buckets, n)
	for i := range s {
		s[i].Key = __KEY_NIL
	}
	return s
}

func (b __Buckets) Size() (n int) {
	for _, e := range b {
		if e.Key != __KEY_NIL {
			n++
		}
	}
	return
}

// var numLookUps, numCollisions int

func (b __Buckets) Find(k __Key) (v *__Value) {
	// numLookUps++
	i := b.start(k)
	for {
		// Maybe switch to range to trade 1 bound check for 1 copy?
		ei := &b[i]
		ki := ei.Key
		if __Equal(ki, k) {
			return &ei.Value
		}
		if __Equal(ki, __KEY_NIL) {
			return nil
		}
		// numCollisions++
		i++
		if i == len(b) {
			i = 0
		}
	}
}

func (b __Buckets) FindEntry(k __Key) *__Entry {
	i := b.start(k)
	for {
		ei := &b[i]
		ki := ei.Key
		if __Equal(ki, k) || __Equal(ki, __KEY_NIL) {
			return ei
		}
		i++
		if i == len(b) {
			i = 0
		}
	}
}

func (b __Buckets) Range() chan __Entry {
	ch := make(chan __Entry)
	go func() {
		for _, e := range b {
			if e.Key != __KEY_NIL {
				ch <- e
			}
		}
		close(ch)
	}()
	return ch
}

func (b __Buckets) start(k __Key) int {
	return int(__Hash(k) % uint(len(b)))
}

func (b __Buckets) nextAvailable(k __Key) *__Entry {
	i := b.start(k)
	for {
		ei := &b[i]
		if __Equal(ei.Key, __KEY_NIL) {
			return ei
		}
		i++
		if i == len(b) {
			i = 0
		}
	}
}
