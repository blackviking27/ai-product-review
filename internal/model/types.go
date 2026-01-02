package model

type Review struct {
	Text   string
	Rating string
}

type ProductSpecs struct {
	Name        string
	Specs       map[string]string
	Description []string
}

type ExtractedProdictDetails struct {
	Review []Review
	Specs  ProductSpecs
}

type ErrorResponse struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type ScrappedDataRequest struct {
	Url string `json:"url"`
}

type ScrappedDataResponse struct {
	Name       string                  `json:"name"`
	ProductUrl string                  `json:"productUrl"`
	Data       ExtractedProdictDetails `json:"data"`
}

const (
	ERROR_NO_PRODUCT_URL        = "ERROR_NO_PRODUCT_URL"
	ERROR_FAILED_TO_SCRAPE_DATA = "ERROR_FAILED_TO_SCRAPE_DATA"
)

type ApiReviewRequestBody struct {
	ApiKey     string `json:"apiKey"`
	ProductUrl string `json:"productUrl"`
}

type ProductVerdict struct {
	Verdict     string
	Rating      string
	ProductName string
	ProductType string
}
