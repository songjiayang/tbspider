package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
)

var (
	q, loc, minPrice, maxPrice string
	s, limit                   int
	mall                       bool

	sortMap map[int]string = map[int]string{
		0: "renqi-desc",
		1: "sale-desc",
		2: "credit-desc",
		3: "price-asc",
		4: "price-desc",
	}
)

func main() {
	flag.StringVar(&q, "q", "", "查询关键词")
	flag.StringVar(&loc, "loc", "", "发货地")
	flag.IntVar(&s, "sort", 0, "商品排序，可选参数: 0 人气, 1 销量，2 信用，3 价格由低到高，4 价格由高到低")
	flag.StringVar(&minPrice, "min", "", "最低价格")
	flag.StringVar(&maxPrice, "max", "", "最高价格")
	flag.IntVar(&limit, "l", 200, "最多获取商品数量")
	flag.BoolVar(&mall, "mall", false, "是否为抓取天猫数据，默认只抓取淘宝数据")
	flag.Parse()

	if q == "" {
		panic("Invalid args")
	}

	var (
		finish     bool
		skip       = 0
		foundItems []string
		page       = "https://s.taobao.com/search"
	)

	for !finish {
		page = buildUrl(page, skip)
		body, err := fetch(page)
		if err != nil {
			panic(err)
		}

		items, err := parse(body)
		if err != nil {
			panic(err)
		}

		foundItems = append(foundItems, items...)

		skip += len(items)
		if len(items) == 0 || skip > limit {
			finish = true
		}

		time.Sleep(time.Second)
	}

	err := save(foundItems)
	if err != nil {
		panic(err)
	}
}

func buildUrl(page string, skip int) string {
	u, _ := url.Parse(page)

	query := u.Query()

	if skip == 0 {
		query.Set("q", q)
		query.Set("sort", sortMap[s])

		if minPrice != "" || maxPrice != "" {
			query.Set("filter", fmt.Sprintf("reserve_price[%s,%s]", minPrice, maxPrice))
		}

		if loc != "" {
			query.Set("loc", loc)
		}

		if mall {
			query.Set("tab", "mall")
		}
	}

	query.Set("s", strconv.Itoa(skip))

	u.RawQuery = query.Encode()

	return u.String()
}

func fetch(page string) (body []byte, err error) {
	req, err := http.NewRequest("GET", page, nil)
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", randomdata.UserAgentString())
	req.Header.Add("Host", randomdata.IpV4Address())
	req.Header.Add("Referer", "https://s.taobao.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	return
}

var (
	parseRex = regexp.MustCompile("auctions\":\\[\\{.*\"recommendAucti")
)

func parse(data []byte) (items []string, err error) {
	fdata := string(parseRex.Find(data))
	fdata = strings.Replace(fdata, "auctions\":", "", 1)
	fdata = strings.Replace(fdata, ",\"recommendAucti", "", 1)

	var list []map[string]interface{}

	err = json.Unmarshal([]byte(fdata), &list)
	if err != nil {
		return
	}

	items = make([]string, len(list))

	for i, item := range list {
		itemUrl := item["detail_url"].(string)
		if !strings.HasPrefix(itemUrl, "http") {
			itemUrl = "https:" + itemUrl
		}
		items[i] = itemUrl
	}

	return
}

func save(items []string) error {
	t := time.Now()
	filename := fmt.Sprintf("./%d-%02d-%02d.txt", t.Year(), t.Month(), t.Day())
	return ioutil.WriteFile(filename, []byte(strings.Join(items, "\n")), 0644)
}
