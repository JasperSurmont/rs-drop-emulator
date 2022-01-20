package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	zilyanaUrl             = "https://runescape.wiki/images/General_Graardor.png?c6b33"
	zilyanaCommonDroptable = []Drop{
		{
			Name:   "Grimy ranarr",
			Amount: 5,
		},
		{
			Name:   "Ranarr seed",
			Amount: 2,
		},
		{
			Name:        "Rune dart",
			AmountRange: [2]int{30, 40},
		},
		{
			Name:   "Large plated rune salvage",
			Amount: 1,
		},
	}
	zilyanaUncommonDroptable = []Drop{
		{
			Name:   "Diamond",
			Amount: 6,
		},
		{
			Name:   "Prayer potion (4)",
			Amount: 3,
		},
		{
			Name:   "Saradomin brew (3)",
			Amount: 3,
			OtherDrops: []Drop{
				{
					Name:   "Super restore (4)",
					Amount: 3,
				},
			},
		},
		{
			Name:   "Super defence (3)",
			Amount: 3,
			OtherDrops: []Drop{
				{
					Name:   "Super magic potion (3)",
					Amount: 3,
				},
			},
		},
		{
			Name:        "Unicorn horn",
			AmountRange: [2]int{5, 10},
		},
		{
			Name:   "Battlestaff",
			Amount: 2,
		},
		{
			Name:   "Huge plated adamant salvage",
			Amount: 1,
		},
		{
			Name:   "Magic seed",
			Amount: 1,
		},
	}
	zilyanaRareDroptable = []Drop{
		{
			Name: "Saradomin sword",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Armadyl crossbow",
			Rate: 1.0 / 768.0,
			Bold: true,
		},
		{
			Name: "Off-hand Armadyl crossbow",
			Rate: 1.0 / 768.0,
			Bold: true,
		},
		{
			Name: "Saradomin's murmur",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Saradomin's hiss",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Saradomin's whisper",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
	}
	saradominHilt = Drop{
		Name:   "Saradomin hilt",
		Amount: 1,
	}
	zilyanaAlwaysDroptable = []Drop{
		{
			Name:   "Bones",
			Amount: 1,
		},
	}
)

var ZilyanaCommand = &discordgo.ApplicationCommand{
	Name:        "zilyana",
	Description: "Emulate a Command Zilyana drop",
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

func Zilyana(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	drops := emulateDropGwd1(amount, zilyanaRareDroptable, zilyanaUncommonDroptable, zilyanaCommonDroptable, saradominHilt)
	if enableGuarantees := GetOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		AddAlwaysDroptable(amount, &drops, zilyanaAlwaysDroptable)
	}

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: zilyanaUrl,
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
	log.Info("command executed", "command", ZilyanaCommand.Name)
}
