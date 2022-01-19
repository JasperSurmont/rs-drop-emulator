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

## Project structure

Entry point is in [main.go](./main.go). It loads all the application commands found in [beasts](./runescape/beasts) and [general](./general) and starts the bot.

In [beasts](./runescape/beasts), every file has a similar structure specifying the droptables, and emulating the drop. Since a lot of the drop mechanisms are repeated throughout the game, we generalize some of them in [core](./runescape/core). There we have all sorts of useful functions in [general](./runescape/core/general.go) and more specific ones in files like [gwd1](./runescape/core/gwd1.go).

Interacting with the RS Api and cache (to get the GE Prices) is specified in [util]('./runescape/util). Because the RS API is very unstable and not really trustworthy, it's important that when any change happens, it is first tested using [rsapi_test]('./runescape/util/rsapi_test.go).

## TODO

- More tests
- Add bosses
- Add option for uniques only
- _Add clues?_
