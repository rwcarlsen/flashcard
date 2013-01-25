
package flash

import (
	"fmt"
	"io"
	"time"
	"math/rand"
	"text/tabwriter"
	"bufio"
	"strings"
	"strconv"
)

const (
	initScore = 0.00001
	timeFmt = "2006-01-02 15:04:05.999999999 -0700 MST"
)

type Set struct {
	Cards []*Card
	bounds []float64
}

func Load(r io.Reader) (*Set, error) {
	buf := bufio.NewReader(r)
	s := &Set{}
	for i := 0; true; i++ {
		line, isPrefix, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if isPrefix {
			return nil, fmt.Errorf("line %v: line is too long", i)
		}

		cells := strings.Split(string(line), "  ")
		items := []string{}
		for _, cell := range cells {
			data := strings.TrimSpace(cell)
			if len(data) == 0 {
				continue
			}
			items = append(items, data)
		}

		if len(items) < 2 {
			return nil, fmt.Errorf("line %v: Need at least 2 items per column, got %v", i, len(items))
		}

		score := initScore
		if len(items) > 3 {
			score, err = strconv.ParseFloat(items[3], 64)
			if err != nil {
				return nil, err
			}
		}

		date := time.Now()
		if len(items) > 2 {
			date, err = time.Parse(timeFmt, items[2])
			if err != nil {
				return nil, fmt.Errorf("line %v: %v", i+1, err)
			}
		}

		c := &Card{
			Front: items[0],
			Back: items[1],
			Date: date,
			Score: score,
		}
		s.Cards = append(s.Cards, c)
	}

	return s, nil
}

func (s *Set) Save(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 8, 4, 4, ' ', 0)
	for _, c := range s.Cards {
		if err := c.save(tw); err != nil {
			return err
		}
		if _, err := fmt.Fprint(tw, "\n"); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func (s *Set) AddCard(cards ...*Card) {
	s.Cards = append(s.Cards, cards...)
}

func (s *Set) DuplicateOf(card *Card) *Card {
	for _, c := range s.Cards {
		if c.Front == card.Front {
			return c
		}
	}
	return nil
}

func (s *Set) Next(scoreWeight, timeWeight float64) *Card {
	s.calcBounds(scoreWeight, timeWeight)
	r := rand.Float64()
	index := -1
	for i, v := range s.bounds {
		if r <= v {
			index = i
			break
		}
	}
	return s.Cards[index]
}

func (s *Set) calcBounds(scoreWeight, timeWeight float64) {
	sw := scoreWeight / (scoreWeight + timeWeight)
	tw := 1 - sw

	s.bounds = make([]float64, len(s.Cards))

	oldest := s.oldestView()
	newest := s.newestView()
	norm := 0.0
	for _, c := range s.Cards {
		norm += s.weighted(c, sw, tw, oldest, newest)
	}

	s.bounds[0] = s.weighted(s.Cards[0], sw, tw, oldest, newest) / norm
	for i := 1; i < len(s.Cards); i++ {
		s.bounds[i] = s.bounds[i-1] + (s.weighted(s.Cards[i], sw, tw, oldest, newest) / norm)
	}
}

func (s *Set) oldestView() time.Time {
	oldest := time.Now()
	for _, c := range s.Cards {
		if c.Date.Before(oldest) {
			oldest = c.Date
		}
	}
	return oldest
}

func (s *Set) newestView() time.Time {
	newest := time.Time{}
	for _, c := range s.Cards {
		if c.Date.After(newest) {
			newest = c.Date
		}
	}
	return newest
}

func (s *Set) weighted(c *Card, sw, tw float64, oldest, newest time.Time) float64 {
	// the .0001 is to prevent divide by zero
	return (1 - c.Score) * sw + float64(newest.Sub(c.Date)) / (float64(newest.Sub(oldest))+.00001) * tw
}

type Card struct {
	Front string
	Back string
	// Date card was last viewed
	Date time.Time
	// Score between 0 and 1. 0 means don't know well, 1 means know well
	Score float64
}

func NewCard(front, back string) *Card {
	return &Card{Front: front, Back: back}
}

func (c *Card) save(w io.Writer) error {
	 _, err := fmt.Fprintf(w, "%v\t%v\t%v\t%v", c.Front, c.Back, c.Date.Format(timeFmt), c.Score)
	 return err
}

func (c *Card) AddView(pass bool) {
	if pass {
		c.Score = c.Score + (1 - c.Score) / 2
	} else {
		c.Score = c.Score / 2
	}
	c.Date = time.Now()
}

func (c *Card) String() string {
	return fmt.Sprintf("Front: %v\tBack: %v", c.Front, c.Back)
}

