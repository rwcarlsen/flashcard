
package main

import (
	"fmt"
	"log"
	"flag"
	"os"
	"strings"
	"math/rand"

	"github.com/rwcarlsen/flashcard/flash"
)

var sw = flag.Float64("sw", 3, "score weight for flash probability")
var tw = flag.Float64("tw", 1, "last view weight for flash probability")
var back = flag.Bool("back", false, "flash card backs instead of fronts")
var bi = flag.Bool("bi", false, "randomly show either front or back of cards")
var count = flag.Int("n", 5, "number of times/cards to flash")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
		if *bi {
			*back = rand.Float64() > 0.5
		}
		if *back {
			fmt.Printf(" Back: %v", c.Back)
			scanline()
			fmt.Printf("Front: %v\n Pass: ", c.Front)
			if err := scanline(&ans); err != nil {
				save(set, path)
				log.Fatal(err)
			}
		} else {
			fmt.Printf("Front: %v", c.Front)
			scanline()
			fmt.Printf(" Back: %v\n Pass: ", c.Back)
			if err := scanline(&ans); err != nil {
				save(set, path)
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
	save(set, path)
}

func save(set *flash.Set, path string) {
	// save changes
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	} else if err := set.Save(f); err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func scanline(v ...interface{}) error {
	if _, err := fmt.Scanln(v...); err != nil {
		if _, err := fmt.Scanln(); err != nil {
			return err
		}
	}
	return nil
}
