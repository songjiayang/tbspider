package spider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
)

type Spider struct {
	query *Query

	items    []*Item
	isFinish bool
}

func NewSpider(query *Query) *Spider {
	return &Spider{
		query: query,
		items: make([]*Item, 0),
	}
}

func (this *Spider) Result() []*Item {
	return this.items
}

func (this *Spider) Run() error {

	for !this.isFinish {
		body, err := this.fetch()
		if err != nil {
			return err
		}

		items, err := parse(body)
		if err != nil {
			return err
		}

		this.items = append(this.items, items...)

		this.query.SetSkip(this.query.Skip + len(items))
		this.isFinish = (len(items) == 0 || this.query.IsFinish())

		time.Sleep(time.Second)
	}

	return nil
}

func (this *Spider) fetch() (body []byte, err error) {
	req, err := http.NewRequest("GET", this.page(), nil)
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

func (this *Spider) page() string {
	return fmt.Sprintf("https://s.taobao.com/search?%s", this.query.Values().Encode())
}
