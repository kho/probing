package probing

type __Key uint64
type __Value uint64

const __KEY_NIL = ^__Key(0)

func __hash(k __Key) uint {
	// https://code.google.com/p/fast-hash
	h := uint64(k)
	h ^= h >> 23
	h *= 0x2127599bf4325c37
	h ^= h >> 47
	return uint(h)
}

func __equal(a, b __Key) bool {
	return a == b
}
