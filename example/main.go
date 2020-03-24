package main

import (
	"fmt"
	"os"

	"github.com/sj14/csvutil"
)

/////////////////////////////////////////////
// The examples ignore all error handling! //
/////////////////////////////////////////////
func main() {

	a := [][]string{
		{"a", "b", "c"},
		{"d", "e", "f"},
	}

	b := [][]string{
		{"g", "h", "i"},
		{"j", "k", "l"},
	}

	result := csvutil.Equals(a, b)
	fmt.Println(result)

	records := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	ds := csvutil.New(records)

	ds.AddRow([]string{"my first name", "my last name", "my nick"})

	ds.AddRows([][]string{
		{"my first name 0", "my last name 0", "my nick 0"},
		{"my first name 1", "my last name 1", "my nick 1"},
		{"my first name 2", "my last name 2", "my nick 2"},
	})
	ds.Write(os.Stdout)
	return

	fmt.Println(ds.Raw())

	lastNames, _ := ds.ExtractCol("last_name")
	fmt.Println(lastNames)

	ds.AddCol([]string{"column_headline", "my ow 1", "my row 2", "my row 3"}, 1)

	ds.RenameCol("username", "nick")

	ds.DeleteCol("first_name")

	var addRowNumber = func(val string, i int) string { return fmt.Sprintf("%v (%v)", val, i) }
	ds.ModifyCol("first_name", addRowNumber)

	ds.Write(os.Stdout, csvutil.Delimiter('|'), csvutil.UseCLRF(true))
}
