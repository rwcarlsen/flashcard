
package main

import (
	"os"
	"log"
	"flag"

	"github.com/rwcarlsen/flashcard/flash"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 3 {
		log.Fatal("Need 3 args: <set-path> <front-content> <back-content>")
	}

	path := flag.Arg(0)
	front := flag.Arg(1)
	back := flag.Arg(2)

	var set *flash.Set
	f, err := os.Open(path)
	if err != nil {
		log.Printf("set '%v' doesn't exist - new set created.", path)
		set = &flash.Set{}
	} else {
		if set, err = flash.Load(f); err != nil {
			log.Fatal(err)
		}
		f.Close()
	}

	set.AddCard(flash.NewCard(front, back))

	// save changes
	if f, err = os.Create(path); err != nil {
		log.Fatal(err)
	} else if err := set.Save(f); err != nil {
		log.Fatal(err)
	}
	f.Close()
}

