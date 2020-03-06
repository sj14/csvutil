package csvutil

import (
	"testing"
)

func TestEquals(t *testing.T) {
	testCases := []struct {
		description string
		datasetA    [][]string
		datasetB    [][]string
		expected    bool
	}{
		{
			description: "empty",
			datasetA:    [][]string{[]string{}},
			datasetB:    [][]string{[]string{}},
			expected:    true,
		},
		{
			description: "happy",
			datasetA:    [][]string{[]string{"A", "B", "C"}, []string{"1", "2", "3"}},
			datasetB:    [][]string{[]string{"A", "B", "C"}, []string{"1", "2", "3"}},
			expected:    true,
		},
		{
			description: "fail/value",
			datasetA:    [][]string{[]string{"A", "B", "C"}, []string{"1", "99", "3"}},
			datasetB:    [][]string{[]string{"A", "B", "C"}, []string{"1", "2", "3"}},
			expected:    false,
		},
		{
			description: "fail/row/len",
			datasetA:    [][]string{[]string{"A", "B", "C"}, []string{"1", "2", "3"}},
			datasetB:    [][]string{[]string{"A", "B", "C"}},
			expected:    false,
		},
		{
			description: "fail/col/len",
			datasetA:    [][]string{[]string{"A", "B", "C"}, []string{"1", "3", "3"}},
			datasetB:    [][]string{[]string{"A", "B", "C"}, []string{"1", "2"}},
			expected:    false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			if Equals(tt.datasetA, tt.datasetB) != tt.expected {
				t.Fail()
			}
		})
	}
}

func TestNew(t *testing.T) {
	raw := [][]string{
		[]string{"Row 1 Col 1", "Row 1 Col 2"},
		[]string{"Row 2 Col 1", "Row 2 Col 2"},
	}
	ds := New(raw)

	if !Equals(ds.Raw(), raw) {
		t.Fail()
	}
}
