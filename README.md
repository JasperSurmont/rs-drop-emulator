# Runescape 3 Drop simulator

A discord bot that simulates drops. 
The goal is that the user can input a variety of variables that have an influence on the droprate

## Currently implemented bosses

- [Giant Mole]('./giantmole.go)
- [Vorago]('./vorago.go)

All GWD1 bosses:
- [K'ril](./kril.go)
- [Graardor](./graardor.go)
- [Kree'arra](./kreearra.go)
- [Zilyana](./zilyana.go)
- [Nex](./nex.go)

All GWD2 bosses (except telos):
- [Vindicta](./vindicta.go)
- [Helwyr](./helwyr.go)
- [Gregorovic](./gregorovic.go)
- [Twin Furies](./twinfuries.go)

## Project structure

Entry point is in [main.go](./main.go). It loads all the application commands and starts the bot.

Every boss has its own file (eg [Vindicta](./vindicta.go)), but sometimes common drop mechanisms are grouped (eg [Gwd1](./gwd1.go)).
Maybe if this gets too messy to put everything in the root we can try to put it in different packages.

Interacting with the RS Api and cache (to get the GE Prices) is specified in [util]('./runescape/util). Because the RS API is very unstable and not really trustworthy, it's important that when any change happens, it is first tested using [rsapi_test]('./runescape/util/rsapi_test.go).

## TODO

### General

- Make discord icon
- _Make website?_

### Code

- More tests
- Rare drop table
- Add option for uniques only
- Hard mode
- _Add clues?_
