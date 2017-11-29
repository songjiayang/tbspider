package spider

import (
	"testing"

	"github.com/golib/assert"
)

func TestSpiderPage(t *testing.T) {
	assertion := assert.New(t)

	worker := NewSpider(&Query{
		Kw: "袜子",
	})

	assertion.Equal(worker.page(), "https://s.taobao.com/search?q=%E8%A2%9C%E5%AD%90&sort=renqi-desc")
}

func TestSpiderRun(t *testing.T) {
	assertion := assert.New(t)

	worker := NewSpider(&Query{
		Kw: "袜子",
	})

	err := worker.Run()

	assertion.Nil(err)
	assertion.NotEmpty(worker.Result())
}
