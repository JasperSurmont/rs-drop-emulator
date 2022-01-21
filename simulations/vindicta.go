package simulations

import (
	"github.com/bwmarrin/discordgo"
)

var (
	vindictaUrl    = "https://runescape.wiki/images/Vindicta.png?41b58"
	vindictaTables = dropTables{
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
				Name:        "Large plated rune salvage",
				AmountRange: [2]int{8, 15},
			},
			{
				Name:        "Dragon bones",
				AmountRange: [2]int{150, 250},
			},
			{
				Name:        "Black dragonhide",
				AmountRange: [2]int{25, 44},
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
				Name: "Dragon Rider lance",
				Bold: true,
			},
			{
				Rate: 1.0 / 256.0,
				Name: "Crest of Zaros",
				Bold: true,
			},
			{
				Rate: 1.0 / 256.0,
				Name: "Zarosian essence",
				Bold: true,
			},
		},
		alwaysDroptable: []Drop{
			{
				Name:   "Dragon bones",
				Amount: 1,
			},
		},
	}
)

var VindictaCommand = &discordgo.ApplicationCommand{
	Name:        "vindicta",
	Description: "Simulate a vindicta drop with full reputation",
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

func Vindicta(s *discordgo.Session, i *discordgo.InteractionCreate) {
	simulateDrop(s, i, simulateDropGwd2, vindictaTables, vindictaUrl, VindictaCommand.Name, "Vindicta")
}
