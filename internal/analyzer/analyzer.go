package analyzer

import (
	"context"
	"errors"
	"strings"

	"github.com/blackviking27/ai-product-reviwer/internal/model"
)

func isValidParams(params model.ApiReviewRequestBody) []string {
	var errors []string

	if params.ProductUrl == "" {
		errors = append(errors, "No product URL provided")
	}
	return errors

}
func GetAiReviewForProduct(params model.ApiReviewRequestBody) (model.ProductOpinion, error) {
	ctx := context.Background()

	// Validating the params
	validationErrors := isValidParams(params)

	if len(validationErrors) > 0 {
		return model.ProductOpinion{}, errors.New(strings.Join(validationErrors, ", "))
	}

	// Add the different model selection logic
	var opinion model.ProductOpinion
	var err error
	switch params.AICompany {
	case "gemini":
	case "google":
		opinion, err = getReviewFromGemini(ctx, params)
	}

	return opinion, err
}
