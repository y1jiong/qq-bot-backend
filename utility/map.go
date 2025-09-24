package utility

import (
	"sort"
	"strings"
)

func SortArrayReverseFromMapKey(m map[string]any) (arr []string) {
	arr = make([]string, 0, len(m))
	for k := range m {
		arr = append(arr, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(arr)))
	return
}

func MatchAllKeywords(str string, m map[string]any) (eureka bool, hit, mValue string) {
	for _, k := range SortArrayReverseFromMapKey(m) {
		const exactSuffix = "$"
		if strings.HasSuffix(k, exactSuffix) {
			kExact := k[:len(k)-len(exactSuffix)]
			if str == kExact {
				eureka = true
				hit = kExact
				if vv, ok := m[k].(string); ok {
					mValue = vv
				}
				break
			}
			continue
		}

		fields := strings.Fields(k)
		var notContains bool
		for _, kk := range fields {
			if !strings.Contains(str, kk) {
				notContains = true
				break
			}
		}
		if notContains {
			continue
		}

		eureka = true
		hit = k
		if vv, ok := m[k].(string); ok {
			mValue = vv
		}
		if strings.HasPrefix(str, k) {
			break
		}
	}
	return
}
