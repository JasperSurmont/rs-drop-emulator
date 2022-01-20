package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	kreearraUrl             = "https://runescape.wiki/images/Kree%27arra.png?fcdb7"
	kreearraCommonDroptable = []Drop{
		{
			Name:   "Small spiky rune salvage",
			Amount: 1,
		},
		{
			Name:        "Rune bolts",
			AmountRange: [2]int{18, 25},
		},
		{
			Name:        "Rune arrow",
			AmountRange: [2]int{100, 105},
		},
		{
			Name:   "Black dragonhide body",
			Amount: 1,
		},
	}
	kreearraUncommonDroptable = []Drop{
		{
			Name:   "Medium bladed rune salvage",
			Amount: 1,
		},
		{
			Name:        "Dragon bolts (e)",
			AmountRange: [2]int{2, 15},
		},
		{
			Name:        "Grimy dwarf weed",
			AmountRange: [2]int{5, 22},
		},
		{
			Name:   "Super ranging potion (3)",
			Amount: 3,
			OtherDrops: []Drop{
				{
					Name:   "Super defence (3)",
					Amount: 3,
				},
			},
		},
		{
			Name:   "Dwarf weed seed",
			Amount: 3,
		},
		{
			Name:        "Crushed nest",
			AmountRange: [2]int{12, 15},
		},
		{
			Name:   "Crystal key",
			Amount: 1,
		},
		{
			Name:   "Yew seed",
			Amount: 1,
		},
	}
	kreearraRareDroptable = []Drop{
		{
			Name: "Armadyl helmet",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Armadyl chestplate",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Armadyl chainskirt",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Armadyl gloves",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Armadyl boots",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Armadyl buckler",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
	}
	aramadylHilt = Drop{
		Name:   "Armadyl hilt",
		Amount: 1,
	}
	kreearraAlwaysDroptable = []Drop{
		{
			Name:   "Big bones",
			Amount: 1,
		},
		{
			Name:        "Feather",
			AmountRange: [2]int{1, 15},
		},
	}
)

var KreearraCommand = &discordgo.ApplicationCommand{
	Name:        "kreearra",
	Description: "Emulate a Command Kreearra drop",
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

func Kreearra(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	drops := emulateDropGwd1(amount, kreearraRareDroptable, kreearraUncommonDroptable, kreearraCommonDroptable, aramadylHilt)
	if enableGuarantees := GetOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		AddAlwaysDroptable(amount, &drops, kreearraAlwaysDroptable)
	}

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: kreearraUrl,
		},
		Title:       fmt.Sprintf("You killed General Graardor %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
	log.Info("command executed", "command", KreearraCommand.Name)

}
