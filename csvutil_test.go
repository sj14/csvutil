package csvutil

import (
	"testing"
)

func TestNew(t *testing.T) {
	raw := [][]string{
		[]string{"Row 1 Col 1", "Row 1 Col 2"},
		[]string{"Row 2 Col 1", "Row 2 Col 2"},
	}
	ds := New(raw, false)

	if ds.Raw()[0][0] != raw[0][0] {
		t.Fail()
	}
	if ds.Raw()[0][1] != raw[0][1] {
		t.Fail()
	}
	if ds.Raw()[1][0] != raw[1][0] {
		t.Fail()
	}
	if ds.Raw()[1][1] != raw[1][1] {
		t.Fail()
	}
}
