package core

import "math/rand"

// There are 8 uncommons, 3 commons and 1 rare
// let's say common is 50%, uncommon is 48%, then each uncommon has 6% chance, and rare has 2% chance
const (
	giantmoleRareRate     = 2.0 / 100.0
	giantmoleUncommonRate = 48.0 / 100
	giantmoleTertiaryRate = 1.0 / 10.0
)

func EmulateDropGiantMole(amount int64, tertiaryDroptable []Drop, rareDroptable []Drop, uncommonDroptable []Drop, commonDroptable []Drop) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var drop Drop

		roll := rand.Float64()

		// We split up the interval: [0, sum) = rare, [sum, sum+uncommonrate) = uncommon, [sum+uncommon, 1) = common
		if roll < giantmoleRareRate {
			drop = rareDroptable[rand.Intn(len(rareDroptable))]
		} else if roll < giantmoleRareRate+giantmoleUncommonRate {
			drop = uncommonDroptable[rand.Intn(len(uncommonDroptable))]
		} else {
			drop = commonDroptable[rand.Intn(len(commonDroptable))]
		}

		drop.SetAmount()
		addDropValueToMap(drops, &drop)

		// We roll for tertiary drop
		roll = rand.Float64()
		if roll < giantmoleTertiaryRate {
			roll = rand.Float64() // Roll again to decide which drop
			curr := 0.0

			for _, d := range tertiaryDroptable {
				curr += d.Rate
				if roll < curr {
					d.SetAmount()
					addDropValueToMap(drops, &d)
					break
				}
			}
		}

	}
	return drops
}
