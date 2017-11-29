package spider

import (
	"testing"

	"github.com/golib/assert"
)

func TestSortType(t *testing.T) {
	assertion := assert.New(t)

	assertion.Equal(SortType(0).String(), "renqi-desc")
	assertion.Equal(SortType(1).String(), "sale-desc")
	assertion.Equal(SortType(2).String(), "credit-desc")
	assertion.Equal(SortType(3).String(), "price-asc")
	assertion.Equal(SortType(4).String(), "price-desc")
}

func TestPrice(t *testing.T) {
	assertion := assert.New(t)

	assertion.Equal(Price(0).String(), "")
	assertion.Equal(Price(1.0).String(), "1.00")
}

func TestQueryValues(t *testing.T) {
	assertion := assert.New(t)

	query := &Query{
		Kw:       "袜子",
		Loc:      "上海",
		SType:    SortType(1),
		MinPrice: 1.0,
		MaxPrice: 200,
		IsTMall:  true,
	}

	v := query.Values()

	assertion.Equal(v.Get("q"), "袜子")
	assertion.Equal(v.Get("loc"), "上海")
	assertion.Equal(v.Get("sort"), "sale-desc")
	assertion.Equal(v.Get("filter"), "reserve_price[1.00,200.00]")
	assertion.Equal(v.Get("tab"), "mall")
}
