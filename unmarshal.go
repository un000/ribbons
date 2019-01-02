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

	// make Add efficient
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	for i := 0; i < len(list); i++ {
		u.Add(list[i])
	}

	return nil
}
