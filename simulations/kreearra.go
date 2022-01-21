package simulations

import (
	"github.com/bwmarrin/discordgo"
)

var (
	kreearraUrl    = "https://runescape.wiki/images/Kree%27arra.png?fcdb7"
	kreearraTables = dropTables{
		commonDroptable: []Drop{
			{
				Name:   "Small spiky rune salvage",
				Amount: 1,
			},
			{
				Name:        "Rune bolts",
				AmountRange: [2]int{18, 25},
			},
			{
				Name:        "Rune arrow",
				AmountRange: [2]int{100, 105},
			},
			{
				Name:   "Black dragonhide body",
				Amount: 1,
			},
		},
		uncommonDroptable: []Drop{
			{
				Name:   "Medium bladed rune salvage",
				Amount: 1,
			},
			{
				Name:        "Dragon bolts (e)",
				AmountRange: [2]int{2, 15},
			},
			{
				Name:        "Grimy dwarf weed",
				AmountRange: [2]int{5, 22},
			},
			{
				Name:   "Super ranging potion (3)",
				Amount: 3,
				OtherDrops: []Drop{
					{
						Name:   "Super defence (3)",
						Amount: 3,
					},
				},
			},
			{
				Name:   "Dwarf weed seed",
				Amount: 3,
			},
			{
				Name:        "Crushed nest",
				AmountRange: [2]int{12, 15},
			},
			{
				Name:   "Crystal key",
				Amount: 1,
			},
			{
				Name:   "Yew seed",
				Amount: 1,
			},
		},
		uniqueDroptable: []Drop{
			{
				Name: "Armadyl helmet",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
			{
				Name: "Armadyl chestplate",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
			{
				Name: "Armadyl chainskirt",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
			{
				Name: "Armadyl gloves",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
			{
				Name: "Armadyl boots",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
			{
				Name: "Armadyl buckler",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
		},
		extra: Drop{
			Name:   "Armadyl hilt",
			Amount: 1,
		},
		alwaysDroptable: []Drop{
			{
				Name:   "Big bones",
				Amount: 1,
			},
			{
				Name:        "Feather",
				AmountRange: [2]int{1, 15},
			},
		},
	}
)

var KreearraCommand = &discordgo.ApplicationCommand{
	Name:        "kreearra",
	Description: "Simulate a Command Kreearra drop",
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

func Kreearra(s *discordgo.Session, i *discordgo.InteractionCreate) {
	simulateDrop(s, i, simulateDropGwd1, kreearraTables, kreearraUrl, KreearraCommand.Name, "Kree'arra")

}
