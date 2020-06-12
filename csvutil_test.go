package csvutil

import (
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
			description: "nil",
			datasetA:    nil,
			datasetB:    nil,
			expected:    true,
		},
		{
			description: "empty/empty",
			datasetA:    [][]string{},
			datasetB:    [][]string{{}},
			expected:    false,
		},
		{
			description: "empty",
			datasetA:    [][]string{{}},
			datasetB:    [][]string{{}},
			expected:    true,
		},
		{
			description: "happy",
			datasetA:    [][]string{{"A", "B", "C"}, {"1", "2", "3"}},
			datasetB:    [][]string{{"A", "B", "C"}, {"1", "2", "3"}},
			expected:    true,
		},
		{
			description: "fail/value",
			datasetA:    [][]string{{"A", "B", "C"}, {"1", "99", "3"}},
			datasetB:    [][]string{{"A", "B", "C"}, {"1", "2", "3"}},
			expected:    false,
		},
		{
			description: "fail/row/len",
			datasetA:    [][]string{{"A", "B", "C"}, {"1", "2", "3"}},
			datasetB:    [][]string{{"A", "B", "C"}},
			expected:    false,
		},
		{
			description: "fail/col/len",
			datasetA:    [][]string{{"A", "B", "C"}, {"1", "3", "3"}},
			datasetB:    [][]string{{"A", "B", "C"}, {"1", "2"}},
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
			create:      [][]string{{}},
		},
		{
			description: "normal",
			create: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
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

func TestAddCol(t *testing.T) {
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
				{"Row 1 Col 1", "Row 2 Col 1"},
			},
		},
		{
			description: "add at beginning",
			init: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:   []string{"Row 1 Col 0", "Row 2 Col 0"},
			index: 0,
			want: [][]string{
				{"Row 1 Col 0", "Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 0", "Row 2 Col 1", "Row 2 Col 2"},
			},
		},
		{
			description: "add in-between",
			init: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:   []string{"Row 1 Col 1.5", "Row 2 Col 1.5"},
			index: 1,
			want: [][]string{
				{"Row 1 Col 1", "Row 1 Col 1.5", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 1.5", "Row 2 Col 2"},
			},
		},
		{
			description: "add at end",
			init: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:   []string{"Row 1 Col 3", "Row 2 Col 3"},
			index: 2,
			want: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2", "Row 1 Col 3"},
				{"Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3"},
			},
		},
		{
			description: "add at end/negative",
			init: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:   []string{"Row 1 Col 3", "Row 2 Col 3"},
			index: -1,
			want: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2", "Row 1 Col 3"},
				{"Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3"},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			ds := New(tt.init)
			err := ds.AddCol(tt.add, tt.index)
			require.NoError(t, err)
			require.True(t, Equals(tt.want, ds.Raw()))
		})
	}
}

func TestAddRow(t *testing.T) {
	testCases := []struct {
		description string
		init        [][]string
		add         [][]string
		want        [][]string
		wantErr     error
	}{
		{
			description: "add nothing",
			init: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add:     [][]string{{}},
			wantErr: ErrColLen,
		},
		{
			description: "add nil",
			init: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add: nil,
			want: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
			}},
		{
			description: "add rows to nothing",
			init:        nil,
			add: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
				{"Row 3 Col 1", "Row 3 Col 2"},
			},
			want: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
				{"Row 3 Col 1", "Row 3 Col 2"},
			},
		},
		{
			description: "add single row to existing",
			init: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
			},
			add: [][]string{{"Row 3 Col 1", "Row 3 Col 2"}},
			want: [][]string{
				{"Row 1 Col 1", "Row 1 Col 2"},
				{"Row 2 Col 1", "Row 2 Col 2"},
				{"Row 3 Col 1", "Row 3 Col 2"},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			ds := New(tt.init)
			err := ds.AddRows(tt.add)
			if tt.wantErr != nil {
				require.EqualError(t, tt.wantErr, err.Error())
				return
			}
			require.NoError(t, err)
			require.True(t, Equals(tt.want, ds.Raw()))
		})
	}
}

func TestRenameCol(t *testing.T) {
	testCases := []struct {
		description string
		init        [][]string
		old         string
		new         string
		want        [][]string
		wantErr     error
	}{
		{
			description: "rename col",
			init: [][]string{
				{"Col 1", "Col 2", "Col 3"},
				{"Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3"},
				{"Row 3 Col 1", "Row 3 Col 2", "Row 3 Col 3"},
			},
			old: "Col 2",
			new: "Middle Column",
			want: [][]string{
				{"Col 1", "Middle Column", "Col 3"},
				{"Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3"},
				{"Row 3 Col 1", "Row 3 Col 2", "Row 3 Col 3"},
			}, wantErr: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			ds := New(tt.init)
			err := ds.RenameCol(tt.old, tt.new)
			if tt.wantErr != nil {
				require.EqualError(t, tt.wantErr, err.Error())
				return
			}
			require.NoError(t, err)
			t.Log(tt.want)
			t.Log(ds.Raw())
			require.True(t, Equals(tt.want, ds.Raw()))
		})
	}
}

func TestMoveCol(t *testing.T) {
	testCases := []struct {
		description string
		init        [][]string
		toMove      string
		newIdx      int
		want        [][]string
		wantErr     error
	}{
		{
			description: "move col",
			init: [][]string{
				{"Col 1", "Col 2", "Col 3"},
				{"Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3"},
				{"Row 3 Col 1", "Row 3 Col 2", "Row 3 Col 3"},
			},
			toMove: "Col 1",
			newIdx: 1,
			want: [][]string{
				{"Col 2", "Col 1", "Col 3"},
				{"Row 2 Col 2", "Row 2 Col 1", "Row 2 Col 3"},
				{"Row 3 Col 2", "Row 3 Col 1", "Row 3 Col 3"},
			}, wantErr: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			ds := New(tt.init)
			err := ds.MoveCol(tt.toMove, tt.newIdx)
			if tt.wantErr != nil {
				require.EqualError(t, tt.wantErr, err.Error())
				return
			}
			require.NoError(t, err)
			t.Log(tt.want)
			t.Log(ds.Raw())
			require.True(t, Equals(tt.want, ds.Raw()))
		})
	}
}
