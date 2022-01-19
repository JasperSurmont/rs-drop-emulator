package runescape

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

		roll := rand.Float64()

		// We split up the interval: [0, sum) = rare, [sum, sum+1/128) = Hilt chance, [sum+1/128,sum+1/128+1/128) = Shard chance,
		// [sum+1/128+1/128, sum+1/128+1/128+uncommonrate) = uncommon, [sum+1/128+1/128+uncommonrate,1) = common
		if roll < sum {
			drop = rareDroptable[rand.Intn(len(rareDroptable))]
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
			drop = uncommonDroptable[rand.Intn(len(uncommonDroptable))]
		} else {
			drop = commonDroptable[rand.Intn(len(commonDroptable))]
		}

		drop.SetAmount()
		addDropValueToMap(drops, &drop)

		// Add drops that always go together
		for _, d := range drop.OtherDrops {
			d.SetAmount()
			addDropValueToMap(drops, &d)
		}
	}
	return drops
}
