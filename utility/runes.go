package utility

import "unicode"

func FindNaturalBreakPoint(runes []rune, maxLen, maxOffset int) int {
	if len(runes) <= maxLen {
		return len(runes)
	}
	if maxOffset > maxLen {
		maxOffset = maxLen
	}

	for i := maxLen - 1; i >= maxLen-maxOffset; i-- {
		r := runes[i]
		if unicode.IsPunct(r) || unicode.IsSpace(r) {
			return i + 1
		}
	}
	return maxLen
}
