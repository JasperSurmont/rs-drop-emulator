package core

import "math/rand"

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

func EmulateDropGwd1(amount int64, rareDroptable []Drop, uncommonDroptable []Drop, commonDroptable []Drop, hilt Drop) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var sum float64 = 0
		for _, d := range rareDroptable {
			sum += d.Rate
		}

		var drop Drop

		// We first roll the uniques
		if rand.Float64() < sum {
			drop = rareDroptable[rand.Intn(len(rareDroptable))]
		} else if rand.Float64() < 1.0/128.0 { // Hilt chance
			if rand.Float64() < 1.0/4.0 {
				drop = hilt
			} else {
				drop = Drop{
					Name:        "Coins",
					AmountRange: [2]int{19501, 21000},
				}
			}
		} else if rand.Float64() < 1.0/128.0 { // Shard chance
			if rand.Float64() < 1.0/2.0 {
				drop = godswordShards[rand.Intn(len(godswordShards))]
			} else {
				drop = Drop{
					Name:        "Coins",
					AmountRange: [2]int{20500, 21000},
				}
			}
		} else if rand.Float64() > CommonRateWithoutRare {
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
