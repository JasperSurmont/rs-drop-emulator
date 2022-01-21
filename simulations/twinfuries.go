package simulations

import (
	"github.com/bwmarrin/discordgo"
)

var (
	twinfuriesUrl    = "https://runescape.wiki/images/Nymora%2C_the_Vengeful.png?230f1"
	twinfuriesTables = dropTables{
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
		},
		uniqueDroptable: []Drop{
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
		},
		alwaysDroptable: []Drop{},
	}
)

var TwinfuriesCommand = &discordgo.ApplicationCommand{
	Name:        "twinfuries",
	Description: "Simulate a Twin Furies drop with full reputation",
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
	simulateDrop(s, i, simulateDropGwd2, twinfuriesTables, twinfuriesUrl, TwinfuriesCommand.Name, "Twin Furies")
}
