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
			name:     "empty list",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "list with no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "list with duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := cache.RemoveRepeatedElementsInList(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("RemoveRepeatedElementsInList(%v) = %v; want %v", tc.input, actual, tc.expected)
			}
		})
	}
}
