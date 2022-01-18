package core

import "math/rand"

// Emulate a drop with a rare table that has priority, and without a normal rare table
func EmulateDropGwd2(commonRateWithoutRare float32, amount int64, rareDroptable []Drop, uncommonDroptable []Drop, commonDroptable []Drop) map[string]*Drop {

	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var sum float32 = 0
		for _, d := range rareDroptable {
			sum += d.Rate
		}

		var drop Drop

		if rand.Float32() < sum {
			drop = rareDroptable[rand.Intn(len(rareDroptable))]
		} else if rand.Float32() > commonRateWithoutRare {
			drop = uncommonDroptable[rand.Intn(len(uncommonDroptable))]
		} else {
			drop = commonDroptable[rand.Intn(len(commonDroptable))]
		}

		drop.SetAmount()
		_, ok := drops[drop.Name]

		if ok {
			drops[drop.Name].Amount += drop.Amount
		} else {
			drops[drop.Name] = &drop
		}
	}
	return drops
}
