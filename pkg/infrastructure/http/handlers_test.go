package http_test

import (
	"bytes"
	// "context" // No longer needed in mock signatures directly
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	// Import actual packages to be tested and for domain types
	app_http "github.com/jizumer/expedition-value/pkg/infrastructure/http"
	"github.com/jizumer/expedition-value/pkg/application"
	"github.com/jizumer/expedition-value/pkg/domain/company"
	"github.com/jizumer/expedition-value/pkg/domain/portfolio"

	"github.com/google/uuid"
)

// --- Mock CompanyRepository (for TestCompanyService) ---
type mockCompanyRepository struct {
	FindByTickerFunc       func(ticker string) (*company.Company, error)
	SearchByScoreRangeFunc func(minScore, maxScore float64) ([]*company.Company, error)
	SaveFunc               func(c *company.Company) error
	DeleteFunc             func(ticker string) error
}

func (m *mockCompanyRepository) FindByTicker(ticker string) (*company.Company, error) {
	if m.FindByTickerFunc != nil { return m.FindByTickerFunc(ticker) }
	return nil, errors.New("mockCompanyRepository FindByTicker not implemented")
}
func (m *mockCompanyRepository) SearchByScoreRange(minScore, maxScore float64) ([]*company.Company, error) {
	if m.SearchByScoreRangeFunc != nil { return m.SearchByScoreRangeFunc(minScore, maxScore) }
	return nil, errors.New("mockCompanyRepository SearchByScoreRange not implemented")
}
func (m *mockCompanyRepository) Save(c *company.Company) error {
	if m.SaveFunc != nil { return m.SaveFunc(c) }
	return errors.New("mockCompanyRepository Save not implemented")
}
func (m *mockCompanyRepository) Delete(ticker string) error {
	if m.DeleteFunc != nil { return m.DeleteFunc(ticker) }
	return errors.New("mockCompanyRepository Delete not implemented")
}

// --- TestCompanyService (mock for CompanyHandler, embeds real service) ---
type TestCompanyService struct {
	*application.CompanyService
	mockGetCompanyByTicker func(ticker string) (*company.Company, error)
	mockCreateCompany      func(ticker string, metrics company.FinancialMetrics, sector company.Sector) (*company.Company, error)
    // Add other application.CompanyService methods if they need to be mocked for other tests
    mockSearchCompaniesByScore func(minScore, maxScore float64) ([]*company.Company, error)
    mockUpdateCompanyMetrics   func(ticker string, newMetrics company.FinancialMetrics) error
    mockRefreshCompany         func(ticker string) error
}

func NewTestCompanyService() *TestCompanyService {
	repoMock := &mockCompanyRepository{}
	concreteService := application.NewCompanyService(repoMock)
	return &TestCompanyService{CompanyService: concreteService}
}

func (m *TestCompanyService) GetCompanyByTicker(ticker string) (*company.Company, error) {
	if m.mockGetCompanyByTicker != nil { return m.mockGetCompanyByTicker(ticker) }
	return nil, errors.New("TestCompanyService: GetCompanyByTicker behavior not set")
}
func (m *TestCompanyService) CreateCompany(ticker string, metrics company.FinancialMetrics, sector company.Sector) (*company.Company, error) {
	if m.mockCreateCompany != nil { return m.mockCreateCompany(ticker, metrics, sector) }
	return nil, errors.New("TestCompanyService: CreateCompany behavior not set")
}
// Implement other application.CompanyService methods to use mocks or default behavior
func (m *TestCompanyService) SearchCompaniesByScore(minScore, maxScore float64) ([]*company.Company, error) {
    if m.mockSearchCompaniesByScore != nil { return m.mockSearchCompaniesByScore(minScore, maxScore) }
    return nil, errors.New("TestCompanyService: SearchCompaniesByScore behavior not set")
}
func (m *TestCompanyService) UpdateCompanyMetrics(ticker string, newMetrics company.FinancialMetrics) error {
    if m.mockUpdateCompanyMetrics != nil { return m.mockUpdateCompanyMetrics(ticker, newMetrics) }
    return errors.New("TestCompanyService: UpdateCompanyMetrics behavior not set")
}
func (m *TestCompanyService) RefreshCompany(ticker string) error {
    if m.mockRefreshCompany != nil { return m.mockRefreshCompany(ticker) }
    return errors.New("TestCompanyService: RefreshCompany behavior not set")
}


// --- Mock PortfolioRepository (for TestPortfolioService) ---
type mockPortfolioRepository struct {
	FindByIDFunc func(id string) (*portfolio.Portfolio, error)
	FindAllFunc  func() ([]*portfolio.Portfolio, error)
	SaveFunc     func(p *portfolio.Portfolio) error
	DeleteFunc   func(id string) error
	SearchByRiskProfileFunc func(riskProfile portfolio.RiskProfile) ([]*portfolio.Portfolio, error)
}
func (m *mockPortfolioRepository) FindByID(id string) (*portfolio.Portfolio, error) { if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }; return nil, errors.New("mockPortfolioRepository FindByID not implemented") }
func (m *mockPortfolioRepository) FindAll() ([]*portfolio.Portfolio, error) { if m.FindAllFunc != nil { return m.FindAllFunc() }; return nil, errors.New("mockPortfolioRepository FindAll not implemented") }
func (m *mockPortfolioRepository) Save(p *portfolio.Portfolio) error { if m.SaveFunc != nil { return m.SaveFunc(p) }; return errors.New("mockPortfolioRepository Save not implemented") }
func (m *mockPortfolioRepository) Delete(id string) error { if m.DeleteFunc != nil { return m.DeleteFunc(id) }; return errors.New("mockPortfolioRepository Delete not implemented") }
func (m *mockPortfolioRepository) SearchByRiskProfile(riskProfile portfolio.RiskProfile) ([]*portfolio.Portfolio, error) { if m.SearchByRiskProfileFunc != nil { return m.SearchByRiskProfileFunc(riskProfile) }; return nil, errors.New("mockPortfolioRepository SearchByRiskProfile not implemented")}


// --- TestPortfolioService (mock for PortfolioHandler, embeds real service) ---
type TestPortfolioService struct {
	*application.PortfolioService
	mockCreatePortfolio     func(cashBalance portfolio.Money, riskProfile portfolio.RiskProfile) (*portfolio.Portfolio, error)
	mockGetPortfolioDetails func(portfolioID string) (*portfolio.Portfolio, error)
    // Add other application.PortfolioService methods if they need to be mocked
    mockAddPosition          func(portfolioID string, companyTicker string, shares int, purchasePrice portfolio.Money) error
    mockAdjustPosition       func(portfolioID string, companyTicker string, newShares int) error
    mockRecommendRebalance   func(portfolioID string) (*application.RebalanceRecommendation, error)
    mockExecuteRebalance     func(portfolioID string, recommendation application.RebalanceRecommendation) error
}

func NewTestPortfolioService() *TestPortfolioService {
	portfolioRepoMock := &mockPortfolioRepository{}
	companyRepoMock := &mockCompanyRepository{}
	concreteService := application.NewPortfolioService(portfolioRepoMock, companyRepoMock)
	return &TestPortfolioService{PortfolioService: concreteService}
}
func (m *TestPortfolioService) CreatePortfolio(cashBalance portfolio.Money, riskProfile portfolio.RiskProfile) (*portfolio.Portfolio, error) {
	if m.mockCreatePortfolio != nil { return m.mockCreatePortfolio(cashBalance, riskProfile) }
	return nil, errors.New("TestPortfolioService: CreatePortfolio behavior not set")
}
func (m *TestPortfolioService) GetPortfolioDetails(portfolioID string) (*portfolio.Portfolio, error) {
	if m.mockGetPortfolioDetails != nil { return m.mockGetPortfolioDetails(portfolioID) }
	return nil, errors.New("TestPortfolioService: GetPortfolioDetails behavior not set")
}
// Implement other application.PortfolioService methods
func (m *TestPortfolioService) AddPosition(portfolioID string, companyTicker string, shares int, purchasePrice portfolio.Money) error {
    if m.mockAddPosition != nil { return m.mockAddPosition(portfolioID, companyTicker, shares, purchasePrice) }
    return errors.New("TestPortfolioService: AddPosition behavior not set")
}
func (m *TestPortfolioService) AdjustPosition(portfolioID string, companyTicker string, newShares int) error {
    if m.mockAdjustPosition != nil { return m.mockAdjustPosition(portfolioID, companyTicker, newShares) }
    return errors.New("TestPortfolioService: AdjustPosition behavior not set")
}
func (m *TestPortfolioService) RecommendRebalance(portfolioID string) (*application.RebalanceRecommendation, error) {
    if m.mockRecommendRebalance != nil { return m.mockRecommendRebalance(portfolioID) }
    return nil, errors.New("TestPortfolioService: RecommendRebalance behavior not set")
}
func (m *TestPortfolioService) ExecuteRebalance(portfolioID string, recommendation application.RebalanceRecommendation) error {
    if m.mockExecuteRebalance != nil { return m.mockExecuteRebalance(portfolioID, recommendation) }
    return errors.New("TestPortfolioService: ExecuteRebalance behavior not set")
}

// --- Test Helper ---
func executeRequest(req *http.Request, handler http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// --- CompanyHandler Tests ---
func TestCompanyHandler_GetCompanyByTicker(t *testing.T) {
	serviceMock := NewTestCompanyService()
	handler := app_http.NewCompanyHandler(serviceMock)

	t.Run("Success", func(t *testing.T) {
		expectedCompany, _ := company.NewCompany("AAPL", company.FinancialMetrics{PERatio: 15.5}, company.Technology)
		expectedCompany.UpdatedAt = time.Now()

		serviceMock.mockGetCompanyByTicker = func(ticker string) (*company.Company, error) {
			if ticker == "AAPL" { return expectedCompany, nil }
			return nil, errors.New("company not found")
		}

		req, _ := http.NewRequest("GET", "/company?ticker=AAPL", nil)
		rr := executeRequest(req, handler.GetCompanyByTicker)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		var returnedCompany company.Company
		if err := json.NewDecoder(rr.Body).Decode(&returnedCompany); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
		if returnedCompany.Ticker != "AAPL" {
			t.Errorf("handler returned unexpected body: got ticker %v want %v", returnedCompany.Ticker, "AAPL")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		serviceMock.mockGetCompanyByTicker = func(ticker string) (*company.Company, error) {
			return nil, errors.New("company not found an error")
		}
		req, _ := http.NewRequest("GET", "/company?ticker=UNKNOWN", nil)
		rr := executeRequest(req, handler.GetCompanyByTicker)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
			t.Fatalf("could not decode error response: %v", err)
		}
		if errResp.Error != "company not found" {
			t.Errorf("handler returned unexpected error message: got %q want %q", errResp.Error, "company not found")
		}
	})

	t.Run("EmptyTicker", func(t *testing.T) {
		serviceMock.mockGetCompanyByTicker = nil
		req, _ := http.NewRequest("GET", "/company?ticker=", nil)
		rr := executeRequest(req, handler.GetCompanyByTicker)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		var errResp app_http.ErrorResponse
		err := json.NewDecoder(rr.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("could not decode error response: %v", err)
		}
		if errResp.Error != "ticker query parameter is required" {
			t.Errorf("handler returned unexpected error message: got %q want %q", errResp.Error, "ticker query parameter is required")
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		serviceMock.mockGetCompanyByTicker = func(ticker string) (*company.Company, error) {
			return nil, errors.New("some internal service error")
		}
		req, _ := http.NewRequest("GET", "/company?ticker=ANY", nil)
		rr := executeRequest(req, handler.GetCompanyByTicker)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
			t.Fatalf("could not decode error response: %v", err)
		}
		if errResp.Error != "internal server error" {
			t.Errorf("handler returned unexpected error message: got %q want %q", errResp.Error, "internal server error")
		}
	})
}

func TestCompanyHandler_CreateCompany(t *testing.T) {
	serviceMock := NewTestCompanyService()
	handler := app_http.NewCompanyHandler(serviceMock)

	t.Run("Success", func(t *testing.T) {
		defaultMetrics := company.FinancialMetrics{}
		defaultSector := company.UndefinedSector
		createdComp, _ := company.NewCompany("NEWCO", defaultMetrics, defaultSector)
		serviceMock.mockCreateCompany = func(ticker string, metrics company.FinancialMetrics, sector company.Sector) (*company.Company, error) {
			if ticker == "NEWCO" { return createdComp, nil }
			return nil, errors.New("unexpected ticker for create")
		}
		payload := app_http.CreateCompanyRequest{Ticker: "NEWCO", Name: "New Company Inc."}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/company/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreateCompany)
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
		var returnedCompany company.Company
		if err := json.NewDecoder(rr.Body).Decode(&returnedCompany); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
		if returnedCompany.Ticker != "NEWCO" {
			t.Errorf("handler returned unexpected ticker: got %v want NEWCO", returnedCompany.Ticker)
		}
	})

	t.Run("BadRequest_InvalidPayload", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/company/create", strings.NewReader("{malformed_json"))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreateCompany)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "invalid request payload" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "invalid request payload")
		}
	})

	t.Run("BadRequest_EmptyTickerInPayload", func(t *testing.T) {
		payload := app_http.CreateCompanyRequest{Ticker: "", Name: "No Ticker Inc."}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/company/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreateCompany)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "ticker is required" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "ticker is required")
		}
	})

	t.Run("Conflict_AlreadyExists", func(t *testing.T) {
		serviceMock.mockCreateCompany = func(ticker string, metrics company.FinancialMetrics, sector company.Sector) (*company.Company, error) {
			return nil, errors.New("company already exists")
		}
		payload := app_http.CreateCompanyRequest{Ticker: "EXIST", Name: "Existing Company Inc."}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/company/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreateCompany)
		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "company already exists" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "company already exists")
		}
	})

	t.Run("ServiceError_Generic", func(t *testing.T) {
		serviceMock.mockCreateCompany = func(ticker string, metrics company.FinancialMetrics, sector company.Sector) (*company.Company, error) {
			return nil, errors.New("some other internal service error")
		}
		payload := app_http.CreateCompanyRequest{Ticker: "ANY", Name: "Any Company Inc."}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/company/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreateCompany)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "internal server error" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "internal server error")
		}
	})
}

// --- PortfolioHandler Tests ---
func TestPortfolioHandler_CreatePortfolio(t *testing.T) {
	serviceMock := NewTestPortfolioService()
	handler := app_http.NewPortfolioHandler(serviceMock)

	t.Run("Success", func(t *testing.T) {
		reqCash, _ := portfolio.NewMoney(100000, "USD")
		reqRisk := portfolio.Moderate
		createdPortfolio, _ := portfolio.NewPortfolio(uuid.NewString(), reqRisk, *reqCash)
		serviceMock.mockCreatePortfolio = func(cashBalance portfolio.Money, riskProfile portfolio.RiskProfile) (*portfolio.Portfolio, error) {
			if cashBalance.Amount == reqCash.Amount && riskProfile == reqRisk { return createdPortfolio, nil }
			return nil, errors.New("mock CreatePortfolio called with unexpected params")
		}
		payload := app_http.CreatePortfolioRequest{CashBalance: *reqCash, RiskProfile: reqRisk}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/portfolio/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreatePortfolio)
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
		var respPortfolio portfolio.Portfolio
		if err := json.NewDecoder(rr.Body).Decode(&respPortfolio); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
		if respPortfolio.ID != createdPortfolio.ID {
			t.Errorf("handler returned unexpected portfolio ID: got %v want %v", respPortfolio.ID, createdPortfolio.ID)
		}
	})

	t.Run("BadRequest_InvalidPayload", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/portfolio/create", strings.NewReader("{malformed}"))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreatePortfolio)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "invalid request payload" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "invalid request payload")
		}
	})

	t.Run("BadRequest_ValidationErrorFromService", func(t *testing.T) {
		serviceErrorMsg := "initial cash balance cannot be negative"
		serviceMock.mockCreatePortfolio = func(cb portfolio.Money, rp portfolio.RiskProfile) (*portfolio.Portfolio, error) {
			return nil, errors.New(serviceErrorMsg)
		}
		invalidCash, _ := portfolio.NewMoney(-100, "USD")
		payload := app_http.CreatePortfolioRequest{CashBalance: *invalidCash, RiskProfile: portfolio.Conservative}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/portfolio/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreatePortfolio)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != serviceErrorMsg {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, serviceErrorMsg)
		}
	})

	t.Run("ServiceError_Generic", func(t *testing.T) {
		serviceMock.mockCreatePortfolio = func(cb portfolio.Money, rp portfolio.RiskProfile) (*portfolio.Portfolio, error) {
			return nil, errors.New("some internal repository error")
		}
		reqCash, _ := portfolio.NewMoney(1000, "USD")
		payload := app_http.CreatePortfolioRequest{CashBalance: *reqCash, RiskProfile: portfolio.Aggressive}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/portfolio/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(req, handler.CreatePortfolio)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "internal server error" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "internal server error")
		}
	})
}

func TestPortfolioHandler_GetPortfolioDetails(t *testing.T) {
	serviceMock := NewTestPortfolioService()
	handler := app_http.NewPortfolioHandler(serviceMock)

	t.Run("Success", func(t *testing.T) {
		portfolioID := uuid.NewString()
		cash, _ := portfolio.NewMoney(1000, "USD")
		expectedPortfolio, _ := portfolio.NewPortfolio(portfolioID, portfolio.Conservative, *cash)
		serviceMock.mockGetPortfolioDetails = func(id string) (*portfolio.Portfolio, error) {
			if id == portfolioID { return expectedPortfolio, nil }
			return nil, errors.New("portfolio not found")
		}
		req, _ := http.NewRequest("GET", "/portfolio?id="+portfolioID, nil)
		rr := executeRequest(req, handler.GetPortfolioDetails)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		var respPortfolio portfolio.Portfolio
		if err := json.NewDecoder(rr.Body).Decode(&respPortfolio); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
		if respPortfolio.ID != portfolioID {
			t.Errorf("handler returned unexpected portfolio ID: got %v want %v", respPortfolio.ID, portfolioID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		serviceMock.mockGetPortfolioDetails = func(id string) (*portfolio.Portfolio, error) {
			return nil, errors.New("some portfolio not found error from service")
		}
		req, _ := http.NewRequest("GET", "/portfolio?id=UNKNOWN_ID", nil)
		rr := executeRequest(req, handler.GetPortfolioDetails)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "portfolio not found" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "portfolio not found")
		}
	})

	t.Run("EmptyID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/portfolio?id=", nil)
		rr := executeRequest(req, handler.GetPortfolioDetails)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "portfolio id query parameter is required" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "portfolio id query parameter is required")
		}
	})

	t.Run("ServiceError_Generic", func(t *testing.T) {
		serviceMock.mockGetPortfolioDetails = func(id string) (*portfolio.Portfolio, error) {
			return nil, errors.New("some other service layer error")
		}
		req, _ := http.NewRequest("GET", "/portfolio?id=ANY_ID", nil)
		rr := executeRequest(req, handler.GetPortfolioDetails)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		var errResp app_http.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil { t.Fatalf("could not decode error response: %v", err) }
		if errResp.Error != "internal server error" {
			t.Errorf("unexpected error message: got %q want %q", errResp.Error, "internal server error")
		}
	})
}
// Removed conceptual var _ declarations and placeholder service methods that used old mock types
// Removed "Okay"
