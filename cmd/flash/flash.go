
package main

import (
	"fmt"
	"log"
	"flag"
	"os"
	"strings"

	"github.com/rwcarlsen/flashcard/flash"
)

var sw = flag.Float64("sw", 3, "score weight for flash probability")
var tw = flag.Float64("tw", 1, "last view weight for flash probability")
var back = flag.Bool("back", false, "flash card backs instead of fronts")
var count = flag.Int("n", 5, "number of times/cards to flash")

func main() {
	flag.Parse()
	path := flag.Arg(0)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	set, err := flash.Load(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	for i := 0; i < *count; i++ {
		c := set.Next(*sw, *tw)
		ans := ""
		fmt.Printf("----------- Card %v -----------\n", i+1)
		if *back {
			fmt.Printf(" Back: %v", c.Back)
			fmt.Scanln()
			fmt.Printf("Front: %v\n Pass: ", c.Front)
			if _, err := fmt.Scanln(&ans); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("Front: %v", c.Front)
			fmt.Scanln()
			fmt.Printf(" Back: %v\n Pass: ", c.Back)
			if _, err := fmt.Scanln(&ans); err != nil {
				log.Fatal(err)
			}
		}

		pass := false
		switch strings.ToLower(ans) {
		case "t", "true", "y", "yes":
			pass = true
		}
		c.AddView(pass)
	}

	// save changes
	if f, err = os.Create(path); err != nil {
		log.Fatal(err)
	} else if err := set.Save(f); err != nil {
		log.Fatal(err)
	}
	f.Close()
}

