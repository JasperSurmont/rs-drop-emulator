package runescape

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	UncommonRateWithoutRare float64 = 1.0 / 3.0
	MAX_AMOUNT_ROLLS                = 500
)

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

func (d *Drop) SetAmount() {
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
func AmountToPrice(drops map[string]*Drop) (res []NamedRSPrice, total RSPrice, ok bool) {
	ok = true
	total = RSPrice("0")

	ch := make(chan NamedRSPrice)

	for _, d := range drops {
		go GetItemPrice(d.Name, ch, d.Untradeable)
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
		res = append(res, NamedRSPrice{
			Name:  d.Name,
			Price: price,
		})
		total.Add(price)
	}
	return
}

// Sort given map in decreasing value
func SortDrops(m *[]NamedRSPrice) {
	sort.Slice((*m)[:], func(i, j int) bool {
		comp, _ := (*m)[i].Price.Compare((*m)[j].Price)
		return comp == 1
	})
}

// Add the drops that should always be dropped to the drops pointer
func AddAlwaysDroptable(amount int64, drops *map[string]*Drop, alwaysDroptable []Drop) {
	for _, d := range alwaysDroptable {
		add := d // Add copy, otherwise it adjusts original value
		_, ok := (*drops)[add.Name]
		add.SetAmount()
		if ok {
			(*drops)[add.Name].Amount += add.Amount * int(amount)
		} else {
			add.Amount *= int(amount)
			(*drops)[add.Name] = &add
		}
	}
}

// Given the drops with values and drops with prices, make the drop list to be printed
func MakeDropList(n []NamedRSPrice, m map[string]*Drop, total RSPrice, ok bool) string {
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
func RemoveDropFromTable(s []Drop, i int) []Drop {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Get the option with the given name. If not found returns an empty option, so check if the returned option has eg a non-empty name
func GetOptionWithName(opt []*discordgo.ApplicationCommandInteractionDataOption, name string) discordgo.ApplicationCommandInteractionDataOption {
	for i, o := range opt {
		if o.Name == name {
			return *opt[i]
		}
	}
	return discordgo.ApplicationCommandInteractionDataOption{}
}

type RSPrice string

type NamedRSPrice struct {
	Name  string
	Price RSPrice
	Error error
}

// Returns 1 if p2 is smaller, 0 if equal and -1 if bigger.
func (p1 RSPrice) Compare(p2 RSPrice) (int, error) {
	p1o, p2o := p1[len(p1)-1:], p2[len(p2)-1:]               //order
	p1m, p2m := orderCharToFloat(p1o), orderCharToFloat(p2o) //multiplier

	p1v, p2v := string(p1[:len(p1)-1]), string(p2[:len(p2)-1]) //value

	// Adjust for numbers without an order: every character is part of the price
	if p1m == 1 {
		p1v = string(p1)
	}
	if p2m == 1 {
		p2v = string(p2)
	}

	p1f, err := strconv.ParseFloat(p1v, 64) //float
	if err != nil {
		return 0, err
	}
	p2f, err := strconv.ParseFloat(p2v, 64) //float
	if err != nil {
		return 0, err
	}

	p1r, p2r := p1f*p1m, p2f*p2m //real
	diff := p1r - p2r
	if diff < 0 {
		return -1, nil
	} else if diff > 0 {
		return 1, nil
	} else {
		return 0, nil
	}
}

// Format the price such that values that are too big for its order get formatted into a smaller value with bigger order
func (p *RSPrice) Format() (err error) {
	value, err := strconv.ParseFloat(string((*p)[:len(*p)-1]), 64)
	order := string((*p)[len(*p)-1:])

	// Adjust for numbers without an order: every character is part of the price
	if order != "b" && order != "m" && order != "k" {
		value, err = strconv.ParseFloat(string(*p), 64)
	}

	if err != nil {
		log.Errorf("couldn't format rsprice %v.\n%v", p, err)
		return
	}
	switch order {
	case "b": // Billion is the highest we can go
	case "m":
		if value >= 1000 {
			tmp := float64(value) / 100.0
			tmp = math.Floor(tmp)
			*p = RSPrice(fmt.Sprintf("%.1f", tmp/10) + "b")
		}
	case "k":
		if value >= 1000 {
			tmp := float64(value) / 100.0
			tmp = math.Floor(tmp)
			*p = RSPrice(fmt.Sprintf("%.1f", tmp/10) + "m")
			err = p.Format() // We format again in case it can be formatted further
		}
	default:
		if value >= 10000 { // It only gets converted to k once it hits 10k
			tmp := float64(value) / 100.0
			tmp = math.Floor(tmp)
			*p = RSPrice(fmt.Sprintf("%.1f", tmp/10) + "k")
			err = p.Format() // We format again in case it can be formatted further
		} else {
			*p = RSPrice(fmt.Sprintf("%.0f", value))
		}
	}
	return
}

// Multiply the number with x and format afterwards
func (p *RSPrice) Multiply(x int) error {
	value, err := strconv.ParseFloat(string((*p)[:len(*p)-1]), 64)
	order := (*p)[len(*p)-1:]

	// Adjust for numbers without an order: every character is part of the price
	if order != "b" && order != "m" && order != "k" {
		value, err = strconv.ParseFloat(string(*p), 64)
	}

	if err != nil {
		log.Errorf("couldn't format rsprice %v.\n%v", p, err)
		return err
	}

	multiplier := orderCharToFloat(order)
	*p = RSPrice(fmt.Sprintf("%f", value*multiplier*float64(x)))
	p.Format()

	return nil
}

// Add the number with x and format afterwards
func (p1 *RSPrice) Add(p2 RSPrice) error {
	p1o, p2o := (*p1)[len(*p1)-1:], p2[len(p2)-1:]           //order
	p1m, p2m := orderCharToFloat(p1o), orderCharToFloat(p2o) //multiplier

	p1v, p2v := string((*p1)[:len(*p1)-1]), string(p2[:len(p2)-1]) //value

	// Adjust for numbers without an order: every character is part of the price
	if p1m == 1 {
		p1v = string(*p1)
	}
	if p2m == 1 {
		p2v = string(p2)
	}

	p1f, err := strconv.ParseFloat(p1v, 64) //float
	if err != nil {
		return err
	}
	p2f, err := strconv.ParseFloat(p2v, 64) //float
	if err != nil {
		return err
	}

	p1r, p2r := p1f*p1m, p2f*p2m //real
	*p1 = RSPrice(fmt.Sprintf("%f", p1r+p2r))
	p1.Format()
	return nil
}

func (p *RSPrice) UnmarshalJSON(b []byte) error {
	var s string

	// We don't know if they return a string (everything above 9999gp) or an int (under 10k)
	if err := json.Unmarshal(b, &s); err != nil {
		var i int
		if err2 := json.Unmarshal(b, &i); err2 != nil {
			return fmt.Errorf("failed to unmarshal rsprice:\n%v\n%v", err, err2)
		}
		*p = RSPrice(fmt.Sprintf("%v", i))
		return nil
	}
	*p = RSPrice(strings.Replace(s, ",", "", -1))
	return nil
}

func orderCharToFloat(o RSPrice) float64 {
	switch o {
	case "b":
		return math.Pow10(9)
	case "m":
		return math.Pow10(6)
	case "k":
		return math.Pow10(3)
	default:
		return 1
	}
}

type ItemCacheEntry struct {
	LastUpdated string
	Price       RSPrice
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func determineDropWithRates(roll float64, drops *[]Drop) Drop {
	var totalSum float64 = 0
	for _, d := range *drops {
		totalSum += d.Rate
	}

	var sum float64 = 0
	for _, d := range *drops {
		sum += d.Rate / totalSum
		if roll < sum {
			return d
		}
	}

	// This should never happen, the sum will add up to 1
	return Drop{Name: "should not happen, check determineDropWithRates"}
}
