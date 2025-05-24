# Value Investment Analysis and Portfolio Management MVP

## Goal

This application aims to provide a Minimum Viable Product (MVP) for tools related to:
*   **Value Investment Analysis:** Tracking financial metrics of companies and calculating value investment scores.
*   **Portfolio Management:** Managing investment positions and providing recommendations for automated rebalancing based on value scores.

## Proposed Tech Stack

*   **Backend Language:** Go
    *   **Modularity:** Organized into packages reflecting Domain Driven Design layers:
        *   `pkg/domain`: Core business logic, aggregates, entities, value objects, domain events, and repository interfaces.
        *   `pkg/application`: Application services that orchestrate domain logic to fulfill use cases.
        *   `pkg/infrastructure`: Concrete implementations of repository interfaces (e.g., in-memory for MVP, database adapters), external service clients, and UI-facing components (e.g., HTTP handlers, CLI commands).
    *   **Main Application(s):** Located in `cmd/` directory (e.g., `cmd/server` for an HTTP server, `cmd/cli` for a command-line interface).
*   **API (if applicable):** RESTful HTTP API or gRPC. (To be decided based on specific MVP features).
*   **Data Persistence (MVP):** In-memory storage to simplify initial development.
*   **Deployment:**
    *   Local: Direct execution of compiled binaries.
    *   Cloud: Containerized deployment (e.g., using Docker) suitable for platforms like Google Cloud Run, AWS Fargate, or Kubernetes.

**Note to Reviewer:** Please confirm if Go is the preferred backend language. If not, specify the desired language and framework, and I will adjust the plan and implementation accordingly.

## Initial Constraints

*   **MVP Scope:** Focus on core features for company analysis and basic portfolio management as defined in the DDD blueprints (`docs/domain/`).
*   **DDD Principles:** Development will adhere to Domain Driven Design principles, using the guidance in `docs/guidelines.md` and the aggregate designs in `docs/domain/`.
*   **Deployability:** The application must be runnable locally and be structured for straightforward cloud deployment.
*   **Documentation:** Instructions for running the MVP in both environments will be provided. DDD documents will be updated to reflect any design decisions made during development.

## How to Run Locally

### Prerequisites
*   Go (version 1.18 or higher recommended) installed on your system.
*   A Go workspace set up.

### Steps
1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone <repository_url>
    cd expedition-value 
    ```
    (Replace `<repository_url>` with the actual URL of the repository. The project module is `github.com/jizumer/expedition-value`.)

2.  **Navigate to the server directory:**
    ```bash
    cd cmd/server
    ```

3.  **Run the application:**
    ```bash
    go run main.go
    ```
    This will start the HTTP server, typically on port `8080`. You should see a log message indicating the server has started.

4.  **Access the API:**
    You can now access the API endpoints using a tool like `curl` or Postman. For example:
    *   Health check: `curl http://localhost:8080/health`
    *   Get company (example, requires data to be created first via API): `curl http://localhost:8080/company?ticker=AAPL`
    *   Create company (example JSON body):
        ```bash
        curl -X POST -H "Content-Type: application/json" -d '{"ticker":"MSFT","name":"Microsoft Corp"}' http://localhost:8080/company/create 
        ``` 
        (Note: The actual JSON structure for creating a company will depend on the `CreateCompanyRequest` DTO defined in the HTTP handlers and the fields required by the `CompanyService.CreateCompany` method. The example above is a generic placeholder.)

## Deployment to Cloud (Conceptual)

This application is designed to be deployable to cloud platforms that support Go applications or containers. Here's a general approach:

1.  **Containerize the Application (using Docker):**
    *   Create a `Dockerfile` in the root of the project. A basic Go Dockerfile might look like this:

        ```Dockerfile
        # Use an official Go runtime as a parent image
        FROM golang:1.21-alpine AS builder

        # Set the Current Working Directory inside the container
        WORKDIR /app

        # Copy go mod and sum files
        COPY go.mod go.sum ./

        # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
        RUN go mod download

        # Copy the source code into the container
        COPY . .

        # Build the Go app
        # Build for the server, adjust if you have multiple main packages
        RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server/main.go

        # Start a new stage from scratch for a GOCR-style image (Google Cloud Optimized Runtimes)
        # Or use a distroless image for a smaller footprint
        FROM gcr.io/distroless/static-debian11 
        # FROM alpine:latest #  Alternative if you need a shell in the final image

        WORKDIR /

        # Copy the Pre-built binary file from the previous stage
        COPY --from=builder /app/server /server

        # Expose port 8080 to the outside world
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
