package runescape

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	graardorUrl             = "https://runescape.wiki/images/General_Graardor.png?c6b33"
	graardorCommonDroptable = []Drop{
		{
			Name:   "Medium bladed rune salvage",
			Amount: 1,
		},
		{
			Name:   "Huge plated rune salvage",
			Amount: 1,
		},
		{
			Name:   "Orichalcite stone spirit",
			Amount: 3,
		},
		{
			Name:   "Drakolith stone spirit",
			Amount: 3,
		},
		{
			Name:   "Snapdragon seed",
			Amount: 1,
		},
	}
	graardorUncommonDroptable = []Drop{
		{
			Name:   "Ourg bones",
			Amount: 3,
		},
		{
			Name:   "Medium spiky rune salvage",
			Amount: 1,
		},
		{
			Name:   "Large bladed rune salvage",
			Amount: 1,
		},
		{
			Name:   "Runite stone spirit",
			Amount: 3,
		},
		{
			Name:        "Magic logs",
			AmountRange: [2]int{15, 20},
		},
		{
			Name:   "Super restore (4)",
			Amount: 3,
		},
		{
			Name:   "Grimy snapdragon",
			Amount: 3,
		},
	}
	graardorRareDroptable = []Drop{
		{
			Name: "Bandos helmet",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Bandos chestplate",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Bandos tassets",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Bandos gloves",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Bandos boots",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
		{
			Name: "Bandos warshield",
			Rate: 1.0 / 384.0,
			Bold: true,
		},
	}
	bandosHilt = Drop{
		Name:   "Bandos hilt",
		Amount: 1,
	}
	graardorAlwaysDroptable = []Drop{
		{
			Name:   "Ourg bones",
			Amount: 1,
		},
	}
)

var GraardorCommand = &discordgo.ApplicationCommand{
	Name:        "graardor",
	Description: "Emulate a General Graardor drop",
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

func Graardor(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	drops := emulateDropGwd1(amount, graardorRareDroptable, graardorUncommonDroptable, graardorCommonDroptable, bandosHilt)
	if enableGuarantees := GetOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		AddAlwaysDroptable(amount, &drops, graardorAlwaysDroptable)
	}

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: graardorUrl,
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
	log.Infow("command executed", "command", GraardorCommand.Name)

}
