Help me build a microservice generation platform with the following requirements:

  Core Features:
  1. Web-based UI with a textarea for JSON configuration input and an "Evaluate" button
  2. Backend service that validates JSON input against predefined schemas
  3. Code generation engine that creates Golang microservices based on JSON config
  4. CRUD endpoint generator with customizable routes
  5. MySQL database integration for storing configurations and generated service metadata
  6. Response display showing generated service summary or sample API requests/responses

  Technical Stack:
  - Backend: Golang (Gin or Echo framework)
  - Database: MySQL
  - Frontend: HTML/CSS/JavaScript (or React)
  - Code generation: Go templates
  - API documentation: OpenAPI/Swagger

  JSON Configuration Schema should support:
  - Service name and description
  - Database table definitions
  - CRUD endpoint specifications
  - Field validations and types
  - Authentication requirements
  - Custom business logic hooks

  Expected Outputs:
  - For service generation: Complete Golang service with folder structure, main.go, handlers, models, and database connections
  - For endpoint generation: Sample curl requests and JSON responses
  - Service summary including: endpoints list, database schema, and configuration used

  Please create a detailed todo list for building this platform, including all necessary components, validation logic, code generation templates, and database
  schema design.