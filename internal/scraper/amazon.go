package scraper

import (
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blackviking27/ai-product-reviwer/config"
	"github.com/blackviking27/ai-product-reviwer/internal/model"
)

// // This is will return the product review page
// // which is hidden behind a login page
// func GetAmazonReviewLinkPath(url url.URL, config config.Config) (string, error) {

// 	re := regexp.MustCompile(`\/dp\/([a-zA-Z0-9]{10})`)

// 	matchedString := re.FindStringSubmatch(url.Path)

// 	if len(matchedString) == 0 {
// 		slog.Error("Unable to extract product Id from url")
// 		return "", errors.New("Unable to extract product Id from url")
// 	}

// 	fullProductReviewLink := fmt.Sprintf(url.Scheme+"://"+url.Host+"/product-reviews/%s/ref=cm_cr_dp_d_show_all_btm?ie=UTF8&reviewerType=all_reviews", matchedString[1])

// 	return fullProductReviewLink, nil
// }

func extractDetailsFromPage(htmlPageReader io.Reader) (model.ExtractedProdictDetails, error) {

	var extractedDetail model.ExtractedProdictDetails

	doc, err := goquery.NewDocumentFromReader(htmlPageReader)

	if err != nil {
		slog.Error("Unable to parse the html page for scraping")
		return extractedDetail, err
	}

	// var reviews []model.Review

	reviews := getProductReviewFromDocForAamazon(doc)
	productSpecs := getProductDetailsFromDocForAmazon(doc)

	extractedDetail.Review = reviews
	extractedDetail.Specs = productSpecs

	return extractedDetail, nil
}

func getProductReviewFromDocForAamazon(doc *goquery.Document) (reviews []model.Review) {
	doc.Find(`div[id^="customer_review"]`).Each(func(i int, s *goquery.Selection) {

		reviewRating := strings.TrimSpace(s.Find(`i[data-hook="review-star-rating"] span.a-icon-alt`).Text())
		if reviewRating == "" {
			reviewRating = strings.TrimSpace(s.Find(`i.a-icon-star span.a-icon-alt`).Text())
		}
		bodyNode := s.Find(`[data-hook="review-body"] span`)
		bodyNode.Find("br").ReplaceWithHtml("\n")
		reviewBody := strings.TrimSpace(bodyNode.Text())

		reviews = append(reviews, model.Review{
			Text:   strings.TrimSpace(reviewBody),
			Rating: reviewRating,
		})
	})
	return reviews
}

func getProductDetailsFromDocForAmazon(doc *goquery.Document) model.ProductSpecs {

	specs := model.ProductSpecs{
		Name:        "",
		Specs:       map[string]string{},
		Description: []string{},
	}

	// Product Name
	specs.Name = strings.TrimSpace(doc.Find("#productTitle").Text())

	// doc, err := goquery.NewDocumentFromReader(htmlPageReader)

	// if err != nil {
	// 	slog.Error("Unable to parse the html page for scraping")
	// 	return specs, err
	// }

	// EXTRACT PRODUCT OVERVIEW
	// doc.Find("#productOverview_feature_div tr").Each(func(i int, s *goquery.Selection) {
	// 	key := strings.TrimSpace(s.Find("td.a-span3 span").Text())
	// 	val := strings.TrimSpace(s.Find("td.a-span9 span").Text())

	// 	if key != "" && val != "" {
	// 		specs.Overview[key] = val
	// 	}
	// })

	// EXTRACT TECHNICAL SPECIFICATIONS
	doc.Find("#productDetails_techSpec_section_1 tr").Each(func(i int, s *goquery.Selection) {
		// Key is in <th>, Value is in <td>
		key := strings.TrimSpace(s.Find("th").Text())
		val := strings.TrimSpace(s.Find("td").Text())

		// Clean up keys (remove invisible characters)
		key = strings.ReplaceAll(key, "\n", "")
		val = strings.ReplaceAll(val, "\n", "")
		val = strings.ReplaceAll(val, "\u200e", "") // Remove Left-to-Right marks

		if key != "" {
			specs.Specs[key] = val
		}
	})

	// EXTRACT BULLET POINTS
	doc.Find("#feature-bullets li span.a-list-item").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			specs.Description = append(specs.Description, text)
		}
	})

	return specs
}

func ScrapteDataFromAmazonUrl(url string, config config.Config) (model.ExtractedProdictDetails, error) {

	htmlPageReader, err := GetProuductReviewHtml(url)

	if err != nil {
		slog.Error(err.Error())
		return model.ExtractedProdictDetails{}, err
	}

	extractedDetails, err := extractDetailsFromPage(htmlPageReader)

	if err != nil {
		slog.Error("Error while scraping data")
		return extractedDetails, errors.New("Error while scraping data")
	}

	if len(extractedDetails.Review) == 0 {
		slog.Error("No reviews found for product")
		return extractedDetails, errors.New("No reviews found for product")

	}

	return extractedDetails, nil

}
