# tbspider
A simple spider for taobao products, write with golang.

### Usage

1. Use `go get github.com/songjiayang/tbspider/spider`

```golang

package main

import (
	"fmt"

	"github.com/songjiayang/tbspider/spider"
)

func main() {
	q := &spider.Query{
		Kw: "袜子",
	}
	
	worker := spider.NewSpider(q, "./cookie")

	err := worker.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println(worker.Result())
}

```

2. Or use [cli](https://github.com/songjiayang/tbspider/releases/tag/v0.1) command tool

```bash
$ tbspider -h

Usage of tbspider:
  -cookie string
    	淘宝登录cookie (default "./cookie")
  -l int
    	最多获取商品数量 (default 200)
  -loc string
    	发货地
  -mall
    	是否只抓取天猫数据，默认抓取淘宝数据
  -max float
    	最高价格
  -min float
    	最低价格
  -q string
    	查询关键词
  -sort int
    	商品排序，可选参数: 0 人气, 1 销量，2 信用，3 价格由低到高，4 价格由高到低
    	
```

when you run `tbspider -q 袜子`, you will get data in filename `袜子-2017xxxx.json`.

### Data struct

``` golang
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

```


