
package main

import (
	"flag"
	"log"
	"os"

	"github.com/rwcarlsen/flashcard/flash"
)

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

	words := map[string]int{}
	clean := true
	for _, c := range set.Cards {
		if line, ok := words[c.Front]; ok {
			log.Printf("line %v: duplicate of line %v", c.Line, line)
			clean = false
		} else {
			words[c.Front] = c.Line
		}
	}

	if clean {
		log.Printf("File %v checks out clean", path)
	}
}
