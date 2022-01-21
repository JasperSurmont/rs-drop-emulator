package simulations

import (
	"github.com/bwmarrin/discordgo"
)

var (
	krilUrl    = "https://runescape.wiki/images/K%27ril_Tsutsaroth.png?11873"
	krilTables = dropTables{
		commonDroptable: []Drop{
			{
				Name:   "Infernal ashes",
				Amount: 5,
			},
			{
				Name:   "Medium bladed rune salvage",
				Amount: 1,
			},
			{
				Name:   "Large plated rune salvage",
				Amount: 1,
			},
			{
				Name:   "Huge plated adamant salvage",
				Amount: 1,
			},
			{
				Name:   "Super attack (3)",
				Amount: 3,
				OtherDrops: []Drop{
					{
						Name:   "Super strength (3)",
						Amount: 3,
					},
				},
			},
			{
				Name:   "Super restore (3)",
				Amount: 1,
				OtherDrops: []Drop{
					{
						Name:   "Zamorak brew (3)",
						Amount: 3,
					},
				},
			},
		},
		uncommonDroptable: []Drop{
			{
				Name:   "Lantadyme seed",
				Amount: 3,
			},
			{
				Name:   "Grimy lantadyme",
				Amount: 10,
			},
			{
				Name:   "Orichalcite stone spirit",
				Amount: 3,
			},
			{
				Name:        "Wine of Zamorak",
				AmountRange: [2]int{2, 10},
			},
		},
		uniqueDroptable: []Drop{
			{
				Name: "Hood of subjugation",
				Rate: 1.0 / 512.0,
				Bold: true,
			},
			{
				Name: "Garb of subjugation",
				Rate: 1.0 / 512.0,
				Bold: true,
			},
			{
				Name: "Gown of subjugation",
				Rate: 1.0 / 512.0,
				Bold: true,
			},
			{
				Name: "Gloves of subjugation",
				Rate: 1.0 / 512.0,
				Bold: true,
			},
			{
				Name: "Boots of subjugation",
				Rate: 1.0 / 512.0,
				Bold: true,
			},
			{
				Name: "Ward of subjugation",
				Rate: 1.0 / 512.0,
				Bold: true,
			},
			{
				Name: "Zamorakian spear",
				Rate: 1.0 / 512.0,
				Bold: true,
			},
			{
				Name: "Steam battlestaff",
				Rate: 1.0 / 512.0,
				Bold: true,
			},
		},
		extra: Drop{
			Name:   "Zamorak hilt",
			Amount: 1,
		},
		alwaysDroptable: []Drop{
			{
				Name:   "Infernal ashes",
				Amount: 1,
			},
		},
	}
)

var KrilCommand = &discordgo.ApplicationCommand{
	Name:        "kril",
	Description: "Simulate a General Kril drop",
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

func Kril(s *discordgo.Session, i *discordgo.InteractionCreate) {
	simulateDrop(s, i, simulateDropGwd1, krilTables, krilUrl, KrilCommand.Name, "K'ril Tsutsaroth")
}
