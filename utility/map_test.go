package utility

import (
	"reflect"
	"testing"
)

func TestSortArrayReverseFromMapKey(t *testing.T) {
	tests := []struct {
		name    string
		m       map[string]any
		wantArr []string
	}{
		{
			name:    "Empty map",
			m:       map[string]any{},
			wantArr: []string{},
		},
		{
			name:    "Single element map",
			m:       map[string]any{"a": 1},
			wantArr: []string{"a"},
		},
		{
			name:    "Multiple elements map",
			m:       map[string]any{"c": 3, "a": 1, "b": 2},
			wantArr: []string{"c", "b", "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotArr := SortArrayReverseFromMapKey(tt.m); !reflect.DeepEqual(gotArr, tt.wantArr) {
				t.Errorf("SortArrayReverseFromMapKey() = %v, want %v", gotArr, tt.wantArr)
			}
		})
	}
}
