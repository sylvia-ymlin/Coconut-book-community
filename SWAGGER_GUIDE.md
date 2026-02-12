# Swagger API Documentation Guide

## Overview

BookCommunity uses Swagger/OpenAPI 3.0 for API documentation. The documentation is automatically generated from code annotations and provides an interactive UI for testing endpoints.

## Access Swagger UI

### Local Development
```
http://localhost:8080/swagger/index.html
```

### Production
```
https://your-domain.com/swagger/index.html
```

## Features

✅ **Interactive API Testing** - Test endpoints directly in the browser
✅ **Auto-generated** - Documentation generated from code annotations
✅ **OpenAPI 3.0 Standard** - Industry-standard specification
✅ **JWT Authentication** - Built-in auth testing support
✅ **Request/Response Examples** - See example payloads

## Available Endpoints

### User Management
- `POST /douyin/user/register/` - Register new user
- `POST /douyin/user/login/` - User login
- `GET /douyin/user/` - Get user information

### Recommendations
- `GET /douyin/recommend` - Get personalized book recommendations
- `GET /douyin/search` - Search books by keyword
- `GET /douyin/book/:isbn` - Get book details

### System
- `GET /health` - Health check endpoint

## Authentication

Most endpoints require JWT authentication. To test authenticated endpoints:

1. **Register/Login** to get a token
2. Click **Authorize** button in Swagger UI
3. Enter: `Bearer <your-token>`
4. Click **Authorize**
5. Now you can test protected endpoints

## Generating Documentation

### Prerequisites
```bash
# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Docs
```bash
# Using Makefile
make swagger

# Or directly
swag init -g main.go --output ./docs
```

### Format Swagger Comments
```bash
make swagger-fmt
```

## Writing Swagger Annotations

### Handler Example

```go
// CreateBook creates a new book
// @Summary Create a new book
// @Description Add a new book to the database
// @Tags Book
// @Accept json
// @Produce json
// @Param book body BookDTO true "Book information"
// @Security BearerAuth
// @Success 200 {object} response.BookResponse
// @Failure 400 {object} response.CommonResponse
// @Failure 401 {object} response.CommonResponse
// @Router /book [post]
func CreateBookHandler(c *gin.Context) {
    // handler implementation
}
```

### Model Example

```go
// BookDTO represents book creation request
type BookDTO struct {
    Title  string `json:"title" binding:"required" example:"Go Programming"`
    Author string `json:"author" binding:"required" example:"John Doe"`
    ISBN   string `json:"isbn" binding:"required" example:"978-1234567890"`
} // @name BookDTO
```

### Common Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@Summary` | Short description | `@Summary Get user info` |
| `@Description` | Detailed description | `@Description Retrieves user by ID` |
| `@Tags` | Group endpoints | `@Tags User` |
| `@Accept` | Request content type | `@Accept json` |
| `@Produce` | Response content type | `@Produce json` |
| `@Param` | Parameter definition | `@Param id path int true "User ID"` |
| `@Success` | Success response | `@Success 200 {object} UserResponse` |
| `@Failure` | Error response | `@Failure 404 {object} ErrorResponse` |
| `@Router` | Route path | `@Router /user/{id} [get]` |
| `@Security` | Auth requirement | `@Security BearerAuth` |

### Parameter Types

```go
// Query parameter
// @Param name query string false "User name"

// Path parameter
// @Param id path int true "User ID"

// Body parameter
// @Param user body UserDTO true "User data"

// Header parameter
// @Param Authorization header string true "Bearer token"
```

## Project Structure

```
BookCommunity/
├── docs/                    # Generated Swagger files
│   ├── docs.go             # Generated Go code
│   ├── swagger.json        # OpenAPI JSON spec
│   └── swagger.yaml        # OpenAPI YAML spec
├── main.go                 # Main API info annotations
└── internal/
    └── app/
        └── handlers/       # Handler annotations
```

## CI/CD Integration

### GitHub Actions

The CI pipeline automatically validates Swagger docs:

```yaml
# .github/workflows/ci.yaml
- name: Validate Swagger
  run: |
    swag init -g main.go --output ./docs
    git diff --exit-code docs/
```

### Pre-commit Hook

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/sh
make swagger
git add docs/
```

## Best Practices

### 1. Keep Annotations Updated
✅ Update Swagger comments when changing handlers
❌ Don't let docs drift from implementation

### 2. Use Meaningful Descriptions
✅ `@Summary Get user's reading history`
❌ `@Summary Get data`

### 3. Document All Responses
```go
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Not found"
// @Failure 500 {object} ErrorResponse "Internal error"
```

### 4. Group Related Endpoints
```go
// @Tags User
// @Tags Book
// @Tags Recommendation
```

### 5. Provide Examples
```go
type BookDTO struct {
    Title  string `json:"title" example:"The Go Programming Language"`
    Author string `json:"author" example:"Alan Donovan"`
    Price  float64 `json:"price" example:"45.99"`
}
```

## Common Issues & Solutions

### Issue: Docs not updating
```bash
# Clear and regenerate
rm -rf docs/
make swagger
```

### Issue: Model not showing in Swagger UI
```go
// Add @name annotation
type UserDTO struct {
    // ...
} // @name UserDTO
```

### Issue: Enum values not showing
```go
// Use @Enum annotation
// @Param status query string false "User status" Enums(active, inactive, banned)
```

### Issue: Authentication not working
```go
// Ensure security definition in main.go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// Apply to handler
// @Security BearerAuth
```

## Testing API with Swagger

### 1. Health Check
```
GET /health
```
Expected: `{"status": "healthy"}`

### 2. Register User
```
POST /douyin/user/register/
Query params:
  username: testuser
  password: testpass123
```

### 3. Get Recommendations
```
GET /douyin/recommend?top_k=5
Headers:
  Authorization: Bearer <token>
```

### 4. Search Books
```
GET /douyin/search?q=golang&top_k=10
```

## Export Documentation

### JSON
```bash
curl http://localhost:8080/swagger/doc.json > api-spec.json
```

### YAML
```
cp docs/swagger.yaml api-spec.yaml
```

### Postman Collection

1. Visit Swagger UI
2. Export OpenAPI spec
3. Import to Postman

## Advanced Features

### Custom Headers
```go
// @Param X-Request-ID header string false "Request ID"
```

### File Upload
```go
// @Param file formData file true "File to upload"
// @Accept multipart/form-data
```

### Array Responses
```go
// @Success 200 {array} BookDTO
```

### Nested Objects
```go
// @Success 200 {object} response.DataResponse{data=BookDTO}
```

## Comparison with Alternatives

| Tool | Pros | Cons | BookCommunity |
|------|------|------|---------------|
| **Swagger** | Interactive UI, Standard | Learning curve | ✅ Used |
| Postman | Easy testing | Not code-first | ❌ |
| Redoc | Beautiful docs | No testing UI | ❌ |
| API Blueprint | Markdown-based | Less popular | ❌ |

## Resources

- **Swagger Official**: https://swagger.io/
- **swaggo/swag**: https://github.com/swaggo/swag
- **OpenAPI Spec**: https://spec.openapis.org/oas/v3.0.0
- **Gin Swagger**: https://github.com/swaggo/gin-swagger

## Next Steps

1. ✅ Access Swagger UI at `/swagger/index.html`
2. ✅ Test all endpoints
3. ⬜ Add more handler annotations
4. ⬜ Export API spec for client SDKs
5. ⬜ Share docs with frontend team

---

**📚 Happy documenting!**
