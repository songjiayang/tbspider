package spider

import (
	"encoding/json"
	"regexp"
	"strings"
)

type Item struct {
	Title        string `json:"raw_title"`
	Link         string `json:"detail_url"`
	Price        string `json:"view_price"`
	Free         string `json:"view_fee"`
	Loction      string `json:"item_loc"`
	CommentCount string `json:"comment_count"`

	ShopName string `json:"nick"`
	ShopLink string `json:"shopLink"`
}

var (
	parseRegexp = regexp.MustCompile("auctions\":\\[\\{.*\"recommendAucti")
)

func parse(data []byte) (items []*Item, err error) {
	fdata := string(parseRegexp.Find(data))
	fdata = strings.Replace(fdata, "auctions\":", "", 1)
	fdata = strings.Replace(fdata, ",\"recommendAucti", "", 1)

	err = json.Unmarshal([]byte(fdata), &items)

	for _, item := range items {
		item.Link = prefixLink(item.Link)
		item.ShopLink = prefixLink(item.ShopLink)
	}

	return
}

func prefixLink(link string) string {
	if !strings.HasPrefix(link, "http") {
		link = "https:" + link
	}

	return link
}
