package beasts

import (
	"math/rand"
)

const (
	commonRateWithoutRare float32 = 2.0 / 3.0
	commonRateWithRare    float32 = 6.0 / 10.0
	uncommonRateWithRare  float32 = 3.0 / 10.0
)

type Drop struct {
	Rate        float32
	Id          string
	Name        string
	AmountRange [2]int
	Amount      int
}

func (d *Drop) SetAmount() {
	if d.AmountRange == [2]int{0, 0} {
		d.Amount = 1
	} else {
		diff := d.AmountRange[1] - d.AmountRange[0]
		d.Amount = d.AmountRange[0] + rand.Intn(diff+1) // We do +1 because it works with an open interval
	}
}

// Emulate a drop with a rare table that has priority, and without a normal rare table
func EmulateDrop(commonRateWithoutRare float32, amount int64, rareDroptable []Drop, uncommonDroptable []Drop, commonDroptable []Drop) map[string]*Drop {

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

func AddAlwaysDroptable(amount int64, drops *map[string]*Drop, alwaysDroptable []Drop) {
	for _, d := range alwaysDroptable {
		_, ok := (*drops)[d.Name]
		if ok {
			(*drops)[d.Name].Amount += d.Amount * int(amount)
		} else {
			d.Amount *= int(amount)
			(*drops)[d.Name] = &d
		}
	}
}
