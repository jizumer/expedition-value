package http

import (
	"encoding/json"
	"net/http"

	// "github.com/gorilla/mux" // Example router, not strictly needed for placeholders

	"github.com/jizumer/expedition-value/pkg/application"
	"github.com/jizumer/expedition-value/pkg/domain/portfolio" // For request/response types and annotations
)

// ErrorResponse represents a generic error response.
type ErrorResponse struct {
	Error string `json:"error" example:"Detailed error message"`
}

// CompanyHandler holds dependencies for company-related HTTP handlers, primarily the CompanyService.
type CompanyHandler struct {
	service *application.CompanyService
}

// NewCompanyHandler creates a new CompanyHandler.
func NewCompanyHandler(cs *application.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: cs}
}

// CreateCompanyRequest defines the structure for creating a new company.
type CreateCompanyRequest struct {
	Ticker string `json:"ticker" example:"AAPL"`
	Name   string `json:"name" example:"Apple Inc."`
	// For MVP, sector and initial metrics might be set via other means or have defaults.
	// If they were to be included:
	// Sector string `json:"sector" example:"Technology"`
	// PERatio float64 `json:"peRatio" example:"15.5"`
}

// GetCompanyByTicker godoc
// @Summary      Get company by ticker
// @Description  Get company details by its stock ticker
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        ticker query string true "Company Ticker"
// @Success      200  {object}  company.Company "Successfully retrieved company"
// @Failure      400  {object}  ErrorResponse "Invalid request (e.g., missing ticker)"
// @Failure      404  {object}  ErrorResponse "Company not found"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /company [get]
func (h *CompanyHandler) GetCompanyByTicker(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract ticker from r.URL.Query().Get("ticker")
	// TODO: Call h.service.GetCompanyByTicker(ticker)
	// TODO: Handle errors from service (e.g., not found, validation)
	// TODO: Marshal response to JSON and write to w
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "GetCompanyByTicker Not Implemented"})
}

// CreateCompany godoc
// @Summary      Create a new company
// @Description  Adds a new company to the system.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        company body CreateCompanyRequest true "Company data to create"
// @Success      201  {object}  company.Company "Successfully created company"
// @Failure      400  {object}  ErrorResponse "Invalid company data provided"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /company/create [post]
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode request body into a CreateCompanyRequest DTO
	// TODO: Call h.service.CreateCompany(dto.Ticker, company.FinancialMetrics{...} /* from DTO or default */, company.ParseSector(dto.Sector))
	// TODO: Handle errors
	// TODO: Marshal created company to JSON and write response with 201 status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "CreateCompany Not Implemented"})
}

// PortfolioHandler holds dependencies for portfolio-related HTTP handlers.
type PortfolioHandler struct {
	service *application.PortfolioService
}

// NewPortfolioHandler creates a new PortfolioHandler.
func NewPortfolioHandler(ps *application.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{service: ps}
}

// CreatePortfolioRequest DTO for creating a portfolio
type CreatePortfolioRequest struct {
	CashBalance portfolio.Money       `json:"cashBalance"` // e.g. {"amount": 100000, "currency": "USD"}
	RiskProfile portfolio.RiskProfile `json:"riskProfile" example:"Moderate" enums:"Conservative,Moderate,Aggressive,UndefinedProfile"`
}

// CreatePortfolio godoc
// @Summary      Create a new portfolio
// @Description  Creates a new investment portfolio.
// @Tags         portfolios
// @Accept       json
// @Produce      json
// @Param        portfolio body CreatePortfolioRequest true "Portfolio data to create"
// @Success      201  {object}  portfolio.Portfolio "Successfully created portfolio"
// @Failure      400  {object}  ErrorResponse "Invalid portfolio data provided"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /portfolio/create [post]
func (ph *PortfolioHandler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode request body into CreatePortfolioRequest
	// TODO: Call ph.service.CreatePortfolio(dto.CashBalance, dto.RiskProfile)
	// TODO: Handle errors
	// TODO: Marshal created portfolio to JSON and write response with 201 status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "CreatePortfolio Not Implemented"})
}

// GetPortfolioDetails godoc
// @Summary      Get portfolio details
// @Description  Get details of a specific portfolio by its ID.
// @Tags         portfolios
// @Accept       json
// @Produce      json
// @Param        id query string true "Portfolio ID"
// @Success      200  {object}  portfolio.Portfolio "Successfully retrieved portfolio"
// @Failure      400  {object}  ErrorResponse "Invalid request (e.g., missing ID)"
// @Failure      404  {object}  ErrorResponse "Portfolio not found"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /portfolio [get]
func (ph *PortfolioHandler) GetPortfolioDetails(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract portfolioID from r.URL.Query().Get("id")
	// TODO: Call ph.service.GetPortfolioDetails(portfolioID)
	// TODO: Handle errors
	// TODO: Marshal portfolio to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "GetPortfolioDetails Not Implemented"})
}

// --- Utility functions for handlers (optional, can be in a separate file) ---

// respondWithError is a helper function to send a JSON error response.
// For Swaggo, if ErrorResponse is used in @Failure, this function should marshal ErrorResponse.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

// respondWithJSON is a helper function to send a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// HealthCheckHandler godoc
// @Summary      Show the status of server.
// @Description  Get the status of server.
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string "Successfully retrieved health status"
// @Router       /health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
