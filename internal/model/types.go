package model

type Review struct {
	Text   string
	Rating string
}

type ErrorResponse struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type ScrappedDataResponse struct {
	Name       string   `json:"name"`
	ProductUrl string   `json:"productUrl"`
	Reviews    []Review `json:"reviews"`
}

const (
	ERROR_NO_PRODUCT_URL        = "ERROR_NO_PRODUCT_URL"
	ERROR_FAILED_TO_SCRAPE_DATA = "ERROR_FAILED_TO_SCRAPE_DATA"
)
