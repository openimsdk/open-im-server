package cache

import (
	"testing"
	"reflect"
)

func TestRemoveRepeatedElementsInList(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{},
			expected: []string{},
		},
		{
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			input:    []string{"a", "a", "b", "b", "c", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, testCase := range testCases {
		result := RemoveRepeatedElementsInList(testCase.input)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("RemoveRepeatedElementsInList(%v) = %v; want %v", testCase.input, result, testCase.expected)
		}
	}
}
