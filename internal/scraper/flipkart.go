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

	// -----------------------------------------------------------------------
	// 1. EXTRACT PRODUCT NAME
	// Selector: h1.CEn5rD > span.LMizgS
	// -----------------------------------------------------------------------
	productSpecs.Name = strings.TrimSpace(doc.Find("h1.CEn5rD span.LMizgS").Text())

	// -----------------------------------------------------------------------
	// 2. EXTRACT SPECIFICATIONS
	// Structure: .QZKsWF (Section) -> .ZRVDNa (Title) -> table.n7infM (Table)
	// -----------------------------------------------------------------------
	doc.Find("div.QZKsWF").Each(func(i int, section *goquery.Selection) {
		// Section Title (e.g., "Processor And Memory Features")
		secTitle := strings.TrimSpace(section.Find("div.ZRVDNa").Text())
		if secTitle == "" {
			return
		}

		// if _, ok := productSpecs.TechSpecs[secTitle]; !ok {
		// 	specs = make(map[string]string)
		// }

		// Iterate rows in the table inside this section
		section.Find("table.n7infM tr.row").Each(func(j int, row *goquery.Selection) {
			// Key: td.JMeybS
			key := strings.TrimSpace(row.Find("td.JMeybS").Text())
			// Value: td.QPlg21 -> ul -> li
			val := strings.TrimSpace(row.Find("td.QPlg21 li").Text())

			if key != "" && val != "" {
				productSpecs.Specs[key] = val
			}
		})
	})

	// -----------------------------------------------------------------------
	// 3. EXTRACT DESCRIPTION
	// Structure: The description is split into multiple blocks with class .wV4AgS
	// Each block has a Title (.J9EKOQ) and Body (.QsvKzw p)
	// -----------------------------------------------------------------------
	var descBuilder strings.Builder

	// Find the container following "Product Description" label
	doc.Find("div.wV4AgS").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("div.J9EKOQ").Text())
		body := strings.TrimSpace(s.Find("div.QsvKzw p").Text())

		if title != "" || body != "" {
			if title != "" {
				descBuilder.WriteString(title + ": ")
			}
			descBuilder.WriteString(body + "\n\n")
		}
	})

	productSpecs.Description = []string{strings.TrimSpace(descBuilder.String())}
	return productSpecs
}

func ScrapteDataFromFlipkartUrl(url string, config config.Config) (model.ExtractedProdictDetails, error) {

	var productDetailsAndReviews model.ExtractedProdictDetails

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
