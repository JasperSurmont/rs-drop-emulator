package simulations

import (
	"github.com/bwmarrin/discordgo"
)

var (
	graardorUrl    = "https://runescape.wiki/images/General_Graardor.png?c6b33"
	graardorTables = dropTables{
		commonDroptable: []Drop{
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
		},
		uncommonDroptable: []Drop{
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
		},
		uniqueDroptable: []Drop{
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
		},
		extra: Drop{
			Name:   "Bandos hilt",
			Amount: 1,
		},
		alwaysDroptable: []Drop{
			{
				Name:   "Ourg bones",
				Amount: 1,
			},
		},
	}
)

var GraardorCommand = &discordgo.ApplicationCommand{
	Name:        "graardor",
	Description: "Simulate a General Graardor drop",
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
	simulateDrop(s, i, simulateDropGwd1, graardorTables, graardorUrl, GraardorCommand.Name, "General Graardor")
}
