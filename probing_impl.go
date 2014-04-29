package probing

import ()

type entry struct {
	Key   Key
	Value Value
}

type Map struct {
	buckets               buckets
	numEntries, threshold int
}

func NewMap(initSize int) *Map {
	// numBuckets = max(2, initSize*ALLOC, initSize+1)
	numBuckets := int(float64(initSize) * _ALLOC_MULTIPLIER)
	if numBuckets < 2 {
		numBuckets = 2
	}
	if numBuckets < initSize+1 {
		numBuckets = initSize + 1
	}
	// threshold = min(max(1, initSize*THRESH), numBuckets-1)
	threshold := int(float64(initSize) * _THRESHOLD_MULTIPLIER)
	if threshold < 1 {
		threshold = 1
	}
	if threshold > numBuckets-1 {
		threshold = numBuckets - 1
	}
	return &Map{initBuckets(numBuckets), 0, threshold}
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

func (m *Map) insert(k Key) *Value {
	if m.numEntries >= m.threshold {
		m.double()
	}
	ei := m.buckets.nextAvailable(k)
	*ei = entry{Key: k}
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

type buckets []entry

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

func (b buckets) nextAvailable(k Key) *entry {
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
