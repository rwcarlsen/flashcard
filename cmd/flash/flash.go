
package main

import (
	"fmt"
	"log"
	"flag"
	"os"

	"github.com/rwcarlsen/flashcard/flash"
)

var sw = flag.Float64("sw", 1, "score weight for flash probability")
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
		var pass bool
		fmt.Printf("----------- Card %v -----------\n", i+1)
		if *back {
			fmt.Printf(" Back: %v\n Pass: ", c.Back)
			if _, err := fmt.Scanf("%t\n", &pass); err != nil {
				log.Fatal()
			}
			fmt.Printf("Front: %v\n", c.Front)
		} else {
			fmt.Printf("Front: %v\n Pass: ", c.Front)
			if _, err := fmt.Scanf("%t\n", &pass); err != nil {
				log.Fatal()
			}
			fmt.Printf(" Back: %v\n", c.Back)
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

// fw returns a string of len w by appending spaces to string val of v.
func fw(v interface{}, width int) string {
	s := fmt.Sprint(v)
	ss := ""
	for i := 0; i < width - len(s); i++ {
		ss += " "
	}
	return s + ss
}
