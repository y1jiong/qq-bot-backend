package utility

import (
	"testing"
)

func TestFindNaturalBreakPoint(t *testing.T) {
	tests := []struct {
		name      string
		runes     []rune
		maxLen    int
		maxOffset int
		expected  int
	}{
		{
			name:      "Shorter than maxLen",
			runes:     []rune("Hello"),
			maxLen:    10,
			maxOffset: 5,
			expected:  5, // The length of the string itself is the breakpoint
		},
		{
			name:      "No punctuation or space before maxLen",
			runes:     []rune("HelloWorld"),
			maxLen:    10,
			maxOffset: 5,
			expected:  10, // No natural break, return maxLen
		},
		{
			name:      "Natural break with space",
			runes:     []rune("Hello World"),
			maxLen:    10,
			maxOffset: 5,
			expected:  6, // "Hello " is a natural break
		},
		{
			name:      "Natural break with punctuation",
			runes:     []rune("Hello,world!"),
			maxLen:    10,
			maxOffset: 5,
			expected:  6, // "Hello," is a natural break
		},
		{
			name:      "Edge case maxOffset equals maxLen",
			runes:     []rune("This is a test!"),
			maxLen:    15,
			maxOffset: 15,
			expected:  15, // Found "!"
		},
		{
			name:      "Empty string",
			runes:     []rune(""),
			maxLen:    10,
			maxOffset: 5,
			expected:  0, // Empty string should return 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindNaturalBreakPoint(tt.runes, tt.maxLen, tt.maxOffset)
			if result != tt.expected {
				t.Errorf("For input %q, expected %d, but got %d", string(tt.runes), tt.expected, result)
			}
		})
	}
}
