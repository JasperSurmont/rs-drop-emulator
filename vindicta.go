package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	vindictaUrl             = "https://runescape.wiki/images/Vindicta.png?41b58"
	vindictaCommonDroptable = []Drop{
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

	vindictaUncommonDroptable = []Drop{
		{
			Name:        "Phasmatite stone spirit",
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
			Name:        "Large plated rune salvage",
			AmountRange: [2]int{8, 15},
		},
		{
			Name:        "Dragon bones",
			AmountRange: [2]int{150, 250},
		},
		{
			Name:        "Black dragonhide",
			AmountRange: [2]int{25, 44},
		},
	}

	vindictaRareDroptable = []Drop{
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
			Name: "Dragon Rider lance",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Crest of Zaros",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Zarosian essence",
			Bold: true,
		},
	}

	vindictaAlwaysDroptable = []Drop{
		{
			Name:   "Dragon bones",
			Amount: 1,
		},
	}
)

var VindictaCommand = &discordgo.ApplicationCommand{
	Name:        "vindicta",
	Description: "Emulate a vindicta drop with full reputation",
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

func Vindicta(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1
	if amountOpt := GetOptionWithName(i.ApplicationCommandData().Options, "amount"); amountOpt.Name != "" {
		amount = amountOpt.IntValue()
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

	drops := emulateDropGwd2(amount, vindictaRareDroptable, vindictaUncommonDroptable, vindictaCommonDroptable)
	if enableGuarantees := GetOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		AddAlwaysDroptable(amount, &drops, vindictaAlwaysDroptable)
	}

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: vindictaUrl,
		},
		Title:       fmt.Sprintf("You killed vindicta %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
	log.Info("command executed", "command", VindictaCommand.Name)
}
