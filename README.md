## Overview
This service implements value investment principles to automate portfolio management through dynamic scoring and decision automation. The system monitors financial indicators from multiple integrated data sources to generate real-time investment recommendations.

## Core Functionality
- **Automated Scoring Engine**: Continuously evaluates assets using Graham-style valuation formulas combining:
  - Fundamental analysis (P/E ratio, book value, earnings growth)
  - Market position indicators
  - Sector-specific metrics
- **Portfolio Optimization**: Maintains:
  - Current invested positions with cost basis tracking
  - Watchlist of potential investments
  - Historical performance analytics
- **Decision Automation**: Generates four recommendation types:
  1. New position entry signals
  2. Position size increase alerts
  3. Partial reduction warnings
  4. Full liquidation triggers

## Key Features
- **Multi-source Data Integration**:  
  Unified pipeline for financial APIs (Alpha Vantage, Yahoo Finance) and custom data sources
- **Configurable Refresh Intervals**:  
  Adjustable evaluation cycles with default hourly intervals
- **Cloud-native Architecture**:  
  Designed for GCP Free Tier deployment with Firestore persistence
- **Extensible Scoring Model**:  
  Modular algorithm design supporting custom valuation formulas

## Domain-Driven Design Components:
	•	Portfolio Aggregate: Manages positions and rebalancing logic
	•	Company Aggregate: Handles financial data and score calculations
	•	Recommendation Context: Implements decision workflows
## Architecture:
	•	Clean Architecture with three-layer separation
	•	Event-driven design using Pub/Sub patterns
	•	REST API surface for portfolio management

## Non-Functional Requirements
- **Fault Tolerance**: Automatic retries for failed API calls
- **Audit Logging**: Immutable record of all portfolio changes
- **Configuration First**: All parameters exposed via environment variables
- **Testability**: Built-in support for golden master testing

**Architecture**: Modular monolith → Cloud-native evolution  
**Tech Stack**: Go 1.22, Fiber, Viper, PostgreSQL, GCP

## Agentic Generation Notes
This system requires:
- Pure domain models in Go 1.22
- Cloud-agnostic interfaces for persistence
- Table-driven test implementations
- Fiber-compatible handler signatures
- Viper configuration mapping

**Code Generation Rules**:
1. Follow Clean Architecture layers
2. Implement Value Objects for:
   - Financial metrics
   - Investment scores
   - Money operations on the portfolio
3. Use Fiber for REST API endpoints
4. Apply Viper for cloud-agnostic config

**Testing Requirements**:
- 90%+ coverage on domain layer
- Table-driven tests for score calculations
- Golden files for complex financial models

---

## Credits & Thanks

- Inspired by [CodelyTV DDD Blueprints](https://codely.com/en/blog/how-to-implement-ddd-code-using-ai)
- Built with [Go](https://golang.org/), [Fiber](https://gofiber.io/), [Viper](https://github.com/spf13/viper), [PostgreSQL](https://www.postgresql.org/), and [Google Cloud Platform](https://cloud.google.com/)
- Project structure and DDD approach based on community best practices
