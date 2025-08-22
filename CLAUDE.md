# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

API Gladiatore is a microservice generation platform that takes JSON configuration as input and generates complete Golang microservices with CRUD operations, database integration, and API documentation.

**Current Status:** Requirements/Planning phase - implementation not yet started.

## Project Requirements

The platform should:
1. Accept JSON configuration defining service structure, database schema, and endpoints
2. Generate complete Golang microservices using Gin or Echo framework
3. Produce MySQL-compatible database schemas and migrations
4. Create CRUD endpoints with proper validation and error handling
5. Generate API documentation (OpenAPI/Swagger)

## Technology Stack

### Backend (Proposed)
- **Language:** Golang
- **Web Framework:** Gin or Echo (to be decided)
- **Database:** MySQL
- **Code Generation:** Go templates
- **API Documentation:** OpenAPI/Swagger

### Frontend (Proposed)
- **Option 1:** Vanilla HTML/CSS/JavaScript
- **Option 2:** React
- Web-based interface with JSON input textarea and evaluation functionality

## Development Setup

Since the project is not yet initialized, start with:

```bash
# Initialize Go module
go mod init github.com/api-gladiatore/api-gladiatore

# Install framework (choose one)
go get -u github.com/gin-gonic/gin  # For Gin
# OR
go get -u github.com/labstack/echo/v4  # For Echo

# Install MySQL driver
go get -u github.com/go-sql-driver/mysql

# Install other dependencies as needed
go get -u github.com/go-playground/validator/v10
```

## Architecture Guidelines

### Directory Structure (Recommended)
```
api-gladiatore/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── generator/       # Code generation logic
│   ├── validator/       # JSON schema validation
│   └── templates/       # Go templates for code generation
├── pkg/                 # Public packages
├── web/                 # Frontend files
├── migrations/          # Database migrations
└── generated/           # Output directory for generated services
```

### Key Implementation Considerations

1. **JSON Schema Validation:** Implement robust validation for input JSON configurations before processing
2. **Template Engine:** Use Go's `text/template` or `html/template` for code generation
3. **Database Design:** Store service configurations and metadata for tracking generated services
4. **Error Handling:** Provide clear error messages for invalid JSON or generation failures
5. **Security:** Sanitize all inputs and prevent code injection in generated services

## Core Features to Implement

1. **JSON Configuration Parser**
   - Validate against predefined schema
   - Extract service definitions, database tables, and endpoints

2. **Code Generator Engine**
   - Generate Go project structure
   - Create main.go with router setup
   - Generate handlers for CRUD operations
   - Create model structs from table definitions
   - Generate database connection and migration files

3. **Web Interface**
   - JSON input textarea with syntax highlighting
   - "Evaluate" button to trigger generation
   - Display generated service summary and download links

4. **Database Integration**
   - Store configuration history
   - Track generated services
   - Manage service metadata

## JSON Configuration Schema

The platform expects JSON with:
- Service name and description
- Database table definitions with fields and types
- CRUD endpoint specifications
- Field validations
- Authentication requirements
- Custom business logic hooks

## Testing Strategy

When implementing, ensure:
- Unit tests for JSON validation logic
- Integration tests for code generation
- End-to-end tests for the complete flow
- Generated code should be compilable and testable

## Important Notes

- The project is currently in planning phase with only requirements documented
- No existing codebase to maintain compatibility with
- Focus on clean, modular architecture for easy extension
- Generated microservices should follow Go best practices and be production-ready