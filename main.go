package main

import (
	"fmt"
	initlog "log"
	"os"
	"os/signal"
	"syscall"

	"rs-drop-emulator/general"
	"rs-drop-emulator/runescape"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	log      *zap.SugaredLogger
	discord  *discordgo.Session
	commands = []*discordgo.ApplicationCommand{
		runescape.GiantMoleCommand,
		runescape.ZilyanaCommand,
		runescape.GraardorCommand,
		runescape.KreearraCommand,
		runescape.KrilCommand,
		runescape.VindictaCommand,
		runescape.HelwyrCommand,
		runescape.TwinfuriesCommand,
		runescape.GregorovicCommand,
		runescape.VoragoCommand,
		general.HelpCommand,
		general.ContributeCommand,
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"giantmole":  runescape.GiantMole,
		"graardor":   runescape.Graardor,
		"zilyana":    runescape.Zilyana,
		"kreearra":   runescape.Kreearra,
		"kril":       runescape.Kril,
		"vindicta":   runescape.Vindicta,
		"helwyr":     runescape.Helwyr,
		"twinfuries": runescape.Twinfuries,
		"gregorovic": runescape.Gregorovic,
		"vorago":     runescape.Vorago,
		"help":       general.Help,
		"contribute": general.Contribute,
	}
	botToken string
)

func init() {
	configLogger()
	general.ConfigLogger()
	runescape.ConfigLogger()

	godotenv.Load()
	botToken = os.Getenv("DISCORD_BOT_TOKEN")

	var err error
	discord, err = discordgo.New(fmt.Sprintf("Bot %v", botToken))
	if err != nil {
		log.Fatalf("couldn't set up bot; %v", err)
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	defer log.Sync()
	defer discord.Close()

	startBot()
}

// Starts the connection and awaits a termination signal
func startBot() {
	err := discord.Open()
	if err != nil {
		log.Fatalf("error opening connection %v", err)
	}

	// Use guild only commands when testing, to propagate changes faster
	env := os.Getenv("RS_DROP_EMULATOR_ENV")
	guildId := "512644466281152526"
	if env == "PROD" {
		guildId = ""
	}

	for _, v := range commands {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, guildId, v)
		if err != nil {
			log.Fatalf("cannot create '%v' command: %v", v.Name, err)
		}
	}

	log.Info("Bot succesfully started up and listening")

	// Await termination
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
	s := <-sc

	log.Infof("shutting down because with signal %v", s)
	discord.Close()
}

func configLogger() {
	var err error
	var logger *zap.Logger

	env := os.Getenv("RS_DROP_EMULATOR_ENV")

	if env == "PROD" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		initlog.Fatalf("couldn't config logger; %v", err)
	}

	log = logger.Sugar().Named("main")
	zap.ReplaceGlobals(logger)
}
