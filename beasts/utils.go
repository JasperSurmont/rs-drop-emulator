package beasts

import (
	"math/rand"
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
