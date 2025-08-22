# Changelog

All notable changes to API Gladiatore will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.1] - 2025-08-22

### Changed
- Initial release of API Gladiatore

### Summary
- Files added: 0
- Files modified: 2
- Files removed: 0


## [v1.0.0] - 2025-08-22

### Changed
- Minor updates and improvements

### Summary
- Files added: 0
- Files modified: 0
- Files removed: 0


## [Unreleased]

### Added
- Initial release of API Gladiatore platform
- JSON-to-Microservice generation engine with Go templates
- Service management dashboard with React SPA and Tailwind CSS
- Real-time service and endpoint status management (boolean stored as 0/1 in DB)
- MySQL database integration with GORM ORM
- JSON configuration validation system
- CRUD endpoint generation for microservices
- Service enable/disable functionality
- Endpoint management with individual toggle controls
- JSON evaluation panel (not a form) with syntax highlighting
- Docker support for generated services
- Automated CHANGELOG generation in release script
- Frontend with modern UI/UX using Tailwind CSS animations
- CORS support for API endpoints
- Health check endpoint
- Service repository pattern for data access
- Template-based code generation system

### Features Implemented
- **Backend**: Go with Gin framework, GORM ORM, MySQL
- **Frontend**: React with TypeScript, Tailwind CSS, Heroicons
- **Code Generation**: Go templates for microservice scaffolding
- **API Management**: Full CRUD operations for services
- **Status Management**: Boolean status with database persistence
- **Configuration**: JSON-based service configuration with validation

### Technical Details
- Modular architecture with clean separation of concerns
- Repository pattern for database operations
- Validation layer for JSON configurations
- Template engine for code generation
- RESTful API design
- SPA frontend with React hooks
- Responsive design with Tailwind CSS


