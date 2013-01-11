
package main

import (
	"os"
	"fmt"
	"log"
	"flag"

	"github.com/rwcarlsen/flashcard/flash"
)

const w = 30 // fixed col width

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Need 1 arg: <set-path>")
	}
	path := flag.Arg(0)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	var set *flash.Set
	if set, err = flash.Load(f); err != nil {
		log.Fatal(err)
	}
	f.Close()

	fmt.Println(fw("Front", w), fw("Back", w), fw ("Score", 8))
	fmt.Println(fw("-------", w), fw("------", w), fw("------", 8))
	for _, c := range set.Cards {
		fmt.Println(fw(c.Front, w), fw(c.Back, w), fw(c.Score(), 8))
	}
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
