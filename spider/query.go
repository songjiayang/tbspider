package spider

import (
	"fmt"
	"net/url"
	"strconv"
)

type SortType int

func (st SortType) String() string {
	switch st {
	case 0:
		return "renqi-desc"
	case 1:
		return "sale-desc"
	case 2:
		return "credit-desc"
	case 3:
		return "price-asc"
	case 4:
		return "price-desc"
	default:
		return ""
	}
}

type Price float64

func (p Price) String() string {
	if p == 0 {
		return ""
	}

	return fmt.Sprintf("%0.2f", p)
}

type Query struct {
	Kw    string
	Loc   string
	SType SortType

	MinPrice Price
	MaxPrice Price

	IsTMall bool
	Limit   int
	Skip    int

	values url.Values
}

func (q *Query) IsValid() bool {
	return q.Kw != ""
}

func (q *Query) Values() url.Values {
	if q.values != nil {
		return q.values
	}

	v := url.Values{}

	v.Set("q", q.Kw)
	v.Set("sort", q.SType.String())

	if q.MinPrice > 0 || q.MaxPrice > 0 {
		v.Set("filter", fmt.Sprintf("reserve_price[%s,%s]", q.MinPrice, q.MaxPrice))
	}

	if q.Loc != "" {
		v.Set("loc", q.Loc)
	}

	if q.IsTMall {
		v.Set("tab", "mall")
	}

	q.values = v

	return v
}

func (q *Query) SetSkip(skip int) {
	q.Skip = skip
	q.values.Set("s", strconv.Itoa(skip))
}

func (q *Query) IsFinish() bool {
	return q.Skip > q.Limit
}
