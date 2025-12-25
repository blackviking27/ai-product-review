package scraper

import (
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
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

func ScrapteDataFromFlipkartUrl(url string, config config.Config) ([]model.Review, error) {

	htmlPageReader, err := GetProuductReviewHtml(url)
	if err != nil {
		slog.Error(err.Error())
		return []model.Review{}, err
	}
	defer htmlPageReader.Close()

	doc, err := goquery.NewDocumentFromReader(htmlPageReader)

	if err != nil {
		slog.Error("Unable to parse the product page")
		return []model.Review{}, err
	}

	var reviews []model.Review

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

	return reviews, nil
}
