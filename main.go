package main

import (
	"fmt"
	initlog "log"
	"os"
	"os/signal"
	"syscall"

	"rs-drop-emulator/beasts"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	log      *zap.SugaredLogger
	discord  *discordgo.Session
	commands = []*discordgo.ApplicationCommand{
		beasts.VindictaCommand,
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"vindicta": beasts.Vindicta,
	}
	botToken string
)

func init() {
	configLogger()
	beasts.ConfigLogger(log)

	godotenv.Load()
	botToken = os.Getenv("DISCORD_BOT_TOKEN")

	discord = createBot()
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
}

// Creates and returns a new discord session
func createBot() *discordgo.Session {
	// Create new session
	discord, err := discordgo.New(fmt.Sprintf("Bot %v", botToken))
	if err != nil {
		log.Fatalf("couldn't set up bot; %v", err)
	}

	return discord
}

// Starts the connection and awaits a termination signal
func startBot() {
	// Start listening
	err := discord.Open()
	if err != nil {
		log.Fatalf("error opening connection %v", err)
	}

	for _, v := range commands {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "512644466281152526", v)
		if err != nil {
			log.Fatalf("cannot create '%v' command: %v", v.Name, err)
		}
	}

	log.Info("Bot succesfully started up and listening")

	// Await termination
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, os.Interrupt, os.Kill)
	s := <-sc

	log.Infof("shutting down because with signal %v", s)
	discord.Close()
}
