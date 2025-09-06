package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"blackjack/internal/game"
)

type UI struct {
	sc *bufio.Scanner
}

func New() *UI {
	return &UI{sc: bufio.NewScanner(os.Stdin)}
}

func (ui *UI) Prompt(msg string) (string, bool) {
	fmt.Print(msg)
	if !ui.sc.Scan() {
		return "", false
	}
	return ui.sc.Text(), true
}

func (ui *UI) Println(msg string) {
	fmt.Println(msg)
}

func (ui *UI) Redraw(g *game.Game, step string) {
	// ANSI clear + home - https://stackoverflow.com/a/37778152
	fmt.Print("\033[H\033[2J")
	printTable(g, step)
}

func (ui *UI) PrintDealerHit(dealerCards []game.Card) {
	fmt.Println("\n[Dealer hits]")
	up := dealerCards[0].String()
	rest := make([]string, 0, len(dealerCards)-1)
	for _, c := range dealerCards[1:] {
		rest = append(rest, c.String())
	}
	total, soft := handValue(dealerCards)
	fmt.Printf("Dealer: [%s] %s  => %d", up, strings.Join(rest, " "), total)
	if soft {
		fmt.Print(" (soft)")
	}
	fmt.Println()
}

func (ui *UI) PrintDealerStand(dealerCards []game.Card) {
	total, soft := handValue(dealerCards)
	header := "\n[Dealer stands]"
	if total > 21 {
		header = "\n[Dealer busts]"
	}
	fmt.Println(header)

	up := dealerCards[0].String()
	rest := make([]string, 0, len(dealerCards)-1)
	for _, c := range dealerCards[1:] {
		rest = append(rest, c.String())
	}
	fmt.Printf("Dealer: [%s] %s  => %d", up, strings.Join(rest, " "), total)
	if soft {
		fmt.Print(" (soft)")
	}
	fmt.Println()
}

func handValue(cards []game.Card) (int, bool) {
	return gameHandValue(cards)
}

func gameHandValue(cards []game.Card) (total int, soft bool) {
	total = 0
	aces := 0
	for _, c := range cards {
		switch c.Rank {
		case game.Ace:
			aces++
			total += 11
		case game.Ten, game.Jack, game.Queen, game.King:
			total += 10
		default:
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

func printTable(g *game.Game, step string) {
	fmt.Println("\n==================================================")
	fmt.Printf("BLACKJACK — %s | Bankroll: %d | Shoe Remaining: %d\n", step, g.Bankroll, g.CardsRemaining())
	fmt.Println("--------------------------------------------------")
	// dirtbag dealer
	fmt.Print("Dealer:  ")
	if g.Dealer.ShowHole {
		fmt.Printf("[%s] [%s]\n", g.Dealer.UpCard.String(), g.Dealer.HoleCard.String())
	} else {
		fmt.Printf("[%s] [??]\n", g.Dealer.UpCard.String())
	}
	fmt.Println("--------------------------------------------------")
	// seats
	for i, seat := range g.Seats {
		fmt.Printf("Seat %d:\n", i+1)
		if len(seat.Hands) == 0 {
			fmt.Println("  (empty)")
			continue
		}
		for j, h := range seat.Hands {
			tot, soft := gameHandValue(h.Cards)
			status := ""
			switch {
			case h.Busted:
				status = "BUST"
			case isBlackjack(h.Cards):
				status = "BLACKJACK"
			case h.Stood:
				status = "STAND"
			case h.Doubled && h.Finished:
				status = "DOUBLE"
			}
			var cards []string
			for _, c := range h.Cards {
				cards = append(cards, c.String())
			}
			fmt.Printf("  Hand %d (bet=%d, splits=%d): %s  => %d", j+1, h.Bet, h.SplitOrigin, strings.Join(cards, " "), tot)
			if soft {
				fmt.Print(" (soft)")
			}
			if status != "" {
				fmt.Printf("  [%s]", status)
			}
			fmt.Println()
		}
	}
	fmt.Println("==================================================")
}

func isBlackjack(cards []game.Card) bool {
	if len(cards) != 2 {
		return false
	}
	t, _ := gameHandValue(cards)
	return t == 21
}
