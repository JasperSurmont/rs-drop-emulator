package simulations

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var (
	nexUrl    = "https://runescape.wiki/images/Nex.png?67c8a"
	nexTables = dropTables{
		commonDroptable: []Drop{
			{
				Name:   "Magic logs",
				Amount: 375,
			},
			{
				Name:   "Phasmatite stone spirit",
				Amount: 20,
			},
			{
				Name:   "Necrite stone spirit",
				Amount: 20,
			},
			{
				Name:   "Green dragonhide",
				Amount: 400,
			},
			{
				Name:   "Uncut dragonstone",
				Amount: 20,
			},
		},
		uncommonDroptable: []Drop{
			{
				Name:   "Onyx bolts (e)",
				Amount: 375,
			},
			{
				Name:   "Grimy avantoe",
				Amount: 75,
				OtherDrops: []Drop{
					{
						Name:   "Grimy dwarf weed",
						Amount: 75,
					},
				},
			},
			{
				Name:   "Grimy torstol",
				Amount: 40,
			},
			{
				Name:   "Torstol seed",
				Amount: 12,
			},
			{
				Name:   "Magic seed",
				Amount: 5,
			},
		},
		uniqueDroptable: []Drop{
			{
				Rate: 1.0 / 384.0,
				Name: "Torva full helm",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Torva platebody",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Torva platelegs",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Torva boots",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Torva gloves",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Pernix cowl",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Pernix body",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Pernix chaps",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Pernix boots",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Pernix gloves",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Virtus mask",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Virtus robe top",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Virtus robe legs",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Virtus boots",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Virtus gloves",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Virtus wand",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Virtus book",
				Bold: true,
			},
			{
				Rate: 1.0 / 384.0,
				Name: "Zaryte bow",
				Bold: true,
			},
		},
		extra: []Drop{
			{
				Name:   "Saradomin brew (4)",
				Amount: 10,
				Rate:   7.0 / 128.0,
				OtherDrops: []Drop{
					{
						Name:   "Super restore (4)",
						Amount: 30,
					},
				},
			},
			{
				Name:   "Super restore (4)",
				Amount: 10,
				Rate:   29.0 / 128.0,
				OtherDrops: []Drop{
					{
						Name:   "Saradomin brew (4)",
						Amount: 30,
					},
				},
			},
		},
		alwaysDroptable: []Drop{
			{
				Name:   "Big bones",
				Amount: 1,
			},
		},
	}
)

var NexCommand = &discordgo.ApplicationCommand{
	Name:        "nex",
	Description: "Simulate a Nex drop with full reputation",
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

func Nex(s *discordgo.Session, i *discordgo.InteractionCreate) {
	simulateDrop(s, i, simulateDropNex, nexTables, nexUrl, NexCommand.Name, "Nex")
}

func simulateDropNex(amount int64, tables dropTables, i *discordgo.InteractionCreate) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	nexPotions := tables.extra.([]Drop)

	for i := int64(0); i < amount; i++ {
		var sum float64 = 0
		for _, d := range tables.uniqueDroptable {
			sum += d.Rate
		}

		var drop Drop

		roll := rand.Float64()

		// We split up the interval: [0, sum) = rare, [sum, sum+uncommonrate) = uncommon, [sum+uncommon, 1) = common
		if roll < sum {
			drop = tables.uniqueDroptable[rand.Intn(len(tables.uniqueDroptable))]
		} else if roll < sum+UncommonRateWithoutRare {
			drop = tables.uncommonDroptable[rand.Intn(len(tables.uncommonDroptable))]
		} else {
			// Compute if we have a potion drop, if not just common drop
			var sum2 float64 = 0
			for _, d := range nexPotions {
				sum2 += d.Rate
			}
			if roll2 := rand.Float64(); roll2 < sum2 {
				drop = determineDropWithRates(rand.Float64(), nexPotions)
			} else {
				drop = tables.commonDroptable[rand.Intn(len(tables.commonDroptable))]
			}
		}

		drop.setAmount()
		addDropValueToMap(drops, &drop)

		// Add drops that always go together
		for _, d := range drop.OtherDrops {
			d.setAmount()
			addDropValueToMap(drops, &d)
		}
	}
	addGuarantees(amount, &drops, tables.alwaysDroptable, i)
	return drops
}
