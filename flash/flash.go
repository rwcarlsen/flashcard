
package flash

import (
	"io"
	"io/ioutil"
	"encoding/json"
	"time"
	"math/rand"
)

const (
	initScore = 0
)

type Set struct {
	Cards []*Card
	bounds []float64
}

func Load(r io.Reader) (*Set, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var s Set
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Set) Save(w io.Writer) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	} else if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *Set) AddCard(cards ...*Card) {
	s.Cards = append(s.Cards, cards...)
}

func (s *Set) Next(scoreWeight, timeWeight float64) *Card {
	s.calcBounds(scoreWeight, timeWeight)
	r := rand.Float64()
	index := -1
	for i, v := range s.bounds {
		if r < v {
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
		if c.LastView().Before(oldest) {
			oldest = c.LastView()
		}
	}
	return oldest
}

func (s *Set) newestView() time.Time {
	newest := time.Time{}
	for _, c := range s.Cards {
		if c.LastView().After(newest) {
			newest = c.LastView()
		}
	}
	return newest
}

func (s *Set) weighted(c *Card, sw, tw float64, oldest, newest time.Time) float64 {
	return (1 - c.Score()) * sw + float64(newest.Sub(c.LastView())) / float64(newest.Sub(oldest)) * tw
}

type Card struct {
	Front string
	Back string
	Hist []*View
}

func NewCard(front, back string) *Card {
	return &Card{Front: front, Back: back}
}

func (c *Card) Score() float64 {
	if len(c.Hist) == 0 {
		return initScore
	}
	return c.Hist[len(c.Hist)-1].Score
}

func (c *Card) LastView() time.Time {
	if len(c.Hist) == 0 {
		return time.Time{}
	}
	return c.Hist[len(c.Hist)-1].Date
}

func (c *Card) AddView(pass bool) {
	prevScore := c.Score()
	var score float64
	if pass {
		score = prevScore + (1 - prevScore) / 2
	} else {
		score = prevScore / 2
	}

	v := &View{
		Pass: pass,
		Score: score,
		Date: time.Now(),
	}
	c.Hist = append(c.Hist, v)
}

type View struct {
	Pass bool
	Score float64
	Date time.Time
}

