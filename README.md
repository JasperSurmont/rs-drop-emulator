# Runescape 3 Drop simulator

A discord bot that simulates drops. 
The goal is that the user can input a variety of variables that have an influence on the droprate

## Currently implemented bosses

- [Giant Mole]('./giantmole.go)
- [Vorago]('./vorago.go)

All GWD1 bosses:
- [K'ril](./simulations/kril.go)
- [Graardor](./simulations/graardor.go)
- [Kree'arra](./simulations/kreearra.go)
- [Zilyana](./simulations/zilyana.go)
- [Nex](./simulations/nex.go)

All GWD2 bosses (except telos):
- [Vindicta](./simulations/vindicta.go)
- [Helwyr](./simulations/helwyr.go)
- [Gregorovic](./simulations/gregorovic.go)
- [Twin Furies](./simulations/twinfuries.go)

## Project structure

Entry point is in [main.go](./main.go). It loads all the application commands and starts the bot.

Every simulation (like a boss) has its own file in [simulations](./simulations) (eg [Vindicta](./simulations/vindicta.go)), but sometimes common drop mechanisms are grouped (eg [Gwd1](./simulations/gwd1.go)). The basic workflow of adding a new simulation is as follows:
1. Create the different droptables and other variables like the URL
2. Create the command template of type `*discordgo.ApplicationCommand` 
3. Create a function which takes as argument the amount (if applicable), the droptables and the discordgo.InteractionCreate variable. This will return the dropped items (write main drop logic here)
4. Create a function with the name of the boss or other with the first letter capitalized. In this function use [simulatedrop.go](./simulations/simulatedrop.go) to handle all of the discord stuff. Pass as argument the function you created in step 3.
5. Add the command in [main.go](./main.go)

Interacting with the RS Api and cache (to get the GE Prices) is specified in [rsapi.go]('./rsapi/rsapi). Because the RS API is very unstable and not really trustworthy, it's important that when any change happens, it is first tested using [rsapi_test]('./rsapi/rsapi_test.go).

In every package there's a logger available called `log` which is configured to work with the google cloud platform. Please use only that log. When creating a new package make sure to create one using `logger.CreateLogger(<package name>)`.

Commands that have nothing to do with simulations are just put in the root

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
