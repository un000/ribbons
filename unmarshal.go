package ribbons

import (
	"encoding/json"
	"sort"
)

func (u *UINT64Set) UnmarshalJSON(b []byte) error {
	list := make([]uint64, 0)
	err := json.Unmarshal(b, &list)
	if err != nil {
		return err
	}

	if list == nil {
		u.initialized = false
		return nil
	}

	// make Add efficient
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	for i := 0; i < len(list); i++ {
		u.Add(list[i])
	}

	u.initialized = true

	return nil
}
