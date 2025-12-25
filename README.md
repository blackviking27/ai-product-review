# Product Review Analyzer API

A Go-based backend service that intelligently analyzes product reviews from e-commerce platforms (Amazon, Flipkart) to determine if a product is worth buying.

It automates the process of scraping user reviews, feeding them into an LLM (AI), and generating a detailed purchasing verdict including Pros, Cons, and a "Buy/Pass" recommendation.

## 🚀 Key Features

- **Multi-Platform Scraping:** capabilities for Amazon and Flipkart (extensible design).
- **AI-Powered Analysis:** Uses LLMs (e.g., OpenAI) to process sentiment and detect common product faults.
- **Clean Architecture:** Modular code structure separating scraping mechanics, business logic, and HTTP handling.

---

## 📂 Project Structure

This project follows a pragmatic, domain-focused folder structure to ensure separation of concerns. Here is how the code is organized:

```text
product-review-analyzer/
│
├── config/
│   └── config.yaml          # External configuration (CSS selectors, API timeouts, Base URLs).
│
├── internal/
│   │
│   ├── model/               # The "Shared Language". Contains structs used across the app.
│   │                        # (e.g., Review, Product, AnalysisResult). No logic, just data definitions.
│   │
│   ├── scraper/             # The "Collector". Responsible for fetching raw HTML and parsing data.
│   │   ├── amazon.go        # Amazon-specific DOM parsing logic.
│   │   ├── flipkart.go      # Flipkart-specific DOM parsing logic.
│   │   └── scraper.go       # The Interface & Factory. Decides which file to use based on the URL.
│   │
│   ├── analyzer/            # The "Brain". Responsible for decision making.
│   │   └── openai.go        # Sends formatted reviews to the AI API and parses the verdict.
│   │
│   └── api/                 # The "Doorway". Responsible for HTTP communication.
│       ├── routes.go        # URL routing definitions.
│       └── handlers.go      # Controllers that accept requests and orchestrate the Scraper -> Analyzer flow.
│
├── go.mod                   # Go Module definition.
|── main.go                  # Server logic
```
