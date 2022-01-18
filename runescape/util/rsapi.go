package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	_ "embed"
)

const DETAIL_URL = "https://services.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item="

var idList = make(map[string]int)

//go:embed id-list.json
var initData []byte

// Load the IDlist in the map
func init() {
	err := json.Unmarshal(initData, &idList)
	if err != nil {
		log.Fatalf("couldn't unmarshal idlist into map, %v", err)
	}
}

// Look up the item price of name
// Logging is already done, use err solely to check if it succeeded or not
func GetItemPrice(name string) (price RSPrice, err error) {
	lower := strings.ToLower(name)
	if lower == "coin" || lower == "coins" {
		price = "1"
		return
	}
	id, ok := idList[name]
	if !ok {
		err = errors.New("name not found in idlist")
		return
	}
	price, err = GetItemPriceById(id)
	return
}

// Look up the item price of an item id
// Logging is already done, use err solely to check if it succeeded or not
func GetItemPriceById(id int) (price RSPrice, err error) {
	res, err := http.Get(DETAIL_URL + strconv.Itoa(id))
	if err != nil {
		log.Errorw("request failed in GetItemPriceById",
			"id", id,
			"err", err,
		)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorw("could not read body in GetItemPriceById",
			"id", id,
			"err", err,
		)
		return
	}

	if res.StatusCode > 299 {
		log.Errorw("invalid response in GetItemPriceById",
			"id", id,
			"statusCode", res.StatusCode,
		)
		err = errors.New("invalid response in GetItemPriceById")
		return
	}

	var detail DetailResponse
	if err := json.Unmarshal(body, &detail); err != nil {
		log.Errorw("couldn't unmarshal json to DetailResponse in GetItemPriceById",
			"id", id,
			"error", err,
		)
	}
	price = detail.Item.Current.Price
	return
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

type DetailResponse struct {
	Item struct {
		Icon        string `json:"icon"`
		IconLarge   string `json:"icon_large"`
		ID          int    `json:"id"`
		Type        string `json:"type"`
		TypeIcon    string `json:"typeIcon"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Current     struct {
			Trend string  `json:"trend"`
			Price RSPrice `json:"price"`
		} `json:"current"`
		Today struct {
			Trend string  `json:"trend"`
			Price RSPrice `json:"price"`
		} `json:"today"`
		Members string `json:"members"`
		Day30   struct {
			Trend  string `json:"trend"`
			Change string `json:"change"`
		} `json:"day30"`
		Day90 struct {
			Trend  string `json:"trend"`
			Change string `json:"change"`
		} `json:"day90"`
		Day180 struct {
			Trend  string `json:"trend"`
			Change string `json:"change"`
		} `json:"day180"`
	} `json:"item"`
}

type RSPrice string
