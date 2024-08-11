package scraping

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/felipemarinho97/price-crawler/requester"
)

func pichauCashPricePostProcessor(rq *requester.Requester, price string) string {
	sku := strings.Split(price, "SKU: ")[1]
	fmt.Println(price, sku)

	pichauResponse, shouldReturn, returnValue := getPichauCatalog(rq, sku, price)
	if shouldReturn {
		return returnValue
	}
	fmt.Printf("pichauResponse: %v\n", pichauResponse)

	return pichauResponse["data"].(map[string]interface{})["productDetail"].(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})["pichau_prices"].(map[string]interface{})["avista"].(string)
}

func pichauCreditPricePostProcessor(rq *requester.Requester, price string) string {
	sku := strings.Split(price, "SKU: ")[1]

	pichauResponse, shouldReturn, returnValue := getPichauCatalog(rq, sku, price)
	if shouldReturn {
		return returnValue
	}

	return pichauResponse["data"].(map[string]interface{})["productDetail"].(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})["pichau_prices"].(map[string]interface{})["final_price"].(string)
}

func getPichauCatalog(rq *requester.Requester, sku string, price string) (map[string]interface{}, bool, string) {
	data := map[string]interface{}{
		"operationName": "productDetail",
		"variables": map[string]string{
			"sku": sku,
		},
		"query": `query productDetail($sku: String) {\n  productDetail: products(filter: {sku: {eq: $sku}}) {\n    items {\n      __typename\n      sku\n      name\n      only_x_left_in_stock\n      stock_status\n      special_price\n      mysales_promotion {\n        expire_at\n        price_discount\n        price_promotional\n        promotion_name\n        promotion_url\n        qty_available\n        qty_sold\n        __typename\n      }\n      pichauUlBenchmarkProduct {\n        overallScore\n        scoreCPU\n        scoreGPU\n        games {\n          fullHdFps\n          medium4k\n          quadHdFps\n          title\n          ultra1080p\n          ultra4k\n          __typename\n        }\n        __typename\n      }\n      pichau_prices {\n        avista\n        avista_discount\n        avista_method\n        base_price\n        final_price\n        max_installments\n        min_installment_price\n        __typename\n      }\n      price_range {\n        __typename\n      }\n      ... on SimpleProduct {\n        options {\n          option_id\n          required\n          title\n          sort_order\n          __typename\n          ... on CustomizableRadioOption {\n            value {\n              price\n              price_type\n              sku\n              uid\n              title\n              option_type_id\n              __typename\n            }\n            __typename\n          }\n          ... on CustomizableMultipleOption {\n            value {\n              price\n              price_type\n              sku\n              uid\n              title\n              option_type_id\n              __typename\n            }\n            __typename\n          }\n          ... on CustomizableDropDownOption {\n            value {\n              price\n              price_type\n              sku\n              uid\n              title\n              option_type_id\n              __typename\n            }\n            __typename\n          }\n          ... on CustomizableCheckboxOption {\n            value {\n              price\n              price_type\n              sku\n              uid\n              title\n              option_type_id\n              __typename\n            }\n            __typename\n          }\n        }\n        __typename\n      }\n    }\n    __typename\n  }\n}\n`,
	}

	resp, err := rq.Post("https://www.pichau.com.br/api/catalog", data)
	if err != nil {
		fmt.Printf("Error on getting pichau catalog: %v\n", err)
		return nil, true, price
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error on reading pichau response: %v\n", err)
		return nil, true, price
	}
	fmt.Printf("content: %s\n", string(content))

	var pichauResponse map[string]interface{}
	err = json.Unmarshal(content, &pichauResponse)
	if err != nil {
		fmt.Printf("Error on decoding pichau response: %v\n", err)
		return nil, true, ""
	}
	fmt.Printf("pichauResponse: %v\n", pichauResponse)
	return pichauResponse, false, ""
}
