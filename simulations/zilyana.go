package simulations

import (
	"github.com/bwmarrin/discordgo"
)

var (
	zilyanaUrl    = "https://runescape.wiki/images/General_Graardor.png?c6b33"
	zilyanaTables = dropTables{
		commonDroptable: []Drop{
			{
				Name:   "Grimy ranarr",
				Amount: 5,
			},
			{
				Name:   "Ranarr seed",
				Amount: 2,
			},
			{
				Name:        "Rune dart",
				AmountRange: [2]int{30, 40},
			},
			{
				Name:   "Large plated rune salvage",
				Amount: 1,
			},
		},
		uncommonDroptable: []Drop{
			{
				Name:   "Diamond",
				Amount: 6,
			},
			{
				Name:   "Prayer potion (4)",
				Amount: 3,
			},
			{
				Name:   "Saradomin brew (3)",
				Amount: 3,
				OtherDrops: []Drop{
					{
						Name:   "Super restore (4)",
						Amount: 3,
					},
				},
			},
			{
				Name:   "Super defence (3)",
				Amount: 3,
				OtherDrops: []Drop{
					{
						Name:   "Super magic potion (3)",
						Amount: 3,
					},
				},
			},
			{
				Name:        "Unicorn horn",
				AmountRange: [2]int{5, 10},
			},
			{
				Name:   "Battlestaff",
				Amount: 2,
			},
			{
				Name:   "Huge plated adamant salvage",
				Amount: 1,
			},
			{
				Name:   "Magic seed",
				Amount: 1,
			},
		},
		uniqueDroptable: []Drop{
			{
				Name: "Saradomin sword",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
			{
				Name: "Armadyl crossbow",
				Rate: 1.0 / 768.0,
				Bold: true,
			},
			{
				Name: "Off-hand Armadyl crossbow",
				Rate: 1.0 / 768.0,
				Bold: true,
			},
			{
				Name: "Saradomin's murmur",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
			{
				Name: "Saradomin's hiss",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
			{
				Name: "Saradomin's whisper",
				Rate: 1.0 / 384.0,
				Bold: true,
			},
		},
		extra: Drop{
			Name:   "Saradomin hilt",
			Amount: 1,
		},
		alwaysDroptable: []Drop{
			{
				Name:   "Bones",
				Amount: 1,
			},
		},
	}
)

var ZilyanaCommand = &discordgo.ApplicationCommand{
	Name:        "zilyana",
	Description: "Simulate a Command Zilyana drop",
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

func Zilyana(s *discordgo.Session, i *discordgo.InteractionCreate) {
	simulateDrop(s, i, simulateDropGwd1, zilyanaTables, zilyanaUrl, ZilyanaCommand.Name, "Commander Zilyana")
}
