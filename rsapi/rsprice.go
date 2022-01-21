package rsapi

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jaspersurmont/rs-drop-simulator/logger"
)

var log = logger.CreateLogger("rsapi")

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
		log.Error(fmt.Sprintf("couldn't format rsprice %v.%v", p, err))
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
		log.Error(fmt.Sprintf("couldn't format rsprice %v.%v", p, err))
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

func (p *RSPrice) unmarshalJSON(b []byte) error {
	var s string

	// We don't know if they return a string (everything above 9999gp) or an int (under 10k)
	if err := json.Unmarshal(b, &s); err != nil {
		var i int
		if err2 := json.Unmarshal(b, &i); err2 != nil {
			return fmt.Errorf("failed to unmarshal rsprice:\n%v", err2)
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

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
