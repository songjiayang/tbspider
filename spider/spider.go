package spider

import (
	"errors"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Spider struct {
	query *Query

	headers  http.Header
	items    []*Item
	isFinish bool
}

func NewSpider(query *Query, cookieFile string) *Spider {
	// get default headers
	headers := http.Header{}
	headers.Add("Referer", "https://s.taobao.com")
	// set cookie jar
	cookie, err := ioutil.ReadFile(cookieFile)
	if err != nil {
		log.Panicf("load cookie jar with error: %v \n", err)
	}
	headers.Add("cookie", string(cookie))

	return &Spider{
		query:   query,
		headers: headers,
		items:   make([]*Item, 0),
	}
}

func (s *Spider) Result() []*Item {
	return s.items
}

func (s *Spider) Run() error {
	for !s.isFinish {
		body, err := s.fetch()
		if err != nil {
			return err
		}

		items, err := parse(body)
		if err != nil {
			return err
		}

		s.items = append(s.items, items...)

		l := len(items)
		s.query.SetSkip(s.query.Skip + l)
		s.isFinish = (l == 0 || s.query.IsFinish())

		time.Sleep(3 * time.Second)
	}

	return nil
}

func (s *Spider) fetch() (body []byte, err error) {
	req, err := http.NewRequest("GET", s.page(), nil)
	if err != nil {
		return
	}

	req.Header = s.headers
	req.Header.Set("User-Agent", randomdata.UserAgentString())
	req.Header.Set("Host", randomdata.IpV4Address())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, errors.New("Cookie expired")
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *Spider) page() string {
	return fmt.Sprintf("https://s.taobao.com/search?%s", s.query.Values().Encode())
}
