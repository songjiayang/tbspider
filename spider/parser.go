package spider

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

type Item struct {
	Title        string `json:"raw_title"`
	Link         string `json:"detail_url"`
	Price        string `json:"view_price"`
	Free         string `json:"view_fee"`
	Location     string `json:"item_loc"`
	CommentCount string `json:"comment_count"`

	ShopName string `json:"nick"`
	ShopLink string `json:"shopLink"`
}

var (
	itemsRegexp = regexp.MustCompile("auctions\":\\[\\{.*\"recommendAucti")
)

func parse(data []byte) (items []*Item, err error) {
	fdata := string(itemsRegexp.Find(data))
	if fdata == "" {
		log.Panic("Cookie expired, please reset again.")
	}

	fdata = strings.Replace(fdata, "auctions\":", "", 1)
	fdata = strings.Replace(fdata, ",\"recommendAucti", "", 1)

	if err = json.Unmarshal([]byte(fdata), &items); err != nil {
		return
	}

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
