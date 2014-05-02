package probing

import (
	"testing"
)

func Test__MapBasic(t *testing.T) {
	m := New__Map(0, 0)
	// Insert values.
	for i := 0; i < 1000; i++ {
		*m.FindOrInsert(__Key(i)) = __Value(i + 1)
	}
	if n := m.Size(); n != 1000 {
		t.Errorf("incorrect size %d after 1000 insertions", n)
	}
	for i := 0; i < 1000; i++ {
		*m.FindOrInsert(__Key(i)) = __Value(i + 1)
	}
	if n := m.Size(); n != 1000 {
		t.Errorf("incorrect size %d after 2000 insertions", n)
	}
	// Find existent.
	for i := 0; i < 1000; i++ {
		v := m.Find(__Key(i))
		if v == nil {
			t.Errorf("Find(%d) should not be nil", i)
		} else if *v != __Value(i+1) {
			t.Errorf("expect Find(%d) = %d; got %d", i, i+1, *v)
		}
	}
	// Find non-existent.
	for i := 1000; i < 2000; i++ {
		v := m.Find(__Key(i))
		if v != nil {
			t.Errorf("Find(%d) should be nil", i)
		}
	}
	// Ranging over entries.
	for e := range m.Range() {
		if e.Key == __KEY_NIL || e.Key < 0 || e.Key >= 1000 {
			t.Errorf("invalid key: %d", e.Key)
		}
		if __Value(e.Key+1) != e.Value {
			t.Errorf("invalid entry: %+v", e)
		}
	}
}

func Test___bucketsRollOver(t *testing.T) {
	b := __InitBuckets(4)
	b[2] = __Entry{0, 1}
	b[3] = __Entry{0, 1}
	*b.nextAvailable(1) = __Entry{1, 2}
	if b.Find(1) == nil {
		t.Error("cannot find __Key(1)")
	}
	for i := 2; i < 100; i++ {
		if b.Find(__Key(i)) != nil {
			t.Errorf("Find(%d) should be nil", i)
		}
		if b.Find(__Key(i)) != nil {
			t.Errorf("Find(%d) should be nil", i)
		}
		if n := b.nextAvailable(__Key(i)); n.Value != 0 {
			t.Errorf("expect nextAvailable(%d).Value = 0; got %d", i, n.Value)
		}
	}
}
