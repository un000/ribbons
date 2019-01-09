package ribbons

import (
	"encoding/json"
	"fmt"
)

type bucket = uint64

const bucketSize = 64

type UINT64Set struct {
	// min is the value of the first element
	min  uint64
	set  []bucket
	size int

	initialized bool
}

func New() UINT64Set {
	return UINT64Set{}
}

func NewFromJSON(j []byte) (UINT64Set, error) {
	set := New()
	err := json.Unmarshal(j, &set)
	if err != nil {
		return set, err
	}

	return set, nil
}

func (u *UINT64Set) Has(key uint64) bool {
	if !u.initialized || len(u.set) == 0 || key < u.min {
		return false
	}

	bucket := u.getBucketNumber(key)
	if bucket >= len(u.set) {
		return false
	}

	return hasBit(u.set[bucket], u.getPositionInsideBucket(key))
}

func (u *UINT64Set) Add(key uint64) *UINT64Set {
	switch {
	case len(u.set) == 0:
		u.min = key
		u.set = append(u.set, 1)
		u.size++
		u.initialized = true
	case key >= u.min:
		if u.Has(key) {
			return u
		}

		bucketN := u.getBucketNumber(key)
		if bucketN >= len(u.set) {
			u.set = append(u.set, make([]bucket, bucketN-len(u.set)+1)...)
		}

		u.set[bucketN] = setBit(u.set[bucketN], u.getPositionInsideBucket(key))
		u.size++
		u.initialized = true
	case key < u.min:
		newSet := UINT64Set{set: make([]bucket, 0, len(u.set))}
		newSet.Add(key)
		for _, prev := range u.List() {
			newSet.Add(prev)
		}
		*u = newSet
	default:
		panic(fmt.Sprintf("unexpected behaviour; key: %v; set: %+v", key, u))
	}

	return u
}

func (u *UINT64Set) Delete(key uint64) {
	if !u.initialized || len(u.set) == 0 {
		return
	}

	bucket := u.getBucketNumber(key)
	if bucket >= len(u.set) {
		return
	}

	last := u.set[bucket]
	deleted := delBit(u.set[bucket], u.getPositionInsideBucket(key))
	if last != deleted {
		u.set[bucket] = deleted
		u.size--
	}
}

func (u *UINT64Set) Len() int {
	return u.size
}

func (u *UINT64Set) Or(u2 *UINT64Set) {
	second := u2.List()
	for i := range second {
		u.Add(second[i])
	}
}

func (u *UINT64Set) And(u2 *UINT64Set) {
	for i, l := uint64(0), uint64(len(u.set)); i < l; i++ {
		for _, v := range extractToggledBits(bucketSize, u.set[i], u.min+i*bucketSize) {
			if !u2.Has(v) {
				u.Delete(v)
			}
		}
	}
}

func (u *UINT64Set) AndNot(u2 *UINT64Set) {
	for i, l := uint64(0), uint64(len(u.set)); i < l; i++ {
		for _, v := range extractToggledBits(bucketSize, u.set[i], u.min+i*bucketSize) {
			if u2.Has(v) {
				u.Delete(v)
			}
		}
	}
}

func (u *UINT64Set) List() []uint64 {
	res := make([]uint64, 0, u.size)
	for i, l := uint64(0), uint64(len(u.set)); i < l; i++ {
		values := extractToggledBits(bucketSize, u.set[i], u.min+i*bucketSize)
		if len(values) > 0 {
			res = append(res, values...)
		}
	}

	return res
}

func (u *UINT64Set) Initialized() bool {
	return u.initialized
}

func (u *UINT64Set) getPositionInsideBucket(key uint64) uint64 {
	return (key - u.min) % bucketSize
}

func (u *UINT64Set) getBucketNumber(key uint64) int {
	return int((key - u.min) / bucketSize)
}
