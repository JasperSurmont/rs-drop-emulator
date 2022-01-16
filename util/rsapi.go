package util

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strconv"
)

const DETAIL_URL = "https://services.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item="

func GetItemPrice(name string) {
	// look in id-list here
}

func GetItemPriceById(id int) (price string, err error) {
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
			"body", body,
		)
		err = errors.New("invalid response in GetItemPriceById")
		return
	}

	var detail DetailResponse
	if err := json.Unmarshal(body, &detail); err != nil {
		log.Errorw("couldn't unmarshal json to DetailResponse in GetItemPriceById",
			"id", id,
			"body", body,
		)
	}
	price = detail.Item.Current.Price
	return
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
			Trend string `json:"trend"`
			Price string `json:"price"`
		} `json:"current"`
		Today struct {
			Trend string `json:"trend"`
			Price string `json:"price"`
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

// Returns 1 if p2 is smaller, 0 if equal and -1 if bigger.
func (p1 RSPrice) Compare(p2 RSPrice) (int, error) {
	p1o, p2o := p1[len(p1)-1:], p2[len(p2)-1:]               //order
	p1m, p2m := orderCharToFloat(p1o), orderCharToFloat(p2o) //multiplier

	p1v, p2v := string(p1[:len(p1)-1]), string(p2[:len(p2)-1]) //value
	p1f, err := strconv.ParseFloat(p1v, 64)                    //float
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
