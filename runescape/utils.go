package runescape

import "time"

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

type ItemCacheEntry struct {
	LastUpdated string
	Price       RSPrice
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
