<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Copilot Instructions untuk API IRAS

## Context
Ini adalah project REST API menggunakan Golang dengan framework Gin dan GORM. Project mengikuti Clean Architecture pattern dengan struktur folder yang jelas.

## Architecture Patterns
- **Clean Architecture**: Pemisahan yang jelas antara layers (presentation, business, data)
- **Repository Pattern**: Implemented dalam services layer
- **Dependency Injection**: Services di-inject ke controllers
- **Middleware Pattern**: Untuk cross-cutting concerns (CORS, Auth, Logging)

## Code Guidelines

### Naming Conventions
- Package names: lowercase, single word
- Function names: CamelCase untuk exported, camelCase untuk internal
- Variable names: camelCase
- Constants: UPPER_CASE atau CamelCase
- Struct names: CamelCase

### Error Handling
- Selalu return error sebagai nilai terakhir
- Gunakan fmt.Errorf() untuk wrap errors
- Handle errors secara eksplisit, jangan ignore
- Gunakan APIResponse struct untuk consistent response format

### Database Operations
- Gunakan GORM untuk semua database operations
- Implement proper transaction handling untuk complex operations
- Gunakan preloading untuk relationships yang diperlukan
- Implement soft deletes menggunakan gorm.DeletedAt

### API Design
- Follow RESTful principles
- Gunakan proper HTTP status codes
- Implement pagination untuk list endpoints
- Consistent response format menggunakan APIResponse struct
- Validate input menggunakan validator package

### Security
- Implement proper authentication middleware
- Sanitize input data
- Use environment variables untuk sensitive data
- Implement rate limiting untuk production

## Project Structure
```
- cmd/server/: Entry point aplikasi
- internal/config/: Configuration dan database setup
- internal/models/: Data models dan structs
- internal/services/: Business logic layer
- internal/controllers/: HTTP handlers
- internal/middleware/: HTTP middleware
- internal/routes/: Route definitions
- pkg/utils/: Reusable utility functions
```

## Development Practices
- Gunakan dependency injection
- Write unit tests untuk business logic
- Implement proper logging
- Use context.Context untuk cancellation
- Follow Go best practices dan conventions

## Dependencies
- gin-gonic/gin: Web framework
- gorm.io/gorm: ORM
- gorm.io/driver/postgres: PostgreSQL driver
- joho/godotenv: Environment variables
- go-playground/validator: Input validation
- golang.org/x/crypto: Cryptography utilities
