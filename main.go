package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/songjiayang/tbspider/spider"
)

var (
	s                  int
	minPrice, maxPrice float64
)

func main() {
	q := &spider.Query{}

	flag.StringVar(&q.Kw, "q", "", "查询关键词")
	flag.StringVar(&q.Loc, "loc", "", "发货地")
	flag.IntVar(&s, "sort", 0, "商品排序，可选参数: 0 人气, 1 销量，2 信用，3 价格由低到高，4 价格由高到低")
	flag.Float64Var(&minPrice, "min", 0, "最低价格")
	flag.Float64Var(&maxPrice, "max", 0, "最高价格")
	flag.IntVar(&q.Limit, "l", 200, "最多获取商品数量")
	flag.BoolVar(&q.IsTMall, "mall", false, "是否只抓取天猫数据，默认抓取淘宝数据")
	flag.Parse()

	q.SType = spider.SortType(s)
	q.MinPrice = spider.Price(minPrice)
	q.MaxPrice = spider.Price(maxPrice)

	if !q.IsValid() {
		panic("Invalid args")
	}

	worker := spider.NewSpider(q)

	err := worker.Run()
	if err != nil {
		panic(err)
	}

	save(worker.Result(), q.Kw)
}

func save(items []*spider.Item, kw string) error {
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	t := time.Now()
	filename := fmt.Sprintf("./%s-%d%02d%02d.json", kw, t.Year(), t.Month(), t.Day())

	return ioutil.WriteFile(filename, data, 0644)
}
