package core

import "math/rand"

func EmulateDropGwdVorago(amount int64, rareDroptable []Drop, commonDroptable []Drop, energyDrop Drop) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var sum float64 = 0
		for _, d := range rareDroptable {
			sum += d.Rate
		}
		rareDropped := false

		// Vorago has 5 drop piles
		for j := 0; j < 5; j++ {
			var drop Drop

			// Make copy so we don't adjust orig value
			cd := commonDroptable

			roll := rand.Float64()

			// Roll for energy
			if roll < energyDrop.Rate {
				// Make a copy so that we don't adjust the previous one
				addEnergy := energyDrop
				addEnergy.SetAmount()
				addDropValueToMap(drops, &addEnergy)
			}

			roll = rand.Float64()

			if roll < sum && !rareDropped {
				drop = rareDroptable[rand.Intn(len(rareDroptable))]
				rareDropped = true
			} else {
				drop = cd[rand.Intn(len(cd))]
			}

			drop.SetAmount()
			addDropValueToMap(drops, &drop)
		}

	}
	return drops
}
