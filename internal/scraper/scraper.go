package scraper

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/blackviking27/ai-product-reviwer/config"
	"github.com/blackviking27/ai-product-reviwer/internal/model"
)

func getFullProdcutURL(url url.URL, config config.Config) (string, error) {

	host := url.Host
	var err error
	updatedReviewUrl := ""
	switch host {
	case config.Scrapper.Amazon.ReviewLink.Host:
		updatedReviewUrl = url.String()
	case config.Scrapper.Flipkart.ReviewLink.Host:
		updatedReviewUrl, err = GetFlipkartReviewLinkPath(url, config)
	}

	if err != nil {
		return "", err
	}

	if updatedReviewUrl == "" {
		return "", errors.New("No product URL found")
	}

	return updatedReviewUrl, nil
}

// Used by scrapper to get the product HTML page
func GetProuductReviewHtml(url string) (io.ReadCloser, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// html, err := http.Get(url)
	res, err := client.Do(req)

	if err != nil {
		slog.Error("Unable to fetch the product review page")
		return nil, err
	}

	return res.Body, nil

}

func ScrapeDataFromURL(productLink string) ([]model.Review, error) {
	config, err := config.LoadConfig()
	parsedProductLink, err := url.Parse(productLink)
	if err != nil {
		slog.Error("Unable to parse the product link")
		return nil, errors.New("Unable to parse the product link")
	}

	fullProductReviewLink, err := getFullProdcutURL(*parsedProductLink, *config)

	if err != nil {
		slog.Error(err.Error())
		return []model.Review{}, err
	}

	fullProductLinkUrlObject, err := url.Parse(fullProductReviewLink)

	if err != nil {
		slog.Error("Unable to parse full product review link: " + err.Error())
		return nil, err
	}

	var productReviews []model.Review

	switch fullProductLinkUrlObject.Host {
	case config.Scrapper.Amazon.ReviewLink.Host:
		productReviews, err = ScrapteDataFromAmazonUrl(fullProductLinkUrlObject.String(), *config)
	case config.Scrapper.Flipkart.ReviewLink.Host:
		productReviews, err = ScrapteDataFromFlipkartUrl(fullProductLinkUrlObject.String(), *config)
	}

	if err != nil || len(productReviews) == 0 {
		slog.Error("No reviews found for the product")
		return nil, err
	}

	return productReviews, nil
}
