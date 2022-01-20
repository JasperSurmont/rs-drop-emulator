package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaspersurmont/rs-drop-emulator/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	log      logger.LoggerWrapper
	discord  *discordgo.Session
	commands = []*discordgo.ApplicationCommand{
		GiantMoleCommand,
		ZilyanaCommand,
		GraardorCommand,
		KreearraCommand,
		KrilCommand,
		NexCommand,
		VindictaCommand,
		HelwyrCommand,
		TwinfuriesCommand,
		GregorovicCommand,
		VoragoCommand,
		HelpCommand,
		ContributeCommand,
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"giantmole":  GiantMole,
		"graardor":   Graardor,
		"zilyana":    Zilyana,
		"kreearra":   Kreearra,
		"kril":       Kril,
		"nex":        Nex,
		"vindicta":   Vindicta,
		"helwyr":     Helwyr,
		"twinfuries": Twinfuries,
		"gregorovic": Gregorovic,
		"vorago":     Vorago,
		"help":       Help,
		"contribute": Contribute,
	}
	botToken string
)

func init() {
	log = logger.CreateLogger("main")

	godotenv.Load()
	botToken = os.Getenv("DISCORD_BOT_TOKEN")

	var err error
	discord, err = discordgo.New(fmt.Sprintf("Bot %v", botToken))
	if err != nil {
		log.Fatal(fmt.Sprintf("couldn't set up bot; %v", err))
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
		log.Fatal(fmt.Sprintf("error opening connection %v", err))
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
			log.Fatal(fmt.Sprintf("cannot create '%v' command: %v", v.Name, err))
		}
	}

	log.Info("Bot succesfully started up and listening")

	// Await termination
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
	s := <-sc

	log.Info(fmt.Sprintf("shutting down because with signal %v", s))
	discord.Close()
}
