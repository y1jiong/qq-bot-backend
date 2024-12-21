package utility

import "sort"

func ReverseSortedArrayFromMapKey(m map[string]any) (arr []string) {
	arr = make([]string, 0, len(m))
	for k := range m {
		arr = append(arr, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(arr)))
	return
}
