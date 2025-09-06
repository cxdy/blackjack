package game

import "math/rand"

type Seat struct {
	Hands []Hand
}

type Dealer struct {
	UpCard   Card
	HoleCard Card
	ShowHole bool
}

type Config struct {
	Hands      int
	Decks      int
	MaxSplits  int
	StandOnS17 bool
}

type Game struct {
	cfg       Config
	Seats     []Seat
	Dealer    Dealer
	Bankroll  int
	Shoe      []Card
	drawIndex int
	Rng       *rand.Rand
}

func (g *Game) CardsRemaining() int {
	return len(g.Shoe) - g.drawIndex
}

func New(cfg Config, rng *rand.Rand) *Game {
	g := &Game{
		cfg:   cfg,
		Seats: make([]Seat, cfg.Hands),
		Rng:   rng,
	}
	g.Shoe = makeShoe(cfg.Decks)
	shuffle(g.Shoe, g.Rng)
	return g
}

func (g *Game) cardsRemaining() int {
	return len(g.Shoe) - g.drawIndex
}

func (g *Game) draw() Card {
	if g.drawIndex >= len(g.Shoe) {
		shuffle(g.Shoe, g.Rng)
		g.drawIndex = 0
	}
	c := g.Shoe[g.drawIndex]
	g.drawIndex++
	return c
}
