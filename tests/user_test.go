package cache_test

import (
	"reflect"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
)

func TestRemoveRepeatedElementsInList(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "No duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "All duplicates",
			input:    []string{"a", "a", "a"},
			expected: []string{"a"},
		},
		{
			name:     "Some duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Empty list",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Single element",
			input:    []string{"a"},
			expected: []string{"a"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cache.RemoveRepeatedElementsInList(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
