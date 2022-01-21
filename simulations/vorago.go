package simulations

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var (
	voragoUrl    = "https://runescape.wiki/images/Vorago.png?3a75c"
	voragoTables = dropTables{
		commonDroptable: []Drop{
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
		},
		uniqueDroptable: []Drop{
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
		},
		extra: Drop{
			Name:   "Tectonic energy",
			Rate:   4.0 / 5.0,
			Bold:   true,
			Amount: 2,
		},
	}
)

var VoragoCommand = &discordgo.ApplicationCommand{
	Name:        "vorago",
	Description: "Simulate a Vorago team drop. This simulates all 5 piles.",
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
	simulateDrop(s, i, simulateDropVorago, voragoTables, voragoUrl, VoragoCommand.Name, "Vorago")
}

func simulateDropVorago(amount int64, tables dropTables, i *discordgo.InteractionCreate) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var sum float64 = 0
		for _, d := range tables.uniqueDroptable {
			sum += d.Rate
		}
		rareDropped := false

		// Vorago has 5 drop piles
		for j := 0; j < 5; j++ {
			var drop Drop

			// Make copy so we don't adjust orig value
			cd := tables.commonDroptable

			roll := rand.Float64()

			energy := tables.extra.(Drop)
			// Roll for energy
			if roll < energy.Rate {
				// Make a copy so that we don't adjust the previous one
				energy.setAmount()
				addDropValueToMap(drops, &energy)
			}

			roll = rand.Float64()

			if roll < sum && !rareDropped {
				drop = determineDropWithRates(rand.Float64(), tables.uniqueDroptable)
				rareDropped = true
			} else {
				drop = cd[rand.Intn(len(cd))]
			}

			drop.setAmount()
			addDropValueToMap(drops, &drop)
		}

	}
	return drops
}
