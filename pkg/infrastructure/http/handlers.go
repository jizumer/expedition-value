package http

import (
	"encoding/json"
	// "errors" // Unused, removed
	"net/http"
	"strings"

	// "github.com/gorilla/mux" // Example router, not strictly needed for placeholders

	// "context" // No longer needed as service interfaces don't use context yet
	// "github.com/jizumer/expedition-value/pkg/application" // No longer needed as handlers use local interfaces
	"github.com/jizumer/expedition-value/pkg/domain/company"
	"github.com/jizumer/expedition-value/pkg/domain/portfolio"
)

// --- Service Interfaces (for Dependency Injection) ---

// CompanyServiceProvider defines the interface for company service operations needed by handlers.
type CompanyServiceProvider interface {
	GetCompanyByTicker(ticker string) (*company.Company, error)
	CreateCompany(ticker string, metrics company.FinancialMetrics, sector company.Sector) (*company.Company, error)
	// Add other methods from application.CompanyService that handlers might use
}

// PortfolioServiceProvider defines the interface for portfolio service operations needed by handlers.
type PortfolioServiceProvider interface {
	CreatePortfolio(cashBalance portfolio.Money, riskProfile portfolio.RiskProfile) (*portfolio.Portfolio, error)
	GetPortfolioDetails(portfolioID string) (*portfolio.Portfolio, error)
	// Add other methods from application.PortfolioService that handlers might use
}


// ErrorResponse represents a generic error response.
type ErrorResponse struct {
	Error string `json:"error" example:"Detailed error message"`
}

// CompanyHandler holds dependencies for company-related HTTP handlers.
type CompanyHandler struct {
	service CompanyServiceProvider // Use the interface
}

// NewCompanyHandler creates a new CompanyHandler.
func NewCompanyHandler(cs CompanyServiceProvider) *CompanyHandler { // Accept the interface
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
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	comp, err := h.service.GetCompanyByTicker(ticker) // Removed r.Context()
	if err != nil {
		// Assuming a specific error type application.ErrCompanyNotFound or similar might be defined.
		// For now, checking string content is a placeholder.
		// A more robust solution would be to use errors.Is() with a specific error variable.
		if strings.Contains(strings.ToLower(err.Error()), "not found") { // Basic check
			respondWithError(w, http.StatusNotFound, "company not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, comp)
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
// @Failure      409  {object}  ErrorResponse "Company already exists"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /company/create [post]
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate basic input - e.g., ticker is required by the request DTO itself for this handler
	if req.Ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker is required")
		return
	}
	// Name could also be validated here if desired, e.g., if req.Name == "" ...

	// As per subtask, Name from req is not passed to current service signature.
	// Using default FinancialMetrics and UndefinedSector.
	metrics := company.FinancialMetrics{}
	sector := company.UndefinedSector // Assuming company.UndefinedSector is defined.

	comp, err := h.service.CreateCompany(req.Ticker, metrics, sector) // Removed r.Context()
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "already exists") || strings.Contains(errStr, "conflict") {
			respondWithError(w, http.StatusConflict, "company already exists")
		} else if strings.Contains(errStr, "validation failed") || strings.Contains(errStr, "invalid ticker") { // Example validation checks
			respondWithError(w, http.StatusBadRequest, err.Error()) // Or a more generic "invalid data"
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, comp)
}

// PortfolioHandler holds dependencies for portfolio-related HTTP handlers.
type PortfolioHandler struct {
	service PortfolioServiceProvider // Use the interface
}

// NewPortfolioHandler creates a new PortfolioHandler.
func NewPortfolioHandler(ps PortfolioServiceProvider) *PortfolioHandler { // Accept the interface
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
	var req CreatePortfolioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	// Basic validation can be added here if necessary, e.g.
	// if req.RiskProfile == "" { // Assuming RiskProfile could be an empty string for invalid
	//    respondWithError(w, http.StatusBadRequest, "riskProfile is required")
	//    return
	// }
	// if req.CashBalance.Amount < 0 { // Assuming Amount is accessible and comparable
	//    respondWithError(w, http.StatusBadRequest, "cashBalance amount cannot be negative")
	//    return
	// }


	p, err := ph.service.CreatePortfolio(req.CashBalance, req.RiskProfile) // Removed r.Context()
	if err != nil {
		errStr := strings.ToLower(err.Error())
		// Keywords for domain validation errors
		if strings.Contains(errStr, "validation") ||
			strings.Contains(errStr, "invalid") ||
			strings.Contains(errStr, "negative") || // Made this more general to catch "negative cash balance"
			strings.Contains(errStr, "unknown risk profile") {
			respondWithError(w, http.StatusBadRequest, err.Error()) // Send back the specific domain error
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
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
	portfolioID := r.URL.Query().Get("id")
	if portfolioID == "" {
		respondWithError(w, http.StatusBadRequest, "portfolio id query parameter is required")
		return
	}

	p, err := ph.service.GetPortfolioDetails(portfolioID) // Removed r.Context()
	if err != nil {
		// A more robust way would be to use errors.Is(err, portfolio.ErrPortfolioNotFound)
		// if portfolio.ErrPortfolioNotFound is a well-defined error.
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			respondWithError(w, http.StatusNotFound, "portfolio not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
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
