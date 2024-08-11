package scraping

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/felipemarinho97/price-crawler/requester"
)

type Price struct {
	CashPrice   string `json:"cashPrice"`
	CreditPrice string `json:"creditPrice"`
}

var (
	CashPriceSpecs = map[string]WebsiteSelectorSpec{
		"TerabyteShop": {
			Selector: "#valVista",
		},
		"Pichau": {
			Selector:    "div.jss104",
			PostProcess: pichauCashPricePostProcessor,
		},
	}
	CreditPriceSpecs = map[string]WebsiteSelectorSpec{
		"TerabyteShop": {
			Selector: "#valParc",
		},
		"Pichau": {
			Selector:    "div.jss104",
			PostProcess: pichauCreditPricePostProcessor,
		},
	}
)

func GetPriceWithSpec(doc *goquery.Document, rq *requester.Requester, specs map[string]WebsiteSelectorSpec) (string, error) {
	// for each website, check if the document matches the price selector
	for _, spec := range specs {
		price := doc.Find(spec.Selector).Text()
		if price != "" {
			if spec.PostProcess != nil {
				price = spec.PostProcess(rq, price)
			}
			return price, nil
		}
	}

	return "", fmt.Errorf("failed to get price for all specs: %v", specs)
}

func GetPrice(doc *goquery.Document, rq *requester.Requester) (Price, error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var cashPrice, creditPrice string
	var cashErr, creditErr error

	go func() {
		defer wg.Done()
		cashPrice, cashErr = GetPriceWithSpec(doc, rq, CashPriceSpecs)
	}()

	go func() {
		defer wg.Done()
		creditPrice, creditErr = GetPriceWithSpec(doc, rq, CreditPriceSpecs)
	}()

	wg.Wait()

	if cashErr != nil {
		return Price{}, fmt.Errorf("failed to get cash price: %w", cashErr)
	}
	if creditErr != nil {
		return Price{}, fmt.Errorf("failed to get credit price: %w", creditErr)
	}

	return Price{
		CashPrice:   fmt.Sprintf("%.2f", NormalizePrice(cashPrice)),
		CreditPrice: fmt.Sprintf("%.2f", NormalizePrice(creditPrice)),
	}, nil
}

var priceRegex = regexp.MustCompile(`([0-9](,|.))+`)

func NormalizePrice(price string) float64 {
	// find the first non-zero match
	matches := priceRegex.FindAllString(price, -1)
	if len(matches) == 0 {
		return 0
	}

	// for all matches, find the first non-zero match
	for _, match := range matches {
		// replace dot with nothing
		match = strings.ReplaceAll(match, ".", "")

		// replace comma with dot
		match = strings.ReplaceAll(match, ",", ".")

		// parse the match
		price, err := strconv.ParseFloat(match, 64)
		if err != nil {
			return 0
		}

		// if the price is non-zero, return it
		if price != 0 {
			return price
		}
	}

	return 0
}
