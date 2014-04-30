package probing

type __Key uint64
type __Value uint64

const __KEY_NIL = ^__Key(0)

func __hash(k __Key) uint {
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

func __equal(a, b __Key) bool {
	return a == b
}
