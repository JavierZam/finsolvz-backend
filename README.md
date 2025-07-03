# 🚀 Golang Clean Architecture Starterpack

This is a Golang backend starterpack built with **Clean Architecture** principles (also known as Hexagonal Architecture or Port and Adapters). This starterpack is designed to provide a solid, modular, and easy-to-maintain foundation for Go web applications, focusing on *testability*, *scalability*, and *developer experience*.

It is suitable for projects that require a well-organized structure from the outset and are flexible enough to evolve, including multi-tenancy implementation and potential microservices adaptation in the future.

## ✨ Key Features

* **Clean Architecture:** Clear project structure separating `domain`, `service` (usecase), `handler` (transport), and `repository` (persistence) layers for robust, testable, and maintainable code.
* **Go Modules:** Modern Go dependency management for reliable builds.
* **Structured Logging:** Uses `go.uber.org/zap` for consistent, high-performance, and structured logging, crucial for debugging and monitoring.
* **Custom Error Handling:** A structured custom error system (`internal/utils/errors`) with consistent mapping to HTTP status codes, ensuring uniform API error responses.
* **Pagination Helper:** Generic utility (`internal/utils/response`) and DTOs (Data Transfer Objects) for consistent pagination implementation across various list endpoints.
* **Database Ready:** Initial configuration for PostgreSQL with SQLX, equipped with robust schema setup via `golang-migrate/migrate` CLI for reliable database migrations.
    * **Note on Schema:** Uses `id UUID PRIMARY KEY` and `tenant_id UUID` for robust identification and multi-tenancy.
* **Docker Support:** `Dockerfile` for application containerization and Docker (via Makefile) for easy local PostgreSQL database setup and management.
* **Makefile:** Comprehensive automation scripts for common development tasks (build, test, lint, Docker commands, database management including full resets and migrations).
* **GCP Ready:** Architecture is designed to be highly compatible for seamless deployment to Google Cloud Platform services (e.g., Cloud Run, Cloud SQL, Kubernetes), with a `cloudbuild.yaml` example.
* **OpenAPI (Swagger) Documentation:** Includes the `api/` directory for API specifications (`openapi.yaml`), essential for generating and visualizing comprehensive API documentation.
* **Demo API: User Authentication Module:** A fully functional authentication module (Register, Login, Refresh Token) demonstrating the Clean Architecture pattern, JWT implementation, and `bcrypt` for password hashing.
* **Architectural Demo: Employee Module (Disabled by Default):** The `employee` module files are included in the repository for architectural demonstration purposes (showing how a module is structured), but it is **not wired into the application by default** as its schema is not compatible with the `users` table.

## 🛠️ Technology Stack

* **Programming Language:** Go (Golang)
* **Web Framework:** [Gorilla Mux](https://github.com/gorilla/mux)
* **Database:** PostgreSQL
* **ORM/Query Builder:** [SQLX](https://github.com/jmoiron/sqlx)
* **Logging:** [Zap](https://github.com/uber-go/zap)
* **Validation:** [go-playground/validator](https://go-playground.github.io/validator)
* **JWT:** [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
* **Password Hashing:** [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
* **UUID Generation:** [google/uuid](https://github.com/google/uuid)
* **Database Migrations:** [golang-migrate/migrate](https://golang-migrate.run/)
* **Containerization:** Docker
* **Linter:** [golangci-lint](https://golangci-lint.run/)

## 🚀 Getting Started (Local Setup)

Follow the steps below to run the starterpack on your local machine.

### Prerequisites

Ensure you have the following tools installed and properly configured:

* **Git**
* **Visual Studio Code (VS Code)** with **Go**, **Remote - WSL**, and **Docker** extensions.
* **Docker Desktop** with **WSL 2 backend enabled** and integration for your Linux distro.
* **WSL 2** and your preferred Linux Distribution (e.g., Ubuntu).
* **Go SDK** installed within your WSL environment (`go version`).
* **`make`**, `curl`, `tar` installed within your WSL environment.

### Setup Steps

1.  **Clone the Repository:**
    * Navigate to your desired directory in WSL.
    * ```bash
        git clone [https://github.com/yourusername/starterpack-golang-cleanarch.git](https://github.com/yourusername/starterpack-golang-cleanarch.git) # REPLACE WITH YOUR REPO URL
        cd starterpack-golang-cleanarch
        ```
    * **Note:** If adapting for a new project, you'd clone, delete `.git` folder, `git init`, and then adjust Go module name and `APP_NAME` in `Makefile`.

2.  **Install Go Dependencies:**
    * ```bash
        go mod tidy
        go clean -cache
        ```

3.  **Configure Environment Variables:**
    * Create a `.env` file from `.env.example`:
        ```bash
        cp .env.example .env
        # Edit .env as needed.
        ```
    * **IMPORTANT:** Ensure no trailing spaces in `.env` values.

4.  **Perform a Full Database Reset & Migrate:**
    * This command will stop/remove old containers, prune unused volumes (deleting old database data), start a fresh PostgreSQL container, and then apply all database migrations using `golang-migrate`.
    * ```bash
        make docker-reset-db
        ```
    * **Verify Database:** After this command, run `docker exec -it starterpack-golang-cleanarch-db psql -U starteruser -d starterdb -c "\dt"`. You should see `users` and `schema_migrations` tables.

5.  **Run the Go Application:**
    * ```bash
        make run
        ```
    * The application will build and start running at `http://localhost:8080`.

## 🧪 Testing the API

Once your application is running, you can test the available API endpoints.

### 1. Test General Endpoints (No Authentication Required)

* **Health Check:**
    ```bash
    curl http://localhost:8080/health
    ```
    Expected: `OK`
* **Server Info:**
    ```bash
    curl http://localhost:8080/info
    ```
    Expected: JSON response with app/server info.

### 2. Access OpenAPI (Swagger) Documentation

You can view the interactive API documentation using Swagger UI.

1.  **Ensure Docker Desktop is running.**
2.  **Run Swagger UI container:**
    ```bash
    docker run -p 8081:8080 -e SWAGGER_JSON=/app/openapi.yaml -v $(pwd)/api:/app swaggerapi/swagger-ui
    ```
3.  **Open your web browser** and navigate to: `http://localhost:8081`
    You will see the interactive API documentation based on your `api/openapi.yaml` file.

### 3. Test User Authentication Module (Auth Endpoints)

This module demonstrates user registration, login, and token refreshing. Refer to the **Swagger UI** at `http://localhost:8081` for detailed request/response schemas and examples for these endpoints:

* **`POST /auth/register`**: Register a new user.
* **`POST /auth/login`**: Log in a user and get JWT tokens.
* **`POST /auth/refresh`**: Refresh access token using a refresh token.
* **`GET /api/v1/user/me`**: Get current authenticated user's info (requires `access_token`).

Use `curl` or tools like Postman/Insomnia to test these endpoints. For authenticated endpoints, include the `access_token` in the `Authorization` header (e.g., `-H "Authorization: Bearer YOUR_ACCESS_TOKEN"`).

## 📂 Project Structure

This project structure adheres to Clean Architecture principles for clear modularity and separation of concerns:

![Project Architecture](https://github.com/JavierZam/starterpack-golang-cleanarch/blob/master/architecture.png?raw=true)

## ⚙️ Development Guide

### 1. Adapting for a Real Project (e.g., ADVIZ Intranet Backend)

This starterpack is designed to be a strong starting point. Here's what you'll need to adjust when transitioning to your actual ADVIZ project:

1.  **Git Repository & Go Module Name:**
    * After cloning this starterpack, you'll typically delete the `.git` folder (`rm -rf .git`) and initialize a new Git repository (`git init`).
    * Then, you'll link it to your actual ADVIZ Git repository (e.g., on GitLab/GitHub Enterprise).
    * **Crucially, you MUST change the Go module name** in `go.mod` (e.g., `module gitlab.adviz.co.id/intranet/backend`).
    * Perform a **global find-and-replace** (`starterpack-golang-cleanarch` -> `your-new-module-name`) across your entire codebase to update all import paths.
    * Update `APP_NAME` in `Makefile` to your actual application name.

2.  **Database Schema & Migrations:**
    * The `migrations/000001_create_auth_tables.up.sql` creates a `users` table suitable for authentication.
    * For other business domains (e.g., Clients, Projects, Tax Reports), you will **create new migration files** (e.g., `migrations/000002_create_clients_table.up.sql`).
    * **Consider UUID vs. SERIAL:** The current setup uses `UUID` for `id` and `tenant_id` in the `users` table. This is generally recommended for distributed systems. If your project requires `SERIAL PRIMARY KEY` (integer) for IDs, you'll need to adjust the migration SQL and corresponding Go types (`int64`) in `domain`, `repository`, `service`, and DTOs.

3.  **Environment Variables (`.env`):**
    * Update all database credentials (`DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_HOST`, `DB_PORT`) to match your actual development database.
    * **Generate a strong, random `JWT_SECRET`** for your project. Never use the default "this-is-a-super-secret-jwt-key..." in any environment beyond local development.
    * Adjust `JWT_EXPIRES_IN_MINUTES` and `REFRESH_TOKEN_EXPIRES_IN_HOURS` as per your security policy.

4.  **Implement Real Business Modules:**
    * You will create new modules under `internal/app/` (e.g., `internal/app/client`, `internal/app/project`, `internal/app/taxreport`).
    * For each new module, follow the consistent Clean Architecture pattern: `domain` entity, `repository` interface, `repository` implementation, `model` (DTOs), `errors` (custom module errors), `service` (business logic), and `handler` (API endpoints).
    * **Wiring:** In `cmd/server/main.go`, you will add the Dependency Injection (DI) wiring for these new modules to the `authenticatedRouter` (if they require authentication).

5.  **Authentication & Authorization:**
    * The `auth` module provides the core. You might need to expand it with features like password reset, email verification, or more granular role-based access control (RBAC) in `internal/platform/http/middleware/`.

6.  **OpenAPI (Swagger) Documentation:**
    * Update `api/openapi.yaml` to reflect your actual project's `info` (title, description, contact).
    * Add comprehensive definitions for all new API endpoints and their DTOs for each module you build.

7.  **CI/CD Pipeline (`cloudbuild.yaml`):**
    * Adjust the `cloudbuild.yaml` file to use your actual GCP Project ID, desired region, and specific service names for Cloud Run/Kubernetes.
    * **Crucially, configure secrets in Google Secret Manager** for `_DB_PASSWORD`, `_JWT_SECRET`, and any other sensitive credentials, and bind them to your Cloud Build trigger.

### 2. General Consistency Guidelines for Future Development

When implementing new features or modules, always adhere to these guidelines:

* **Layer Separation:**
    * **Domain:** Pure business entities and repository interfaces. No `http.Request`, `sql.DB`, or `json` tags.
    * **Repository:** Database interaction only. Implements `domain` interfaces. Handles `sql.ErrNoRows`.
    * **Service:** Business logic, validation, orchestration. Depends on `domain` interfaces. Maps `domain` to `app` DTOs.
    * **Handler:** HTTP interaction only. Parses requests, validates DTOs, calls `service`, formats responses. Depends on `app` DTOs and `service`.
* **`context.Context`:** Always pass `context.Context` as the first argument in method signatures across layers.
* **`TenantID`:** For multi-tenant data, ensure `tenantID` is extracted from `context.Context` in the handler and passed down to the service and repository layers, where it's used in all database queries (`WHERE tenant_id = ...`).
* **Error Handling:**
    * Use `internal/utils/errors.New(...)` for all custom application errors.
    * `fmt.Errorf("...: %w", err)` to wrap underlying errors.
    * `utils.HandleHTTPError` in handlers for consistent API error responses.
* **DTOs:** Define all request/response structs in `internal/app/[module_name]/model.go`. Use `json` and `validate` tags.
* **Validation:** Use `validator.Validate().Struct(req)` in handlers.
* **Logging:** Use `internal/utils/log` functions (`log.Info`, `log.Error`, `log.Debugf`) with `context.Context`.
* **Testing:** Write unit tests for `service` (mocking repositories) and integration tests for `repository` (with a real test DB) and `handler` (mocking service).
