package runescape

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var (
	voragoUrl             = "https://runescape.wiki/images/Vorago.png?3a75c"
	voragoCommonDroptable = []Drop{
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
	voragoRareDroptable = []Drop{
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
	voragoEnergy = Drop{
		Name:   "Tectonic energy",
		Rate:   4.0 / 5.0,
		Bold:   true,
		Amount: 2,
	}
)

var VoragoCommand = &discordgo.ApplicationCommand{
	Name:        "vorago",
	Description: "Emulate a Vorago team drop. This emulates all 5 piles.",
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

	drops := EmulateDropGwdVorago(amount, voragoRareDroptable, voragoCommonDroptable, voragoEnergy)

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := MakeDropList(dropWithPrice, drops, total, ok)

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

func EmulateDropGwdVorago(amount int64, rareDroptable []Drop, commonDroptable []Drop, energyDrop Drop) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var sum float64 = 0
		for _, d := range rareDroptable {
			sum += d.Rate
		}
		rareDropped := false

		// Vorago has 5 drop piles
		for j := 0; j < 5; j++ {
			var drop Drop

			// Make copy so we don't adjust orig value
			cd := commonDroptable

			roll := rand.Float64()

			// Roll for energy
			if roll < energyDrop.Rate {
				// Make a copy so that we don't adjust the previous one
				addEnergy := energyDrop
				addEnergy.SetAmount()
				addDropValueToMap(drops, &addEnergy)
			}

			roll = rand.Float64()

			if roll < sum && !rareDropped {
				drop = rareDroptable[rand.Intn(len(rareDroptable))]
				rareDropped = true
			} else {
				drop = cd[rand.Intn(len(cd))]
			}

			drop.SetAmount()
			addDropValueToMap(drops, &drop)
		}

	}
	return drops
}
