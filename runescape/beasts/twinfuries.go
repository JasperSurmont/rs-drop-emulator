package beasts

import (
	"fmt"
	"rs-drop-emulator/runescape/core"

	"github.com/bwmarrin/discordgo"
)

var (
	twinfuriesUrl             = "https://runescape.wiki/images/Nymora%2C_the_Vengeful.png?230f1"
	twinfuriesCommonDroptable = []core.Drop{
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
			AmountRange: [2]int{18, 22},
		},
		{
			Name:        "Grimy dwarf weed",
			AmountRange: [2]int{14, 25},
		},
		{
			Name:        "Raw shark",
			AmountRange: [2]int{45, 55},
		},
	}

	twinfuriesUncommonDroptable = []core.Drop{
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
			AmountRange: [2]int{150, 250},
		},
		{
			Name:        "Large bladed rune salvage",
			AmountRange: [2]int{5, 10},
		},
		{
			Name:        "Wine of Zamorak",
			AmountRange: [2]int{18, 22},
		},
		{
			Name:        "Infernal ashes",
			AmountRange: [2]int{150, 249},
		},
	}

	twinfuriesRareDroptable = []core.Drop{
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core helm",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core body",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core legs",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Orb of the Cywir elders",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Blade of Avaryss",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Blade of Nymora",
			Bold: true,
		},
		{
			Rate: 1.0 / 64.0,
			Name: "Zamorakian essence",
			Bold: true,
		},
	}

	twinfuriesAlwaysDroptable = []core.Drop{}
)

var TwinfuriesCommand = &discordgo.ApplicationCommand{
	Name:        "twinfuries",
	Description: "Emulate a Twin Furies drop with full reputation",
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

func Twinfuries(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1
	if lenOpt := len(i.ApplicationCommandData().Options); lenOpt >= 1 {
		amount = i.ApplicationCommandData().Options[0].IntValue()
	}

	// Replace this with max value later
	if amount > core.MAX_AMOUNT_ROLLS || amount < 1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The amount has to be between 1 and %v", core.MAX_AMOUNT_ROLLS),
			},
		})
		return
	}

	drops := core.EmulateDropGwd2(amount, twinfuriesRareDroptable, twinfuriesUncommonDroptable, twinfuriesCommonDroptable)
	if lenOpt := len(i.ApplicationCommandData().Options); lenOpt < 2 || i.ApplicationCommandData().Options[1].BoolValue() {
		core.AddAlwaysDroptable(amount, &drops, twinfuriesAlwaysDroptable)
	}

	dropWithPrice, total, ok := core.AmountToPrice(drops)
	core.SortDrops(&dropWithPrice)
	content := core.MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: twinfuriesUrl,
		},
		Title:       fmt.Sprintf("You killed the Twin Furies %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
