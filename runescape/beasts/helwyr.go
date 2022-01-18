package beasts

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	helwyrUrl             = "https://runescape.wiki/images/Helwyr.png?8740d"
	helwyrCommonDroptable = []Drop{
		{
			Name:        "Drakolith stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Orichalcite stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Uncut diamond",
			AmountRange: [2]int{20, 30},
		},
		{
			Name:        "Grimy dwarf weed",
			AmountRange: [2]int{20, 30},
		},
		{
			Name:        "Raw shark",
			AmountRange: [2]int{45, 60},
		},
	}

	helwyrUncommonDroptable = []Drop{
		{
			Name:        "Phasmatite stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Necrite stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Coins",
			AmountRange: [2]int{60000, 82816},
		},
		{
			Name:        "Uncut dragonstone",
			AmountRange: [2]int{8, 12},
		},
		{
			Name:        "Magic logs",
			AmountRange: [2]int{175, 350},
		},
		{
			Name:        "Large bladed rune salvage",
			AmountRange: [2]int{10, 20},
		},
		{
			Name:        "Crystal key",
			AmountRange: [2]int{2, 4},
		},
		{
			Name:        "Grimy lantadyme",
			AmountRange: [2]int{90, 120},
		},
	}

	helwyrRareDroptable = []Drop{
		{
			Rate: 1.0 / 179.0,
			Name: "Dormant anima core helm",
			Bold: true,
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Dormant anima core body",
			Bold: true,
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Dormant anima core legs",
			Bold: true,
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Orb of the Cywir elders",
			Bold: true,
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Crest of Seren",
			Bold: true,
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Wand of the Cywir elders",
			Bold: true,
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Serenic essence",
			Bold: true,
		},
	}

	helwyrAlwaysDroptable = []Drop{
		{
			Name:   "Bones",
			Amount: 1,
		},
	}
)

var HelwyrCommand = &discordgo.ApplicationCommand{
	Name:        "helwyr",
	Description: "Emulate a helwyr drop with full reputation",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "amount",
			Description: "Amount of kills",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "enable-guarantees",
			Description: "Set to false to remove the guaranteed drops",
			Required:    false,
		},
	},
}

func Helwyr(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1
	if lenOpt := len(i.ApplicationCommandData().Options); lenOpt >= 1 {
		amount = i.ApplicationCommandData().Options[0].IntValue()
	}

	// Replace this with max value later
	if amount > MAX_AMOUNT_ROLLS || amount < 1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The amount has to be between 1 and %v", MAX_AMOUNT_ROLLS),
			},
		})
		return
	}

	drops := EmulateDrop(commonRateWithoutRare, amount, helwyrRareDroptable, helwyrUncommonDroptable, helwyrCommonDroptable)
	if lenOpt := len(i.ApplicationCommandData().Options); lenOpt < 2 || i.ApplicationCommandData().Options[1].BoolValue() {
		AddAlwaysDroptable(amount, &drops, helwyrAlwaysDroptable)
	}

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := makeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: helwyrUrl,
		},
		Title:       fmt.Sprintf("You killed helwyr %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
