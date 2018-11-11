package ribbons

import (
	"encoding/json"
	"sort"
)

type UINT64Set struct {
	// min is the value of the first element
	min  uint64
	set  []uint64
	size int

	initialized bool
}

func (u *UINT64Set) Has(key uint64) bool {
	switch {
	case !u.initialized || len(u.set) == 0:
		return false
	case key < u.min:
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
		u.set = append(u.set, setBit(0, 0))
		u.size++
		u.initialized = true
	case u.initialized && key == u.min:
		return u
	case key > u.min:
		if u.Has(key) {
			return u
		}

		bucket := u.getBucketNumber(key)
		if bucket >= len(u.set) {
			u.set = append(u.set, make([]uint64, bucket-len(u.set)+1)...)
		}

		u.set[bucket] = setBit(u.set[bucket], u.getPositionInsideBucket(key))
		u.size++
		u.initialized = true
	case key < u.min:
		newSet := UINT64Set{set: make([]uint64, 0, len(u.set))}
		newSet.Add(key)
		for _, prev := range u.List() {
			newSet.Add(prev)
		}
		*u = newSet
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

func (u *UINT64Set) Sum(u2 *UINT64Set) {
	second := u2.List()
	for i := range second {
		u.Add(second[i])
	}
}

func (u *UINT64Set) Mul(u2 *UINT64Set) {
	for i, l := uint64(0), uint64(len(u.set)); i < l; i++ {
		if u.set[i] == 0 {
			continue
		}

		for _, v := range extractToggledBits(u.set[i], u.min+i*64) {
			if !u2.Has(v) {
				u.Delete(v)
			}
		}
	}
}

func (u *UINT64Set) List() []uint64 {
	res := make([]uint64, 0, u.size)
	for i, l := uint64(0), uint64(len(u.set)); i < l; i++ {
		if u.set[i] == 0 {
			continue
		}

		values := extractToggledBits(u.set[i], u.min+i*64)
		res = append(res, values...)
	}

	return res
}

func (u *UINT64Set) UnmarshalJSON(b []byte) error {
	list := make([]uint64, 0)
	err := json.Unmarshal(b, &list)
	if err != nil {
		return err
	}

	// make Add efficient
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	for i := 0; i < len(list); i++ {
		u.Add(list[i])
	}

	return nil
}

func (u *UINT64Set) getPositionInsideBucket(key uint64) uint64 {
	return (key - u.min) % 64
}

func (u *UINT64Set) getBucketNumber(key uint64) int {
	return int((key - u.min) / 64)
}
