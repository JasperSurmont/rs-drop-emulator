package simulations

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var (
	rakshaUrl    = "https://runescape.wiki/images/Raksha%2C_the_Shadow_Colossus.png?fba12"
	rakshaTables = dropTables{
		commonDroptable: []Drop{
			{
				Name:        "Spirit weed seed",
				AmountRange: [2]int{5, 10},
			},
			{
				Name:   "Carambola seed",
				Amount: 5,
			},
			{
				Name:        "Golden dragonfruit seed",
				AmountRange: [2]int{5, 10},
			},
			{
				Name:   "Small blunt rune salvage",
				Amount: 50,
			},
			{
				Name:        "Medium spiky orikalkum salvage",
				AmountRange: [2]int{12, 18},
			},
			{
				Name:        "Huge plated orikalkum salvage",
				AmountRange: [2]int{10, 15},
			},
			{
				Name:        "Black dragonhide",
				AmountRange: [2]int{200, 300},
			},
			{
				Name:        "Onyx dust", // Since it's only 1 uncommon, we don't make a diff
				AmountRange: [2]int{75, 131},
			},
			{
				Name:        "Dinosaur bones",
				AmountRange: [2]int{80, 120},
			},
			{
				Name:        "Crystal key",
				AmountRange: [2]int{25, 45},
			},
			{
				Name:        "Inert adrenaline crystal",
				AmountRange: [2]int{12, 18},
			},
			{
				Name:   "Sirenic scale",
				Amount: 3,
			},
			{
				Name:        "Soul rune",
				AmountRange: [2]int{125, 175},
			},
			{
				Name:        "Dark animica stone spirit",
				AmountRange: [2]int{60, 95},
				OtherDrops: []Drop{
					{
						Name:       "Light animica stone spirit",
						SameAmount: true,
					},
				},
			},
		},
		uniqueDroptable: []Drop{
			{
				Name: "Laceration boots",
				Rate: 1.0 / 200.0,
				Bold: true,
			},
			{
				Name: "Blast diffusion boots",
				Rate: 1.0 / 200.0,
				Bold: true,
			},
			{
				Name: "Fleeting boots",
				Rate: 1.0 / 200.0,
				Bold: true,
			},
			{
				Name: "Shadow spike",
				Rate: 1.0 / 500.0,
				Bold: true,
			},
			{
				Name: "Greater Ricochet ability codex",
				Rate: 1.0 / 500.0,
				Bold: true,
			},
			{
				Name: "Greater Chain ability codex",
				Rate: 1.0 / 500.0,
				Bold: true,
			},
			{
				Name: "Divert ability codex",
				Rate: 1.0 / 500.0,
				Bold: true,
			},
		},
	}
	RakshaCommand = &discordgo.ApplicationCommand{
		Name:        "raksha",
		Description: "Simulate a raksha drop",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of kills",
				Required:    false,
			},
		},
	}
)

func Raksha(s *discordgo.Session, i *discordgo.InteractionCreate) {
	simulateDrop(s, i, simulateDropRaksha, rakshaTables, rakshaUrl, RakshaCommand.Name, "Raksha")
}

func simulateDropRaksha(amount int64, tables dropTables, i *discordgo.InteractionCreate) map[string]*Drop {
	var drops = make(map[string]*Drop)

	var sum float64 = 0
	for _, rare := range tables.uniqueDroptable {
		sum += rare.Rate
	}

	for i := int64(0); i < amount; i++ {
		var drop Drop

		roll := rand.Float64() // Check if the rare is rolled
		if roll < sum {
			drop = determineDropWithRates(rand.Float64(), tables.uniqueDroptable)
			drop.setAmount()
			addDropValueToMap(drops, drop)
		}

		drop = tables.commonDroptable[rand.Intn(len(tables.commonDroptable))]
		drop.setAmount()
		addDropValueToMap(drops, drop)
	}

	return drops
}
