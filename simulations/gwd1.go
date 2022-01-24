package simulations

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var (
	godswordShards = []Drop{
		{
			Name:   "Godsword shard 1",
			Amount: 1,
		},
		{
			Name:   "Godsword shard 2",
			Amount: 1,
		},
		{
			Name:   "Godsword shard 3",
			Amount: 1,
		},
	}
)

func simulateDropGwd1(amount int64, tables dropTables, i *discordgo.InteractionCreate) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)
	hilt := tables.extra.(Drop)

	for i := int64(0); i < amount; i++ {
		var sum float64 = 0
		for _, d := range tables.uniqueDroptable {
			sum += d.Rate
		}

		var drop Drop

		roll := rand.Float64()

		// We split up the interval: [0, sum) = rare, [sum, sum+1/128) = Hilt chance, [sum+1/128,sum+1/128+1/128) = Shard chance,
		// [sum+1/128+1/128, sum+1/128+1/128+uncommonrate) = uncommon, [sum+1/128+1/128+uncommonrate,1) = common
		if roll < sum {
			drop = determineDropWithRates(rand.Float64(), tables.uniqueDroptable)
		} else if roll < sum+1.0/128.0 { // Hilt chance
			if rand.Float64() < 1.0/4.0 {
				drop = hilt
			} else {
				drop = Drop{
					Name:        "Coins",
					AmountRange: [2]int{19501, 21000},
				}
			}
		} else if roll < sum+1.0/128.0+1.0/128.0 { // Shard chance
			if rand.Float64() < 1.0/2.0 {
				drop = godswordShards[rand.Intn(len(godswordShards))]
			} else {
				drop = Drop{
					Name:        "Coins",
					AmountRange: [2]int{20500, 21000},
				}
			}
		} else if roll < sum+1.0/128.0+1.0/128.0+UncommonRateWithoutRare {
			drop = tables.uncommonDroptable[rand.Intn(len(tables.uncommonDroptable))]
		} else {
			drop = tables.commonDroptable[rand.Intn(len(tables.commonDroptable))]
		}

		drop.setAmount()
		addDropValueToMap(drops, drop)

		// Add drops that always go together
		for _, d := range drop.OtherDrops {
			d.setAmount()
			addDropValueToMap(drops, d)
		}
	}
	return drops
}
