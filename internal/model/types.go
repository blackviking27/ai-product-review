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

type ExtractedProductDetails struct {
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
	Data       ExtractedProductDetails `json:"data"`
}

const (
	ERROR_NO_PRODUCT_URL        = "ERROR_NO_PRODUCT_URL"
	ERROR_FAILED_TO_SCRAPE_DATA = "ERROR_FAILED_TO_SCRAPE_DATA"
)

type ApiReviewRequestBody struct {
	ApiKey     string `json:"apiKey" validate:"optional"`
	ProductUrl string `json:"productUrl" validate:"required"`
	AIModel    string `json:"aiModel" validate:"required"`
	AICompany  string `json:"aiCompany" validate:"optional"`
}

type ProductOpinion struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	OverallVerdict string   `json:"overall_verdict"`
	Pros           []string `json:"pros"`
	Cons           []string `json:"cons"`
	RatingScore    float64  `json:"rating_score"`
	Summary        string   `json:"executive_summary"`
}
