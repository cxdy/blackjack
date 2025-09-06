package game

type Hand struct {
	Cards       []Card
	Bet         int
	Finished    bool
	Busted      bool
	Stood       bool
	Doubled     bool
	SplitOrigin int
}

func (h *Hand) cloneShallow() Hand {
	return Hand{
		Cards:       append([]Card{}, h.Cards...),
		Bet:         h.Bet,
		Finished:    h.Finished,
		Busted:      h.Busted,
		Stood:       h.Stood,
		Doubled:     h.Doubled,
		SplitOrigin: h.SplitOrigin,
	}
}

func cardValue(r Rank) int {
	switch r {
	case Jack, Queen, King, Ten:
		return 10
	case Ace:
		return 11
	default:
		return int(r)
	}
}

func handValue(cards []Card) (total int, soft bool) {
	total = 0
	aces := 0
	for _, c := range cards {
		if c.Rank == Ace {
			aces++
			total += 11
		} else if c.Rank >= Ten && c.Rank <= King {
			total += 10
		} else {
			total += int(c.Rank)
		}
	}
	soft = false
	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}
	if aces > 0 && total <= 21 {
		soft = true
	}
	return total, soft
}

func isBlackjack(cards []Card) bool {
	if len(cards) != 2 {
		return false
	}
	t, _ := handValue(cards)
	return t == 21
}

func canSplit(h Hand, maxSplits int) bool {
	if len(h.Cards) != 2 {
		return false
	}
	if h.SplitOrigin >= maxSplits {
		return false
	}
	v1 := cardValue(h.Cards[0].Rank)
	v2 := cardValue(h.Cards[1].Rank)
	return v1 == v2
}
