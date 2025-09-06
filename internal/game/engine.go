package game

import (
	"fmt"
	"strings"
)

type UI interface {
	Redraw(g *Game, step string)
	Println(msg string)
	Prompt(msg string) (string, bool)
	PrintDealerHit(dealerCards []Card)
	PrintDealerStand(dealerCards []Card)
}

func (g *Game) Run(ui UI) {
	for {
		ans, ok := ui.Prompt("\nPress ENTER to deal a new round, or 'q' to quit: ")
		if !ok || strings.ToLower(strings.TrimSpace(ans)) == "q" {
			ui.Println("Goodbye.")
			return
		}

		g.resetRound()

		if quit := g.takeBets(ui); quit {
			return
		}

		g.dealInitial()
		ui.Redraw(g, "Initial deal")

		if ended := g.dealerPeekAndResolveIfBJ(ui); ended {
			continue
		}

		if quit := g.playAllPlayers(ui); quit {
			return
		}

		dealerCards := g.playDealer(ui)

		g.settleAndReport(ui, dealerCards)
	}
}

func (g *Game) resetRound() {
	for i := range g.Seats {
		g.Seats[i].Hands = nil
	}
	g.Dealer = Dealer{ShowHole: false}
}

func (g *Game) takeBets(ui UI) (quit bool) {
	for i := range g.Seats {
		bs, ok := ui.Prompt(fmt.Sprintf("Seat %d bet [1]: ", i+1))
		if !ok {
			return true
		}
		bet := parseIntDefault(strings.TrimSpace(bs), 1)
		if bet < 1 {
			bet = 1
		}
		g.Seats[i].Hands = []Hand{{Bet: bet}}
	}
	return false
}

func (g *Game) dealInitial() {
	// one to each player (up)
	for i := range g.Seats {
		g.Seats[i].Hands[0].Cards = append(g.Seats[i].Hands[0].Cards, g.draw())
	}
	// dealer up
	g.Dealer.UpCard = g.draw()
	// second to each player (up)
	for i := range g.Seats {
		g.Seats[i].Hands[0].Cards = append(g.Seats[i].Hands[0].Cards, g.draw())
	}
	// dealer hole
	g.Dealer.HoleCard = g.draw()
	g.Dealer.ShowHole = false
}

// If dealer shows ace or 10 and has blackjack, resolve the round immediately & return true
func (g *Game) dealerPeekAndResolveIfBJ(ui UI) bool {
	upVal := cardValue(g.Dealer.UpCard.Rank)
	if upVal == 10 || g.Dealer.UpCard.Rank == Ace {
		if isBlackjack([]Card{g.Dealer.UpCard, g.Dealer.HoleCard}) {
			g.Dealer.ShowHole = true
			ui.Redraw(g, "Dealer has BLACKJACK")
			net := 0
			for i := range g.Seats {
				for _, h := range g.Seats[i].Hands {
					if isBlackjack(h.Cards) {
						// push (0)
					} else {
						net -= h.Bet
					}
				}
			}
			g.Bankroll += net
			ui.Println(fmt.Sprintf("Round result: %+d (Bankroll: %d)", net, g.Bankroll))
			return true
		}
	}

	// mark natural blackjacks for later
	for i := range g.Seats {
		for j := range g.Seats[i].Hands {
			h := &g.Seats[i].Hands[j]
			if isBlackjack(h.Cards) {
				h.Finished = true
				h.Stood = true
			}
		}
	}
	return false
}

func (g *Game) playAllPlayers(ui UI) (quit bool) {
	for seatIdx := range g.Seats {
		for handIdx := 0; handIdx < len(g.Seats[seatIdx].Hands); handIdx++ {
			if q := g.playOneHand(ui, seatIdx, handIdx); q {
				return true
			}
		}
	}
	return false
}

func (g *Game) playOneHand(ui UI, seatIdx, handIdx int) (quit bool) {
	h := &g.Seats[seatIdx].Hands[handIdx]

	for {
		ui.Redraw(g, fmt.Sprintf("Seat %d — Hand %d", seatIdx+1, handIdx+1))
		if h.Finished || isBlackjack(h.Cards) {
			return false
		}

		tot, soft := handValue(h.Cards)

		// bust check
		if tot > 21 {
			h.Busted = true
			h.Finished = true
			return false
		}

		// auto-stand at 21
		if tot == 21 {
			h.Stood = true
			h.Finished = true
			return false
		}

		// build action list
		actions := []string{"[h]it", "[s]tand"}
		canDouble := !h.Doubled && len(h.Cards) == 2
		if canDouble {
			actions = append(actions, "[d]ouble")
		}
		if canSplit(h.cloneShallow(), g.cfg.MaxSplits) {
			actions = append(actions, "s[p]lit")
		}
		softStr := ""
		if soft {
			softStr = " (soft)"
		}

		ans, ok := ui.Prompt(fmt.Sprintf(
			"Seat %d Hand %d total=%d%s — choose %s: ",
			seatIdx+1, handIdx+1, tot, softStr, strings.Join(actions, "/"),
		))
		if !ok {
			return true
		}

		switch strings.ToLower(strings.TrimSpace(ans)) {
		case "h", "hit":
			h.Cards = append(h.Cards, g.draw())

		case "s", "stand":
			h.Stood = true
			h.Finished = true

		case "d", "double":
			if canDouble {
				h.Bet *= 2
				h.Doubled = true
				h.Cards = append(h.Cards, g.draw())
				nt, _ := handValue(h.Cards)
				if nt > 21 {
					h.Busted = true
				}
				h.Finished = true
			}

		case "p", "split":
			if canSplit(*h, g.cfg.MaxSplits) {
				left := Hand{Bet: h.Bet, SplitOrigin: h.SplitOrigin + 1}
				right := Hand{Bet: h.Bet, SplitOrigin: h.SplitOrigin + 1}
				left.Cards = []Card{h.Cards[0], g.draw()}
				right.Cards = []Card{h.Cards[1], g.draw()}

				g.Seats[seatIdx].Hands[handIdx] = left
				if handIdx+1 >= len(g.Seats[seatIdx].Hands) {
					g.Seats[seatIdx].Hands = append(g.Seats[seatIdx].Hands, right)
				} else {
					tmp := append([]Hand{}, g.Seats[seatIdx].Hands[:handIdx+1]...)
					tmp = append(tmp, right)
					tmp = append(tmp, g.Seats[seatIdx].Hands[handIdx+1:]...)
					g.Seats[seatIdx].Hands = tmp
				}
				h = &g.Seats[seatIdx].Hands[handIdx]
			}

		case "q", "quit":
			return true

		default:
			// ignore unknown inputs, dont feel like error handling
			// gambling degenerates /s
		}

		if h.Finished {
			return false
		}
	}
}

func (g *Game) playDealer(ui UI) []Card {
	g.Dealer.ShowHole = true
	ui.Redraw(g, "Dealer reveals")

	dealerCards := []Card{g.Dealer.UpCard, g.Dealer.HoleCard}

	for {
		dTot, _ := handValue(dealerCards)

		// soft 17 vs hard 17
		stand := dTot >= 17
		if !g.cfg.StandOnS17 && dTot == 17 {
			_, soft := handValue(dealerCards)
			if soft {
				stand = false
			}
		}

		if stand {
			break
		}
		dealerCards = append(dealerCards, g.draw())
		ui.PrintDealerHit(dealerCards)
	}

	if len(dealerCards) >= 2 {
		g.Dealer.UpCard = dealerCards[0]
		g.Dealer.HoleCard = dealerCards[1]
	}
	ui.PrintDealerStand(dealerCards)
	return dealerCards
}

func (g *Game) settleAndReport(ui UI, dealerCards []Card) {
	dTot, _ := handValue(dealerCards)
	dealerBust := dTot > 21

	net := 0
	for i := range g.Seats {
		for _, h := range g.Seats[i].Hands {
			pt, _ := handValue(h.Cards)
			switch {
			case h.Busted:
				net -= h.Bet
			case isBlackjack(h.Cards):
				net += (h.Bet * 3) / 2
			case dealerBust:
				net += 2 * h.Bet
			default:
				if pt > dTot {
					net += h.Bet
				} else if pt < dTot {
					net -= h.Bet
				} // push => 0
			}
		}
	}
	g.Bankroll += net
	ui.Println(fmt.Sprintf("Round result: %+d (Bankroll: %d)", net, g.Bankroll))
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return def
		}
		n = n*10 + int(ch-'0')
	}
	if n == 0 {
		return def
	}
	return n
}
