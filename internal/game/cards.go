package game

type Suit int
type Rank int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

var suitRune = []rune{'♣', '♦', '♥', '♠'}

const (
	Two   Rank = 2
	Three Rank = 3
	Four  Rank = 4
	Five  Rank = 5
	Six   Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine  Rank = 9
	Ten   Rank = 10
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Ace   Rank = 14
)

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	var r string
	switch c.Rank {
	case Ten:
		r = "10"
	case Jack:
		r = "J"
	case Queen:
		r = "Q"
	case King:
		r = "K"
	case Ace:
		r = "A"
	default:
		r = itoa(int(c.Rank))
	}
	return r + string([]rune{rune(' '), suitRune[c.Suit]})
}

func itoa(n int) string {
	if n == 10 {
		return "10"
	}
	return string('0' + rune(n))
}

func makeDeck() []Card {
	d := make([]Card, 0, 52)
	for s := Clubs; s <= Spades; s++ {
		for r := Two; r <= King; r++ {
			d = append(d, Card{Rank: r, Suit: s})
		}
		d = append(d, Card{Rank: Ace, Suit: s})
	}
	return d
}

func makeShoe(n int) []Card {
	shoe := make([]Card, 0, 52*n)
	for i := 0; i < n; i++ {
		shoe = append(shoe, makeDeck()...)
	}
	return shoe
}

func shuffle(cards []Card, rng interface{ Intn(n int) int }) {
	for i := len(cards) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}
