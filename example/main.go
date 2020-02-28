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

	ds := csvutil.New(records)

	if err := ds.AddColumn("ASD", []string{"asd", "asd", "asd"}, 3); err != nil {
		log.Fatalln(err)
	}

	log.Println(ds.Raw())

	if err := ds.DeleteColumn("username"); err != nil {
		log.Fatalln(err)
	}

	if err := ds.WriteAll(); err != nil {
		log.Fatalln(err)
	}
}
