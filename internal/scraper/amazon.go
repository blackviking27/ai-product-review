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

func getReviewFromAPage(htmlPageReader io.Reader) ([]model.Review, error) {

	doc, err := goquery.NewDocumentFromReader(htmlPageReader)

	if err != nil {
		slog.Error("Unable to parse the html page for scraping")
		return nil, err
	}

	var reviews []model.Review

	// reviewContainer := doc.Find("#cm-cr-dp-review-list")
	// reviewContainer.Find("div[data-hook='review']").Each(func(i int, s *goquery.Selection) {
	// 	// text := s.Find("span[data-hook='review-body']").Text()
	// 	// rating := s.Find("i[data-hook='review-star-rating']").Text() // e.g., "5.0 out of 5 stars"

	// 	// Title
	// 	// The title text is usually in a span inside the link
	// 	// reviewTitle = strings.TrimSpace(s.Find(`a[data-hook="review-title"] span`).Last().Text())

	// 	reviewRating := s.Find(`i[data-hook="review-star-rating"] .a-icon-alt`).Text()
	// 	// if review.Rating == "" {
	// 	// 	review.Rating = s.Find("i.a-icon-star .a-icon-alt").Text()
	// 	// }

	// 	// Date
	// 	// review.Date = strings.TrimSpace(s.Find(`[data-hook="review-date"]`).Text())

	// 	// Body
	// 	// Replace <br> tags with newlines to preserve formatting
	// 	bodySel := s.Find(`[data-hook="review-body"] span`)
	// 	bodySel.Find("br").ReplaceWithHtml("\n")
	// 	reviewBody := strings.TrimSpace(bodySel.Text())

	// 	reviews = append(reviews, model.Review{
	// 		Text:   strings.TrimSpace(reviewBody),
	// 		Rating: reviewRating,
	// 	})
	// })

	doc.Find(`div[id^="customer_review"]`).Each(func(i int, s *goquery.Selection) {
		// ID: Extract from the attribute ID (e.g., "customer_review-R3...")
		// review.ID, _ = s.Attr("id")

		// Author: Find the profile name
		// review.Author = strings.TrimSpace(s.Find(".a-profile-name").First().Text())

		// Rating:
		// Priority 1: The 'data-hook' selector
		// Priority 2: The generic 'a-icon-star' class
		reviewRating := strings.TrimSpace(s.Find(`i[data-hook="review-star-rating"] span.a-icon-alt`).Text())
		if reviewRating == "" {
			reviewRating = strings.TrimSpace(s.Find(`i.a-icon-star span.a-icon-alt`).Text())
		}

		// Title:
		// The title is often inside a link tag. We want the text of the span *inside* that link.
		// review.Title = strings.TrimSpace(s.Find(`a[data-hook="review-title"] span`).Last().Text())
		// if review.Title == "" {
		// 	// Fallback for some layouts where title is just bold text
		// 	review.Title = strings.TrimSpace(s.Find(`.review-title`).Text())
		// }

		// Date:
		// review.Date = strings.TrimSpace(s.Find(`[data-hook="review-date"]`).Text())

		// Body:
		// We target the inner span to avoid getting extra whitespace from the container.
		// We also replace <br> with newlines to keep paragraph formatting.
		bodyNode := s.Find(`[data-hook="review-body"] span`)
		bodyNode.Find("br").ReplaceWithHtml("\n")
		reviewBody := strings.TrimSpace(bodyNode.Text())

		// Only add valid reviews (skip empty placeholders)
		// if reviewBody != "" {
		// 	reviews = append(reviews, review)
		// }

		reviews = append(reviews, model.Review{
			Text:   strings.TrimSpace(reviewBody),
			Rating: reviewRating,
		})
	})

	return reviews, nil
}

func ScrapteDataFromAmazonUrl(url string, config config.Config) ([]model.Review, error) {

	var allReviews []model.Review

	htmlPageReader, err := GetProuductReviewHtml(url)

	if err != nil {
		slog.Error(err.Error())
		return []model.Review{}, err
	}

	pageReview, err := getReviewFromAPage(htmlPageReader)

	if err != nil {
		slog.Error("Error while scraping data")
		return []model.Review{}, errors.New("Error while scraping data")
	}

	if len(pageReview) == 0 {
		slog.Error("No reviews found for product")

		return nil, errors.New("No reviews found for product")

	}

	allReviews = append(allReviews, pageReview...)

	// for {
	// 	currPageUrl := url + strconv.Itoa(pageNumber)

	// 	pageReview, err := getReviewFromAPage(htmlPageReader)

	// 	if err != nil {
	// 		slog.Error("Error while scraping data")
	// 		return []model.Review{}, errors.New("Error while scraping data")
	// 	}

	// 	allReviews = append(allReviews, pageReview...)
	// 	break
	// }

	return allReviews, nil

}
