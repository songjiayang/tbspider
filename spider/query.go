package spider

import (
	"fmt"
	"net/url"
	"strconv"
)

type SortType int

func (this SortType) String() string {
	switch this {
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

func (this Price) String() string {
	if this == 0 {
		return ""
	}

	return fmt.Sprintf("%0.2f", this)
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

func (this *Query) IsValid() bool {
	return this.Kw != ""
}

func (this *Query) Values() url.Values {
	if this.values != nil {
		return this.values
	}

	v := url.Values{}

	v.Set("q", this.Kw)
	v.Set("sort", this.SType.String())

	if this.MinPrice > 0 || this.MaxPrice > 0 {
		v.Set("filter", fmt.Sprintf("reserve_price[%s,%s]", this.MinPrice, this.MaxPrice))
	}

	if this.Loc != "" {
		v.Set("loc", this.Loc)
	}

	if this.IsTMall {
		v.Set("tab", "mall")
	}

	this.values = v

	return v
}

func (this *Query) SetSkip(skip int) {
	this.Skip = skip
	this.values.Set("s", strconv.Itoa(skip))
}

func (this *Query) IsFinish() bool {
	return this.Skip > this.Limit
}
