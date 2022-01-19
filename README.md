# Runescape 3 Drop emulator

A discord bot that simulates drops. 
The goal is that the user can input a variety of variables that have an influence on the droprate

## Currently implemented bosses

All GWD1 bosses:
- [K'ril](./runescape/beasts/kril.go)
- [Graardor](./runescape/beasts/graardor.go)
- [Kree'arra](./runescape/beasts/kreearra.go)
- [Zilyana](./runescape/beasts/zilyana.go)

All GWD2 bosses:
- [Vindicta](./runescape/beasts/vindicta.go)
- [Helwyr](./runescape/beasts/helwyr.go)
- [Gregorovic](./runescape/beasts/gregorovic.go)
- [Twin Furies](./runescape/beasts/twinfuries.go)

## TODO

- Implement cache such that not every request has to fetch the price of an item (prices change daily)
- React with embed instead of plain message (and add picture etc)
- More tests
