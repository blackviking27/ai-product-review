package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

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

	productUrlProvidedByUser := r.URL.Query().Get("url")

	if productUrlProvidedByUser == "" {
		slog.Error("No product URL found")
		ThrowErrorInResponse(w, http.StatusBadRequest, model.ERROR_NO_PRODUCT_URL, "No product url found")
		return
	}

	fmt.Println("Query Params: ", productUrlProvidedByUser)
	reviews, err := scraper.ScrapeDataFromURL(productUrlProvidedByUser)

	if err != nil {
		slog.Error("Unable to scrape data :" + err.Error())
		ThrowErrorInResponse(w, http.StatusInternalServerError, model.ERROR_FAILED_TO_SCRAPE_DATA, "Unable to scrape data :"+err.Error())
		return
	}

	scrappedResponse := model.ScrappedDataResponse{
		Name:       "Scrapped Data",
		ProductUrl: productUrlProvidedByUser,
		Reviews:    reviews,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(scrappedResponse)

}
