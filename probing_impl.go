package probing

import (
	"bytes"
	"encoding/gob"
)

type Entry struct {
	Key   Key
	Value Value
}

type Map struct {
	buckets               buckets
	numEntries, threshold int
}

func NewMap(initNumBuckets int, maxUsed float64) *Map {
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
	return &Map{initBuckets(initNumBuckets), 0, threshold}
}

func (m *Map) Size() int {
	return m.numEntries
}

func (m *Map) Find(k Key) *Value {
	return m.buckets.Find(k)
}

func (m *Map) ConstFind(k Key) (Value, bool) {
	return m.buckets.ConstFind(k)
}

func (m *Map) FindOrInsert(k Key) *Value {
	v := m.Find(k)
	if v != nil {
		return v
	}
	return m.insert(k)
}

func (m *Map) Range() chan Entry {
	out := make(chan Entry)
	go func() {
		for _, e := range m.buckets {
			if e.Key != KEY_NIL {
				out <- e
			}
		}
		close(out)
	}()
	return out
}

func (m *Map) MarshalBinary() (data []byte, err error) {
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

func (m *Map) UnmarshalBinary(data []byte) (err error) {
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

func (m *Map) insert(k Key) *Value {
	if m.numEntries >= m.threshold {
		m.double()
	}
	ei := m.buckets.nextAvailable(k)
	*ei = Entry{Key: k}
	m.numEntries++
	return &ei.Value
}

func (m *Map) double() {
	buckets := initBuckets(len(m.buckets) * 2)
	for _, e := range m.buckets {
		k := e.Key
		if !equal(k, KEY_NIL) {
			dst := buckets.nextAvailable(k)
			*dst = e
		}
	}
	m.buckets = buckets
	m.threshold *= 2
}

type buckets []Entry

func initBuckets(n int) buckets {
	s := make(buckets, n)
	for i := range s {
		s[i].Key = KEY_NIL
	}
	return s
}

func (b buckets) Find(k Key) (v *Value) {
	i := b.start(k)
	for {
		// Maybe switch to range to trade 1 bound check for 1 copy?
		ei := &b[i]
		ki := ei.Key
		if equal(ki, k) {
			return &ei.Value
		}
		if equal(ki, KEY_NIL) {
			return nil
		}
		i++
		if i == len(b) {
			i = 0
		}
	}
}

func (b buckets) ConstFind(k Key) (Value, bool) {
	i := b.start(k)
	for _, e := range b[i:] {
		ki := e.Key
		if equal(ki, k) {
			return e.Value, true
		}
		if equal(ki, KEY_NIL) {
			var v Value
			return v, false
		}
	}
	for _, e := range b[:i] {
		ki := e.Key
		if equal(ki, k) {
			return e.Value, true
		}
		if equal(ki, KEY_NIL) {
			var v Value
			return v, false
		}
	}
	panic("impossible")
}

func (b buckets) start(k Key) int {
	return int(hash(k) % uint(len(b)))
}

func (b buckets) nextAvailable(k Key) *Entry {
	i := b.start(k)
	for {
		ei := &b[i]
		if equal(ei.Key, KEY_NIL) {
			return ei
		}
		i++
		if i == len(b) {
			i = 0
		}
	}
}
