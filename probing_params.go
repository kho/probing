package probing

type Key uint64
type Value uint64

const (
	KEY_NIL               = ^Key(0)
	_THRESHOLD_MULTIPLIER = 0.618
)

func hash(k Key) uint {
	// hash64shift from
	// https://web.archive.org/web/20120903003157/http://www.cris.com/~Ttwang/tech/inthash.htm.
	r := uint64(k) // make sure >> is logical.
	r = (^r) + (r << 21)
	r = r ^ (r >> 24)
	r = (r + (r << 3)) + (r << 8)
	r = r ^ (r >> 14)
	r = (r + (r << 2)) + (r << 4)
	r = r ^ (r >> 28)
	r = r + (r << 31)
	return uint(r)
}

func equal(a, b Key) bool {
	return a == b
}
