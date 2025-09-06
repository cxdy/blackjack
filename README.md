# blackjack

I was bored and I like playing blackjack on occasion, figured this could be a fun/harmless way to learn basic strategy

## Defaults
- 3 players
- 6 decks (52x6 = 312 cards)
- Split twice in one hand

### Usage
build it
```
go build cmd/main.go
```
then run it
```
./main
```
You can use the following args:
```
-hands=3 (default=1)
-decks=6 (default=6)
-maxsplits=3 (default=2)
-seed=0 (default=0) // rtfm: https://stats.stackexchange.com/a/354377 
-soft17=false (default=true)
```

