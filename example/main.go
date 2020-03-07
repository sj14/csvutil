package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sj14/csvutil"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	records := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	ds := csvutil.New(records)

	if err := ds.AddCol([]string{"asd", "1", "2", "3"}, 1); err != nil {
		log.Fatalln(err)
	}

	lastNames, err := ds.ExtractCol("last_name")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(lastNames)

	if err := ds.RenameCol("username", "nick"); err != nil {
		log.Fatalln(err)
	}

	// log.Println(ds.Raw())

	if err := ds.DeleteCol("nick"); err != nil {
		log.Fatalln(err)
	}

	var addRowNumber = func(val string, i int) string { return fmt.Sprintf("%v (%v)", val, i) }

	if err := ds.ModifyCol("first_name", addRowNumber); err != nil {
		log.Fatalln(err)
	}

	if err := ds.Write(os.Stdout); err != nil {
		log.Fatalln(err)
	}
}
