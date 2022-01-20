package runescape

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "embed"

	zd "github.com/blendle/zapdriver"
)

const DETAIL_URL = "https://services.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item="

var idList = make(map[string]int)
var itemCache = make(map[int]ItemCacheEntry)

var itemCacheMutex = &sync.RWMutex{}

//go:embed id-list.json
var initIdlistData []byte

// Load the IDlist in the map
func init() {
	err := json.Unmarshal(initIdlistData, &idList)
	if err != nil {
		log.Fatalf("couldn't unmarshal idlist into map, %v", err)
	}
}

// Look up the item price of name
// Logging is already done, use err solely to check if it succeeded or not
func GetItemPrice(name string, ch chan<- NamedRSPrice, untradeable bool) {
	if untradeable {
		ch <- NamedRSPrice{
			Name:  name,
			Price: "0",
		}
		return
	}
	lower := strings.ToLower(name)
	if lower == "coin" || lower == "coins" {
		ch <- NamedRSPrice{
			Name:  "Coins",
			Price: "1",
		}
		return
	}
	id, ok := idList[name]
	if !ok {
		ch <- NamedRSPrice{
			Error: errors.New("name not found in idlist"),
		}
		log.Errorf("name not found in idlist %v", name)
		return
	}
	price, err := GetItemPriceById(id)
	ch <- NamedRSPrice{
		Name:  name,
		Price: price,
		Error: err,
	}
}

// Look up the item price of an item id
// Logging is already done, use err solely to check if it succeeded or not
func GetItemPriceById(id int) (price RSPrice, err error) {

	// First we check the cache
	itemCacheMutex.RLock()
	item, ok := itemCache[id]
	itemCacheMutex.RUnlock()
	if ok {
		if then, err2 := time.Parse(time.UnixDate, item.LastUpdated); err == nil && DateEqual(then, time.Now()) {
			price = item.Price
			return
		} else {
			log.Error("couldn't parse time from cache",
				zd.Label("id", fmt.Sprint(id)),
				zd.Label("time", item.LastUpdated),
				zd.Label("err", err2.Error()),
			)
		}
	}

	res, err := http.Get(DETAIL_URL + strconv.Itoa(id))
	if err != nil {
		log.Error("request failed in GetItemPriceById",
			zd.Label("id", fmt.Sprint(id)),
			zd.Label("err", err.Error()),
		)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("could not read body in GetItemPriceById",
			zd.Label("id", fmt.Sprint(id)),
			zd.Label("err", err.Error()),
		)
		return
	}

	if res.StatusCode > 299 {
		log.Error("invalid response in GetItemPriceById",
			zd.Label("id", fmt.Sprint(id)),
			zd.Label("statusCode", fmt.Sprint(res.StatusCode)),
		)
		err = errors.New("invalid response in GetItemPriceById")
		return
	}

	var detail DetailResponse
	if err := json.Unmarshal(body, &detail); err != nil {
		log.Error("couldn't unmarshal json to DetailResponse in GetItemPriceById",
			zd.Label("id", fmt.Sprint(id)),
			zd.Label("error", err.Error()),
		)
	}
	price = detail.Item.Current.Price

	itemCacheMutex.Lock()
	itemCache[id] = ItemCacheEntry{
		Price:       price,
		LastUpdated: time.Now().Format(time.UnixDate),
	}
	itemCacheMutex.Unlock()
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
