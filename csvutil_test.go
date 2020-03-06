package csvutil

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
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
	testCases := []struct {
		description string
		create      [][]string
	}{
		{
			description: "nil",
			create:      nil,
		},
		{
			description: "empty",
			create:      [][]string{[]string{}},
		},
		{
			description: "normal",
			create: [][]string{
				[]string{"Row 1 Col 1", "Row 1 Col 2"},
				[]string{"Row 2 Col 1", "Row 2 Col 2"},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			ds := New(tt.create)
			if !Equals(ds.Raw(), tt.create) {
				t.Fail()
			}
		})
	}
}

func TestAddColumn(t *testing.T) {
	testCases := []struct {
		description string
		init        [][]string
		add         []string
		index       int
		want        [][]string
	}{
		{
			description: "add to empty",
			init:        [][]string{},
			add:         []string{"Row 1 Col 1", "Row 2 Col 1"},
			index:       0,
			want: [][]string{
				[]string{"Row 1 Col 1", "Row 2 Col 1"},
			},
		},
		{
			description: "add at beginning",
			init: [][]string{
				[]string{"Row 1 Col 1", "Row 1 Col 2"},
				[]string{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:   []string{"Row 1 Col 0", "Row 2 Col 0"},
			index: 0,
			want: [][]string{
				[]string{"Row 1 Col 0", "Row 1 Col 1", "Row 1 Col 2"},
				[]string{"Row 2 Col 0", "Row 2 Col 1", "Row 2 Col 2"},
			},
		},
		{
			description: "add in-between",
			init: [][]string{
				[]string{"Row 1 Col 1", "Row 1 Col 2"},
				[]string{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:   []string{"Row 1 Col 1.5", "Row 2 Col 1.5"},
			index: 1,
			want: [][]string{
				[]string{"Row 1 Col 1", "Row 1 Col 1.5", "Row 1 Col 2"},
				[]string{"Row 2 Col 1", "Row 2 Col 1.5", "Row 2 Col 2"},
			},
		},
		{
			description: "add at end/negative",
			init: [][]string{
				[]string{"Row 1 Col 1", "Row 1 Col 2"},
				[]string{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:   []string{"Row 1 Col 3", "Row 2 Col 3"},
			index: 2,
			want: [][]string{
				[]string{"Row 1 Col 1", "Row 1 Col 2", "Row 1 Col 3"},
				[]string{"Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3"},
			},
		},
		{
			description: "add at end/negative",
			init: [][]string{
				[]string{"Row 1 Col 1", "Row 1 Col 2"},
				[]string{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:   []string{"Row 1 Col 3", "Row 2 Col 3"},
			index: -1,
			want: [][]string{
				[]string{"Row 1 Col 1", "Row 1 Col 2", "Row 1 Col 3"},
				[]string{"Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3"},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			ds := New(tt.init)
			err := ds.AddColumn(tt.add, tt.index)
			require.NoError(t, err)

			log.Printf("want: %v\ngot: %v\n", tt.want, ds.Raw())
			require.True(t, Equals(tt.want, ds.Raw()))
		})
	}

}
