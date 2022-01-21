package simulations

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/jaspersurmont/rs-drop-simulator/logger"
	"github.com/jaspersurmont/rs-drop-simulator/rsapi"

	"github.com/bwmarrin/discordgo"
)

const (
	UncommonRateWithoutRare float64 = 1.0 / 3.0
	MAX_AMOUNT_ROLLS                = 500
)

var log logger.LoggerWrapper

func init() {
	log = logger.CreateLogger("simulations")
}

type Drop struct {
	Rate        float64
	Id          string
	Name        string
	Untradeable bool
	AmountRange [2]int
	Amount      int
	Amounts     []int  // If there are only certain specified amounts
	OtherDrops  []Drop // Other drops that always go with this drop
	Bold        bool   // Whether the drop should be put in bold or not
}

type dropTables struct {
	commonDroptable   []Drop
	uncommonDroptable []Drop
	rareDroptable     []Drop
	uniqueDroptable   []Drop
	alwaysDroptable   []Drop
	tertiaryDroptable []Drop
	extra             interface{} // Maybe find a better solution for this?
}

// Set the amount of a drop. If it already has an amount it does nothing. If it has multiple amounts it chooses one random. If it just has an
// AmountRange then it will randomly select a value in the closed interval
func (d *Drop) setAmount() {
	if d.Amount != 0 {
		return
	}
	if len(d.Amounts) > 0 {
		d.Amount = d.Amounts[rand.Intn(len(d.Amounts))]
		return
	}
	if d.AmountRange == [2]int{0, 0} {
		d.Amount = 1
	} else {
		diff := d.AmountRange[1] - d.AmountRange[0]
		d.Amount = d.AmountRange[0] + rand.Intn(diff+1) // We do +1 because it works with an open interval
	}
}

// Returns an array of structs with all the drops and their price
func amountToPrice(drops map[string]*Drop) (res []rsapi.NamedRSPrice, total rsapi.RSPrice, ok bool) {
	ok = true
	total = rsapi.RSPrice("0")

	ch := make(chan rsapi.NamedRSPrice)

	for _, d := range drops {
		go rsapi.GetItemPrice(d.Name, ch, d.Untradeable)
	}

	// We get same amount of values out of the channel, but continue if there's an error
	for range drops {
		namedPrice := <-ch
		if namedPrice.Error != nil {
			ok = false
			continue
		}

		price := namedPrice.Price
		d := drops[namedPrice.Name]

		err := price.Multiply(d.Amount)
		if err != nil {
			ok = false
			continue
		}
		res = append(res, rsapi.NamedRSPrice{
			Name:  d.Name,
			Price: price,
		})
		total.Add(price)
	}
	return
}

// Sort given map in decreasing value
func sortDrops(m *[]rsapi.NamedRSPrice) {
	sort.Slice((*m)[:], func(i, j int) bool {
		comp, _ := (*m)[i].Price.Compare((*m)[j].Price)
		return comp == 1
	})
}

// Add the drops that should always be dropped to the drops map
func addGuarantees(amount int64, drops *map[string]*Drop, alwaysDroptable []Drop, i *discordgo.InteractionCreate) {
	if enableGuarantees := getOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		for _, d := range alwaysDroptable {
			add := d // Add copy, otherwise it adjusts original value
			_, ok := (*drops)[add.Name]
			add.setAmount()
			if ok {
				(*drops)[add.Name].Amount += add.Amount * int(amount)
			} else {
				add.Amount *= int(amount)
				(*drops)[add.Name] = &add
			}
		}
	}
}

// Make a string in which all the drops are printed with their respective price
func makeDropList(n []rsapi.NamedRSPrice, m map[string]*Drop, total rsapi.RSPrice, ok bool) string {
	var sb strings.Builder
	for _, d := range n {
		if m[d.Name].Bold {
			sb.WriteString(fmt.Sprintf("**%v %v:** %v gp\n", m[d.Name].Amount, d.Name, d.Price))
		} else {
			sb.WriteString(fmt.Sprintf("%v %v: %v gp\n", m[d.Name].Amount, d.Name, d.Price))
		}
	}
	sb.WriteString(fmt.Sprintf("\n**Total GE value: %v**", total))
	if !ok {
		sb.WriteString("\nSomething went wrong; not all items were processed correctly.\n")
		sb.WriteString("If this issue keeps happening please file an issue. See /contribute for more info")
	}
	return sb.String()
}

func addDropValueToMap(m map[string]*Drop, d *Drop) {
	_, ok := m[d.Name]

	if ok {
		m[d.Name].Amount += d.Amount
	} else {
		m[d.Name] = d
	}
}

// Remove index i from the slice s
func removeDropFromTable(s []Drop, i int) []Drop {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Get the option with the given name. If not found returns an empty option, so check if the returned option has eg a non-empty name
func getOptionWithName(opt []*discordgo.ApplicationCommandInteractionDataOption, name string) discordgo.ApplicationCommandInteractionDataOption {
	for i, o := range opt {
		if o.Name == name {
			return *opt[i]
		}
	}
	return discordgo.ApplicationCommandInteractionDataOption{}
}

func determineDropWithRates(roll float64, drops []Drop) Drop {
	var totalSum float64 = 0
	for _, d := range drops {
		totalSum += d.Rate
	}

	var sum float64 = 0
	for _, d := range drops {
		sum += d.Rate / totalSum
		if roll < sum {
			return d
		}
	}

	// This should never happen, the sum will add up to 1
	return Drop{Name: "should not happen, check determineDropWithRates"}
}
