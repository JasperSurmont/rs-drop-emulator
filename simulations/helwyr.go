package simulations

import (
	"github.com/bwmarrin/discordgo"
)

var (
	helwyrUrl    = "https://runescape.wiki/images/Helwyr.png?8740d"
	helwyrTables = dropTables{
		commonDroptable: []Drop{
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
		},
		uncommonDroptable: []Drop{
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
		},
		uniqueDroptable: []Drop{
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
		},
		alwaysDroptable: []Drop{
			{
				Name:   "Bones",
				Amount: 1,
			},
		},
	}
)

var HelwyrCommand = &discordgo.ApplicationCommand{
	Name:        "helwyr",
	Description: "Simulate a Helwyr drop with full reputation",
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
	simulateDrop(s, i, simulateDropGwd2, helwyrTables, helwyrUrl, HelwyrCommand.Name, "Helwyr")
}
