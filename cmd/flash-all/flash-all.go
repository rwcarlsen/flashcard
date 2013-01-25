
package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"text/tabwriter"

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

    writ := tabwriter.NewWriter(os.Stdout, 8, 4, 1, ' ', 0)

	fmt.Fprint(writ, "Front\tBack\tScore\n")
	fmt.Fprint(writ, "-------\t------\t------\n")
	for _, c := range set.Cards {
		fmt.Fprintf(writ, "%v\t%v\t%v\n", c.Front, c.Back, c.Score)
	}

	writ.Flush()
}

