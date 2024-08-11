package scraping

import (
	"github.com/PuerkitoBio/goquery"
)

var ProductNameSpecs = map[string]WebsiteSelectorSpec{
	"TerabyteShop": {
		Selector: "h1.tit-prod",
	},
	"Picheu": {
		Selector: "title",
	},
}

func GetProductName(doc *goquery.Document) (string, error) {
	return GetTextWithSpec(doc, ProductNameSpecs)
}

func GetTextWithSpec(doc *goquery.Document, specs map[string]WebsiteSelectorSpec) (string, error) {
	// for each website, check if the document matches the price selector
	for _, spec := range specs {
		text := doc.Find(spec.Selector).Text()
		if text != "" {
			return text, nil
		}
	}

	return "", nil
}
