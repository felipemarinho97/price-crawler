package scraping

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/felipemarinho97/price-crawler/requester"
)

type Search struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Name      string  `json:"name"`
	CashPrice float64 `json:"cashPrice"`
	Link      string  `json:"link"`
}

var (
	SearchSpecs = map[string]WebsiteSelectorSpec{
		"TerabyteShop": {
			Selector: "#prodarea > div",
		},
	}
)

func GetSearch(doc *goquery.Document, rq *requester.Requester) (Search, error) {
	// get all divs inside the search results
	for _, spec := range SearchSpecs {
		divs := doc.Find(spec.Selector)
		if divs.Length() > 0 {
			search := Search{}
			search.Results = make([]SearchResult, divs.Length())

			divs.Each(func(i int, s *goquery.Selection) {
				search.Results[i].Name = s.Find(".prod-name").AttrOr("title", s.Find(".prod-name").Text())
				search.Results[i].CashPrice = NormalizePrice(s.Find(".prod-new-price").Text())
				search.Results[i].Link, _ = s.Find("a").Attr("href")
			})

			return search, nil
		}
	}

	return Search{}, nil
}
