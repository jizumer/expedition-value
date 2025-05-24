# Domain Driven Design (DDD) Guidelines

This document outlines the general DDD principles and project structure adopted for this application.

## Core DDD Principles

*   **Ubiquitous Language:** We will strive to use a common, rigorous language shared between developers and domain experts (or product owners in this context) for all discussions, code, and documentation related to the domain.
*   **Bounded Contexts:** The application will be designed around clear bounded contexts. For this MVP, we have identified:
    *   **Investment Analysis Context:** Focused on companies, financial metrics, and scoring. (Represented by the `Company` aggregate).
    *   **Portfolio Management Context:** Focused on managing investment portfolios, positions, and rebalancing. (Represented by the `Portfolio` aggregate).
    *   Interactions between contexts will be explicitly defined.
*   **Aggregates:** Aggregates will be used to group entities and value objects that are treated as a single unit for data changes. Each aggregate will have a root entity and a clearly defined boundary. Transactions will not cross aggregate boundaries.
*   **Entities:** Objects with a distinct identity that persists through different states.
*   **Value Objects:** Immutable objects that describe characteristics. They are identified by their values, not a unique ID.
*   **Domain Events:** Significant occurrences within the domain will be modeled as domain events. These can be used to communicate changes between aggregates or to trigger side effects.
*   **Repositories:** Used to abstract the persistence of aggregates. Interfaces will be defined in the domain layer, with implementations in the infrastructure layer.
*   **Application Services:** These will orchestrate domain logic to fulfill specific use cases. They should be thin and primarily delegate work to domain objects.
*   **Domain Services:** For domain logic that doesn't naturally fit within an entity or value object.

## Project Structure (Go specific)

The project will follow a layered architecture, generally organized as follows:

*   **`./cmd`**: Contains the entry points for the application(s).
    *   `./cmd/server/main.go`: Example for an HTTP server application.
    *   `./cmd/cli/main.go`: Example for a command-line interface.
*   **`./pkg`**: Contains the core library code, organized by layer.
    *   **`./pkg/domain`**:
        *   Contains aggregate roots, entities, value objects, domain events, and repository interfaces.
        *   Sub-packages per aggregate (e.g., `pkg/domain/company`, `pkg/domain/portfolio`).
        *   Example: `pkg/domain/company/company.go`, `pkg/domain/company/repository.go`
    *   **`./pkg/application`**:
        *   Contains application services that implement use cases.
        *   These services use repository interfaces from the domain layer and orchestrate domain objects.
        *   Example: `pkg/application/company_service.go`
    *   **`./pkg/infrastructure`**:
        *   Contains concrete implementations of interfaces defined in the domain layer.
        *   Examples:
            *   `pkg/infrastructure/persistence/memory/company_repository.go` (In-memory repository)
            *   `pkg/infrastructure/persistence/postgres/company_repository.go` (PostgreSQL repository - for later)
            *   `pkg/infrastructure/http/company_handler.go` (HTTP handlers if building a web service)
            *   `pkg/infrastructure/ třetí_strany/financial_data_provider.go` (Clients for external services)
*   **`./docs`**: Contains project documentation.
    *   `./docs/domain`: DDD blueprints and aggregate definitions.
    *   `./docs/guidelines.md`: This file.
*   **`./scripts`**: Contains helper scripts for building, testing, deploying, etc.
*   **`go.mod`, `go.sum`**: Go module files.

## Development Process Notes

*   **Test-Driven Development (TDD):** Where practical, TDD is encouraged, especially for domain logic.
*   **Code Reviews:** All code should be reviewed before merging.
*   **Updating DDD Documents:** If design decisions made during implementation impact the definitions in `docs/domain/*.md` or these guidelines, the documents should be updated accordingly.

**Note to Reviewer:** Please review these guidelines and suggest any additions or modifications you deem necessary for the project.
