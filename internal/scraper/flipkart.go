package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blackviking27/ai-product-reviwer/config"
	"github.com/blackviking27/ai-product-reviwer/internal/model"
)

func GetFlipkartReviewLinkPath(url url.URL, config config.Config) (string, error) {
	if !strings.Contains(url.String(), "/p") {
		slog.Error("Not a valid product URL")
		return "", fmt.Errorf("Not a valid product URL")
	}

	matches := regexp.MustCompile(`\/p\/(itm[a-zA-Z0-9]+)`).FindStringSubmatch(url.String())
	if len(matches) > 1 {
		fmt.Println("Found Product ID:", matches[1])
	}

	reviewUrl := strings.Replace(url.String(), "/p/", "/product-reviews/", 1)

	return reviewUrl, nil
}

func getProductReviewFromDocforFlipkart(doc *goquery.Document) (reviews []model.Review) {
	doc.Find("div.col.x_CUu6.QccLnz").Each(func(i int, s *goquery.Selection) {

		ratingStr := s.Find("div.MKiFS6").Text()
		rating := ratingStr

		title := s.Find("p.qW2QI1").Text()

		reviewBody := s.Find("div.G4PxIA").Text()

		// Cleanup: Remove "READ MORE" which Flipkart appends to long reviews
		reviewBody = strings.ReplaceAll(reviewBody, "READ MORE", "")
		reviewBody = strings.TrimSpace(reviewBody)

		if reviewBody != "" {
			// Combine Title and Body for better AI Context
			fullText := title + ". " + reviewBody

			reviews = append(reviews, model.Review{
				Text:   fullText,
				Rating: rating,
			})
		}
	})

	return reviews
}

func getProductDetailsFromDocForFlipkart(doc *goquery.Document) (productSpecs model.ProductSpecs) {

	// Extract the JSON Script
	// The data is inside <script id="is_script">window.__INITIAL_STATE__ = { ... }</script>
	scriptContent := doc.Find("script#is_script").Text()
	if scriptContent == "" {
		log.Fatal("Could not find the script tag containing window.__INITIAL_STATE__")
	}

	// Clean the string to get valid JSON
	jsonStr := strings.TrimSpace(scriptContent)
	jsonStr = strings.TrimPrefix(jsonStr, "window.__INITIAL_STATE__ = ")
	jsonStr = strings.TrimSuffix(jsonStr, ";")

	// Parse JSON flexibly
	// We use map[string]interface{} because the keys (like "10001", "10002") are dynamic
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Navigate the JSON structure to find Product Data
	// Path: pageDataV4 -> page -> data -> [Dynamic Slot ID] -> [Widget List] -> widget -> data -> product -> value
	pageData, _ := data["pageDataV4"].(map[string]interface{})
	page, _ := pageData["page"].(map[string]interface{})
	slots, _ := page["data"].(map[string]interface{})

	var productName string
	productSpecs.Specs = make(map[string]string)

	// Iterate over the slots (e.g., "10001", "10002")
	for _, slotContent := range slots {
		widgets, ok := slotContent.([]interface{})
		if !ok {
			continue
		}

		for _, w := range widgets {
			widgetMap, _ := w.(map[string]interface{})
			widgetDef, _ := widgetMap["widget"].(map[string]interface{})
			widgetType, _ := widgetDef["type"].(string)

			// We are looking for the 'PRODUCT_MIN' widget which contains the summary
			if widgetType == "PRODUCT_MIN" {
				wData, _ := widgetDef["data"].(map[string]interface{})
				product, _ := wData["product"].(map[string]interface{})
				value, _ := product["value"].(map[string]interface{})

				// Extract Title
				if titles, ok := value["titles"].(map[string]interface{}); ok {
					productName, _ = titles["title"].(string)
					productSpecs.Name = productName
				}

				// Extract Key Specs
				if specs, ok := value["keySpecs"].([]interface{}); ok {
					for i, s := range specs {
						productSpecs.Specs[strconv.Itoa(i)] = s.(string)
					}
				}
			}
		}
	}

	return productSpecs
}

func ScrapteDataFromFlipkartUrl(url string, config config.Config) (model.ExtractedProductDetails, error) {

	var productDetailsAndReviews model.ExtractedProductDetails

	htmlPageReader, err := GetProuductReviewHtml(url)
	if err != nil {
		slog.Error(err.Error())
		return productDetailsAndReviews, err
	}
	defer htmlPageReader.Close()

	doc, err := goquery.NewDocumentFromReader(htmlPageReader)

	if err != nil {
		slog.Error("Unable to parse the product page")
		return productDetailsAndReviews, err
	}

	productDetailsAndReviews.Review = getProductReviewFromDocforFlipkart(doc)
	productDetailsAndReviews.Specs = getProductDetailsFromDocForFlipkart(doc)

	// var reviews []model.Review

	// doc.Find("div.col.x_CUu6.QccLnz").Each(func(i int, s *goquery.Selection) {

	// 	ratingStr := s.Find("div.MKiFS6").Text()
	// 	rating := ratingStr

	// 	title := s.Find("p.qW2QI1").Text()

	// 	reviewBody := s.Find("div.G4PxIA").Text()

	// 	// Cleanup: Remove "READ MORE" which Flipkart appends to long reviews
	// 	reviewBody = strings.ReplaceAll(reviewBody, "READ MORE", "")
	// 	reviewBody = strings.TrimSpace(reviewBody)

	// 	if reviewBody != "" {
	// 		// Combine Title and Body for better AI Context
	// 		fullText := title + ". " + reviewBody

	// 		reviews = append(reviews, model.Review{
	// 			Text:   fullText,
	// 			Rating: rating,
	// 		})
	// 	}
	// })

	return productDetailsAndReviews, nil
}
