package main

import (
	"fmt"
	initlog "log"
	"os"
	"os/signal"
	"syscall"

	"rs-drop-emulator/general"
	"rs-drop-emulator/runescape/beasts"
	"rs-drop-emulator/runescape/util"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	log      *zap.SugaredLogger
	discord  *discordgo.Session
	commands = []*discordgo.ApplicationCommand{
		beasts.VindictaCommand,
		beasts.HelwyrCommand,
		beasts.TwinfuriesCommand,
		beasts.GregorovicCommand,
		general.HelpCommand,
		general.ContributeCommand,
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"vindicta":   beasts.Vindicta,
		"helwyr":     beasts.Helwyr,
		"twinfuries": beasts.Twinfuries,
		"gregorovic": beasts.Gregorovic,
		"help":       general.Help,
		"contribute": general.Contribute,
	}
	botToken string
)

func init() {
	configLogger()
	beasts.ConfigLogger()
	general.ConfigLogger()
	util.ConfigLogger()

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
