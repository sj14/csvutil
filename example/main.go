package main

import (
	"log"

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

	ds := csvutil.New(records, true)

	if err := ds.AddColumn([]string{"asd", "1", "2", "3"}, 1); err != nil {
		log.Fatalln(err)
	}

	if err := ds.RenameColumn("username", "nick"); err != nil {
		log.Fatalln(err)
	}

	log.Println(ds.Raw())

	if err := ds.DeleteColumn("nick"); err != nil {
		log.Fatalln(err)
	}

	if err := ds.Write(); err != nil {
		log.Fatalln(err)
	}
}
