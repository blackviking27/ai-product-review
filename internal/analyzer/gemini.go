package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/blackviking27/ai-product-reviwer/internal/model"
	"github.com/blackviking27/ai-product-reviwer/internal/scraper"
	"google.golang.org/genai"
)

func getGeminiContentConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"name":              {Type: genai.TypeString},
				"description":       {Type: genai.TypeString},
				"overall_verdict":   {Type: genai.TypeString},
				"pros":              {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
				"cons":              {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
				"rating_score":      {Type: genai.TypeNumber},
				"executive_summary": {Type: genai.TypeString},
			},
		},
	}
}

func generatePrompt(reviews model.ExtractedProductDetails) string {
	return fmt.Sprintf(`
	Analyze the following product details and reviews and provide a consolidated opinion.
		Ignore any spam or irrelevant text.
		
		Product Details
		- %s

		Reviews:
		- %s
	`, reviews.Specs, reviews.Review)
}

// listAllModels fetches and returns a slice of all model metadata
func listAllModels(ctx context.Context, client *genai.Client) ([]string, error) {

	var modelNames []string
	models, err := client.Models.List(ctx, &genai.ListModelsConfig{})
	if err != nil {
		slog.Error(err.Error())
		return []string{}, err
	}

	for _, m := range models.Items {
		modelNames = append(modelNames, m.Name)
	}

	return modelNames, nil
}

func validateAndReturnAIModelName(name string, models []string) string {
	var isValidModelName bool
	for _, m := range models {
		// API returns names like "models/gemini-pro"
		if m == strings.Split(m, "/")[1] {
			isValidModelName = true
		}
	}

	if !isValidModelName {
		return "gemini-2.5-flash"
	}
	return name

}

func getReviewFromGemini(ctx context.Context, params model.ApiReviewRequestBody) (model.ProductOpinion, error) {
	var opinion model.ProductOpinion

	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  params.ApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("Error while initializing the ai client")
		return opinion, err
	}
	productReviews, err := scraper.ScrapeDataFromURL(params.ProductUrl)
	if err != nil {
		slog.Error("Error while extracting product review")
		return opinion, err
	}

	aiPrompt := genai.Text(generatePrompt(productReviews))

	availableModels, err := listAllModels(ctx, geminiClient)
	if err != nil {
		return opinion, err
	}

	aiModel := validateAndReturnAIModelName(params.AIModel, availableModels)

	result, err := geminiClient.Models.GenerateContent(ctx, aiModel, aiPrompt, getGeminiContentConfig())

	if err != nil {
		slog.Error("Error while fetching the response")
		return opinion, err
	}

	var jsonRaw strings.Builder
	for _, part := range result.Candidates[0].Content.Parts {
		if txt := *part; txt.Text != "" {
			jsonRaw.WriteString(string(txt.Text))
		}
	}

	// Unmarshal into our struct
	if err := json.Unmarshal([]byte(jsonRaw.String()), &opinion); err != nil {
		log.Fatalf("Failed to parse JSON response: %v\nRaw: %s", err, jsonRaw.String())
	}

	return opinion, nil
}
