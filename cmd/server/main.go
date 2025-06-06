// @title           Value Investment Analysis API
// @version         1.0
// @description     This is an API for the Value Investment Analysis and Portfolio Management MVP.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
package main

import (
	"log"
	"net/http"

	// Project packages
	"github.com/jizumer/expedition-value/pkg/application"
	infHttp "github.com/jizumer/expedition-value/pkg/infrastructure/http"
	"github.com/jizumer/expedition-value/pkg/infrastructure/persistence/memory"

	// Swagger imports
	_ "github.com/jizumer/expedition-value/cmd/server/docs" // Generated Swagger docs
	httpSwagger "github.com/swaggo/http-swagger"            // http-swagger
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

	// Swagger UI handler
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")

	log.Println("HTTP routes configured.")

	// 3. Start Server
	port := ":8080"
	log.Printf("Server listening on port %s\n", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
