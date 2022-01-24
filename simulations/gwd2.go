package simulations

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

// Simulate a drop with a rare table that has priority, and without a normal rare table
func simulateDropGwd2(amount int64, tables dropTables, i *discordgo.InteractionCreate) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var sum float64 = 0
		for _, d := range tables.uniqueDroptable {
			sum += d.Rate
		}

		var drop Drop

		roll := rand.Float64()

		// We split up the interval: [0, sum) = rare, [sum, sum+uncommonrate) = uncommon, [sum+uncommon, 1) = common
		if roll < sum {
			drop = determineDropWithRates(rand.Float64(), tables.uniqueDroptable)
		} else if roll < sum+UncommonRateWithoutRare {
			drop = tables.uncommonDroptable[rand.Intn(len(tables.uncommonDroptable))]
		} else {
			drop = tables.commonDroptable[rand.Intn(len(tables.commonDroptable))]
		}

		drop.setAmount()
		addDropValueToMap(drops, drop)
	}
	addGuarantees(amount, drops, tables.alwaysDroptable, i)
	return drops
}
