package beasts

import (
	"fmt"
	"rs-drop-emulator/runescape/core"

	"github.com/bwmarrin/discordgo"
)

var (
	voragoUrl             = "https://runescape.wiki/images/Vorago.png?3a75c"
	voragoCommonDroptable = []core.Drop{
		{
			Name:        "Grimy torstol",
			AmountRange: [2]int{15, 29},
		},
		{
			Name:        "Grapevine seed",
			AmountRange: [2]int{10, 36},
		},
		{
			Name:        "Raw rocktail",
			AmountRange: [2]int{75, 159},
		},
		{
			Name:        "Huge plated rune salvage",
			AmountRange: [2]int{3, 9},
		},
		{
			Name:        "Torstol seed",
			AmountRange: [2]int{4, 10},
		},
		{
			Name:        "Banite stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Black dragonhide",
			AmountRange: [2]int{35, 79},
		},
		{
			Name:        "Onyx bolts (e)",
			AmountRange: [2]int{20, 55},
		},
		{
			Name:        "Magic logs",
			AmountRange: [2]int{35, 128},
		},
		{
			Name:   "Hydrix bolt tips",
			Amount: 50,
		},
		{
			Name:        "Crystal triskelion fragment 1",
			Amount:      1,
			Untradeable: true,
		},
		{
			Name:        "Crystal triskelion fragment 2",
			Amount:      1,
			Untradeable: true,
		}, {
			Name:        "Crystal triskelion fragment 3",
			Amount:      1,
			Untradeable: true,
		},
	}
	voragoRareDroptable = []core.Drop{
		{
			Name: "Seismic wand",
			Rate: 1.0 / 400.0,
			Bold: true,
		},
		{
			Name: "Seismic singularity",
			Rate: 1.0 / 400.0,
			Bold: true,
		},
	}
	voragoEnergy = core.Drop{
		Name:   "Tectonic energy",
		Rate:   4.0 / 5.0,
		Bold:   true,
		Amount: 2,
	}
)

var VoragoCommand = &discordgo.ApplicationCommand{
	Name:        "vorago",
	Description: "Emulate a Vorago team drop. This emulates all 5 piles, so depending on team size and whether you split or not, the value per player changes.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "amount",
			Description: "Amount of kills",
			Required:    false,
		},
	},
}

func Vorago(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1

	if amountOpt := core.GetOptionWithName(i.ApplicationCommandData().Options, "amount"); amountOpt.Name != "" {
		amount = amountOpt.IntValue()
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

	drops := core.EmulateDropGwdVorago(amount, voragoRareDroptable, voragoCommonDroptable, voragoEnergy)

	dropWithPrice, total, ok := core.AmountToPrice(drops)
	core.SortDrops(&dropWithPrice)
	content := core.MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: voragoUrl,
		},
		Title:       fmt.Sprintf("You killed Vorago %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
