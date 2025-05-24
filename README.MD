# Value Investment Analysis and Portfolio Management

## Overview
This service implements value investment principles to automate portfolio management through dynamic scoring and decision automation. The system monitors financial indicators from multiple integrated data sources to generate real-time investment recommendations. This MVP focuses on core analysis and management capabilities.

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
- **Multi-source Data Integration**: Unified pipeline for financial APIs (e.g., Alpha Vantage, Yahoo Finance) and custom data sources.
- **Configurable Refresh Intervals**: Adjustable evaluation cycles with default hourly intervals.
- **Cloud-native Architecture Target**: Designed for GCP deployment with Firestore/PostgreSQL persistence. (MVP currently uses in-memory).
- **Extensible Scoring Model**: Modular algorithm design supporting custom valuation formulas.

## Tech Stack & Architecture

*   **Backend Language:** Go (Version 1.22)
*   **Web Framework:** Fiber (Target)
*   **Configuration Management:** Viper (Target)
*   **Database:** PostgreSQL (Target), In-memory (Current MVP)
*   **Deployment Target:** Google Cloud Platform (GCP)
*   **Architecture Principles:**
    *   **Clean Architecture:** Adhering to a three-layer separation (Domain, Application, Infrastructure).
    *   **Domain-Driven Design (DDD):**
        *   **Modularity:** Organized into packages reflecting DDD layers (`pkg/domain`, `pkg/application`, `pkg/infrastructure`).
        *   **Aggregates:** Core components include `Portfolio` (manages positions, rebalancing) and `Company` (handles financial data, score calculations). A `Recommendation Context` is envisioned for decision workflows.
    *   **Event-Driven Design:** Using Pub/Sub patterns (Target).
*   **API:** RESTful HTTP API (Implemented via Go standard library for MVP, target Fiber).

## Initial Constraints (MVP Focus)

*   **MVP Scope:** Focus on core features for company analysis and basic portfolio management as defined in the DDD blueprints (`docs/domain/`).
*   **DDD Principles:** Development will adhere to Domain Driven Design principles, using the guidance in `docs/guidelines.md` and the aggregate designs in `docs/domain/`.
*   **Deployability:** The application must be runnable locally and be structured for straightforward cloud deployment.
*   **Documentation:** Instructions for running the MVP in both environments will be provided. DDD documents will be updated to reflect any design decisions made during development.

## Non-Functional Requirements (Target)
- **Fault Tolerance**: Automatic retries for failed API calls.
- **Audit Logging**: Immutable record of all portfolio changes.
- **Configuration First**: All parameters exposed via environment variables.
- **Testability**: Built-in support for golden master testing.

## How to Run Locally (Current MVP)

### Prerequisites
*   Go (version 1.22 or higher) installed on your system.
*   A Go workspace set up.

### Steps
1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone <repository_url> # Replace with actual URL
    cd expedition-value 
    ```
    (The project module is `github.com/jizumer/expedition-value`.)

2.  **Run the application (from the project root directory `expedition-value`):**
    ```bash
    go run ./cmd/server/main.go
    ```
    This will start the HTTP server, typically on port `8080`.

3.  **Access the API:**
    *   **Swagger UI (API Documentation):** Once the server is running, access the interactive API documentation via your browser at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html). This UI allows you to view all available API endpoints, their parameters, and test them directly.
    *   Health check: `curl http://localhost:8080/health`
    *   Get company (example, requires data to be created first via API): `curl http://localhost:8080/company?ticker=AAPL`
    *   Create company (example JSON body):
        ```bash
        curl -X POST -H "Content-Type: application/json" -d '{"ticker":"MSFT","name":"Microsoft Corp"}' http://localhost:8080/company/create 
        ``` 
        (Note: The actual JSON structure for creating a company will depend on the `CreateCompanyRequest` DTO defined in the HTTP handlers and the fields required by the `CompanyService.CreateCompany` method. The example above is a generic placeholder.)


## Deployment to Cloud (Conceptual for MVP, Target GCP)

This application is designed to be deployable to cloud platforms that support Go applications or containers. Here's a general approach:

1.  **Containerize the Application (using Docker):**
    *   Create a `Dockerfile` in the root of the project. A basic Go Dockerfile might look like this:

        ```Dockerfile
        # Use an official Go runtime as a parent image
        FROM golang:1.22-alpine AS builder

        # Set the Current Working Directory inside the container
        WORKDIR /app

        # Copy go mod and sum files
        COPY go.mod go.sum ./

        # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
        RUN go mod download

        # Copy the source code into the container
        COPY . .

        # Build the Go app for the server
        RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server/main.go

        # Start a new stage
        FROM gcr.io/distroless/static-debian11 

        WORKDIR /

        # Copy the Pre-built binary file from the previous stage
        COPY --from=builder /app/server /server

        # Expose port 8080
        EXPOSE 8080

        # Command to run the executable
        ENTRYPOINT ["/server"]
        ```
    *   Build the Docker image: `docker build -t expedition-value:latest .`

2.  **Choose a Cloud Platform and Service:**
    *   **Google Cloud:** Cloud Run, Google Kubernetes Engine (GKE)
    *   **AWS:** Elastic Beanstalk, Amazon Elastic Container Service (ECS), Amazon EKS
    *   **Azure:** Azure App Service, Azure Kubernetes Service (AKS)

3.  **Push the Docker Image to a Registry:**
    *   E.g., Google Container Registry (GCR), Docker Hub, AWS Elastic Container Registry (ECR).
    *   `docker tag expedition-value:latest gcr.io/your-gcp-project/expedition-value:latest`
    *   `docker push gcr.io/your-gcp-project/expedition-value:latest`

4.  **Deploy the Image:**
    *   Follow the specific instructions for your chosen cloud service to deploy the container image. This usually involves configuring service settings, environment variables (if any), scaling options, and public access.

**Future Considerations for Cloud Deployment:**
*   **Configuration Management:** Externalize configuration (ports, database connection strings) instead of hardcoding.
*   **Structured Logging:** Implement structured logging for better observability.
*   **Persistent Storage:** Replace in-memory repositories with actual database backends (e.g., PostgreSQL, MySQL) managed by the cloud provider (e.g., Cloud SQL, AWS RDS).
*   **CI/CD Pipeline:** Set up a Continuous Integration/Continuous Deployment pipeline to automate builds and deployments.


## Credits & Thanks

- Inspired by [CodelyTV DDD Blueprints](https://codely.com/en/blog/how-to-implement-ddd-code-using-ai)
- Built with [Go](https://golang.org/), and targeting [Fiber](https://gofiber.io/), [Viper](https://github.com/spf13/viper), [PostgreSQL](https://www.postgresql.org/), and [Google Cloud Platform](https://cloud.google.com/)
- Project structure and DDD approach based on community best practices.
