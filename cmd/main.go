package main

import (
	"flag"
	"math/rand"
	"time"

	"blackjack/internal/game"
	"blackjack/internal/ui/terminal"
)

func main() {
	hands := flag.Int("hands", 1, "number of player hands (1-9)")
	decks := flag.Int("decks", 6, "number of decks in shoe (1-8)")
	maxSplits := flag.Int("maxsplits", 2, "max splits per original hand")
	seed := flag.Int64("seed", 0, "rng seed (0 => random)")
	soft17 := flag.Bool("soft17", true, "stand on soft 17 (default true)")
	flag.Parse()

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(*seed))

	cfg := game.Config{
		Hands:      *hands,
		Decks:      *decks,
		MaxSplits:  *maxSplits,
		StandOnS17: *soft17,
	}

	ui := terminal.New()
	g := game.New(cfg, rng)

	ui.Println("Welcome to CLI Blackjack.")
	ui.Println("Actions: [h]it, [s]tand, [d]ouble (1 card), s[p]lit (if allowed). [q] to quit.")
	ui.Println("Payouts: Blackjack 3:2. If dealer busts: wins pay 2:1. Otherwise even money; push returns bet.")
	ui.Println("Default bet is 1 (you can change at round start). Shoe reshuffles automatically when exhausted.")

	g.Run(ui)
}
