package simulations

import (
	"github.com/bwmarrin/discordgo"
)

var (
	gregorovicUrl    = "https://runescape.wiki/images/Gregorovic.png?d97e5"
	gregorovicTables = dropTables{
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
				Name:        "Medium plated rune salvage",
				AmountRange: [2]int{5, 10},
			},
			{
				Name:        "Battlestaff",
				AmountRange: [2]int{45, 60},
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
				Name: "Shadow glaive",
				Bold: true,
			},
			{
				Rate: 1.0 / 256.0,
				Name: "Off-hand shadow glaive",
				Bold: true,
			},
			{
				Rate: 1.0 / 256.0,
				Name: "Crest of Sliske",
				Bold: true,
			},
			{
				Rate: 1.0 / 64.0,
				Name: "Sliskean essence",
				Bold: true,
			},
		},

		alwaysDroptable: []Drop{},
	}
)

var GregorovicCommand = &discordgo.ApplicationCommand{
	Name:        "gregorovic",
	Description: "Simulate a Gregorovic drop with full reputation",
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

func Gregorovic(s *discordgo.Session, i *discordgo.InteractionCreate) {
	simulateDrop(s, i, simulateDropGwd2, gregorovicTables, gregorovicUrl, GregorovicCommand.Name, "Gregorovic")
}
