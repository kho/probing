package probing

import (
	"testing"
)

func TestMapBasic(t *testing.T) {
	m := NewMap(0)
	// Insert values.
	for i := 0; i < 1000; i++ {
		*m.FindOrInsert(Key(i)) = Value(i + 1)
	}
	if n := m.Size(); n != 1000 {
		t.Errorf("incorrect size %d after 1000 insertions", n)
	}
	for i := 0; i < 1000; i++ {
		*m.FindOrInsert(Key(i)) = Value(i + 1)
	}
	if n := m.Size(); n != 1000 {
		t.Errorf("incorrect size %d after 2000 insertions", n)
	}
	// Find existent.
	for i := 0; i < 1000; i++ {
		v := m.Find(Key(i))
		if v == nil {
			t.Errorf("Find(%d) should not be nil", i)
		} else if *v != Value(i+1) {
			t.Errorf("expect Find(%d) = %d; got %d", i, i+1, *v)
		}
		vv, ok := m.ConstFind(Key(i))
		if !ok {
			t.Errorf("ConstFind(%d) should be ok ", i)
		} else if vv != Value(i+1) {
			t.Errorf("expect ConstFind(%d) = %d; got %d", i, i+1, vv)
		}
	}
	// Find non-existent.
	for i := 1000; i < 2000; i++ {
		v := m.Find(Key(i))
		if v != nil {
			t.Errorf("Find(%d) should be nil", i)
		}
		_, ok := m.ConstFind(Key(i))
		if ok {
			t.Errorf("ConstFind(%d) should be not ok ", i)
		}
	}
}

func Test_bucketsRollOver(t *testing.T) {
	b := initBuckets(4)
	b[2] = entry{0, 1}
	b[3] = entry{0, 1}
	*b.nextAvailable(1) = entry{1, 2}
	if _, ok := b.ConstFind(1); !ok {
		t.Error("cannot find Key(1)")
	}
	for i := 2; i < 100; i++ {
		if b.Find(Key(i)) != nil {
			t.Errorf("Find(%d) should be nil", i)
		}
		if _, ok := b.ConstFind(Key(i)); ok {
			t.Errorf("ConstFind(%d) should not be ok", i)
		}
		if n := b.nextAvailable(Key(i)); n.Value != 0 {
			t.Errorf("expect nextAvailable(%d).Value = 0; got %d", i, n.Value)
		}
	}
}
