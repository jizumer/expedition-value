package main

import (
	"log"
	"net/http"

	// Project packages
	"github.com/jizumer/expedition-value/pkg/application"
	infHttp "github.com/jizumer/expedition-value/pkg/infrastructure/http"
	"github.com/jizumer/expedition-value/pkg/infrastructure/persistence/memory"

	// No external router for MVP, using net/http.ServeMux
)

func main() {
	log.Println("Starting Value Investment Analysis MVP server...")

	// 1. Initialization
	log.Println("Initializing repositories and services...")

	// Instantiate Repositories
	companyRepo := memory.NewInMemoryCompanyRepository()
	// Portfolio repo needs company repo for some operations (e.g., SearchBySector, if implemented fully)
	portfolioRepo := memory.NewInMemoryPortfolioRepository(companyRepo)

	// Instantiate Application Services
	companyService := application.NewCompanyService(companyRepo)
	portfolioService := application.NewPortfolioService(portfolioRepo, companyRepo)

	// Instantiate HTTP Handlers
	companyHandler := infHttp.NewCompanyHandler(companyService)
	portfolioHandler := infHttp.NewPortfolioHandler(portfolioService)

	log.Println("Initialization complete.")

	// 2. HTTP Routing
	log.Println("Setting up HTTP routes...")
	mux := http.NewServeMux()

	// Root handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" { // Basic check to prevent matching all paths
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Value Investment Analysis MVP API Root"}`))
	})

	// Health check
	mux.HandleFunc("/health", infHttp.HealthCheckHandler)

	// Company routes
	// GetCompanyByTicker expects GET with ?ticker=XYZ
	// The handler infHttp.CompanyHandler.GetCompanyByTicker needs to be implemented
	// to parse r.URL.Query().Get("ticker")
	mux.HandleFunc("/company", companyHandler.GetCompanyByTicker)

	// CreateCompany expects POST
	// The handler infHttp.CompanyHandler.CreateCompany needs to be implemented
	// to check r.Method == http.MethodPost and parse the request body.
	mux.HandleFunc("/company/create", companyHandler.CreateCompany)

	// Portfolio routes
	// GetPortfolioDetails expects GET with ?id=XYZ
	// The handler infHttp.PortfolioHandler.GetPortfolioDetails needs to be implemented
	// to parse r.URL.Query().Get("id")
	mux.HandleFunc("/portfolio", portfolioHandler.GetPortfolioDetails)

	// CreatePortfolio expects POST
	// The handler infHttp.PortfolioHandler.CreatePortfolio needs to be implemented
	// to check r.Method == http.MethodPost and parse the request body.
	mux.HandleFunc("/portfolio/create", portfolioHandler.CreatePortfolio)

	log.Println("HTTP routes configured.")

	// 3. Start Server
	port := ":8080"
	log.Printf("Server listening on port %s\n", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
