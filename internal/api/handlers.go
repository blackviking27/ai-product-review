package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/blackviking27/ai-product-reviwer/internal/analyzer"
	"github.com/blackviking27/ai-product-reviwer/internal/model"
	"github.com/blackviking27/ai-product-reviwer/internal/scraper"
)

func ThrowErrorInResponse(w http.ResponseWriter, code int, name, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := model.ErrorResponse{
		Name: name,
		Text: msg,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Error while encoding data" + err.Error())
		http.Error(w, msg, code)
	}

}

func GetScrapedDataForProduct(w http.ResponseWriter, r *http.Request) {

	var reqBody model.ScrappedDataRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		slog.Error("Error while parsing the request body for parsing request body")
		ThrowErrorInResponse(w, http.StatusBadRequest, "INTERNAL_SERVER_ERROR", "Error while parsing the request body for parsing request body")
		return
	}

	productUrlProvidedByUser := reqBody.Url

	if productUrlProvidedByUser == "" {
		slog.Error("No product URL found")
		ThrowErrorInResponse(w, http.StatusBadRequest, model.ERROR_NO_PRODUCT_URL, "No product url found")
		return
	}

	extractedDetail, err := scraper.ScrapeDataFromURL(productUrlProvidedByUser)

	if err != nil {
		slog.Error("Unable to scrape data :" + err.Error())
		ThrowErrorInResponse(w, http.StatusInternalServerError, model.ERROR_FAILED_TO_SCRAPE_DATA, "Unable to scrape data :"+err.Error())
		return
	}

	scrappedResponse := model.ScrappedDataResponse{
		Name:       "Scrapped Data",
		ProductUrl: productUrlProvidedByUser,
		Data:       extractedDetail,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(scrappedResponse)

}

func GetAiReviewForProduct(w http.ResponseWriter, r *http.Request) {

	var requestBody model.ApiReviewRequestBody

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		slog.Error("Missing/Incorrect fields being passed ")
		ThrowErrorInResponse(w, http.StatusBadRequest, "Missing fields", "Required fields are missing in the request body")
		return
	}

	opinion, err := analyzer.GetAiReviewForProduct(requestBody)
	if err != nil {
		ThrowErrorInResponse(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(opinion)

}
