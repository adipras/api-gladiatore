# ⚔️ API Gladiatore

A powerful microservice generation and management platform that transforms JSON configurations into production-ready Golang microservices.

## 🚀 Features

- **JSON-to-Microservice Generation**: Define your service structure in JSON and get a complete Go microservice
- **Service Management Dashboard**: Beautiful React SPA for managing generated services
- **CRUD Operations**: Automatic generation of Create, Read, Update, Delete endpoints
- **MySQL Integration**: Built-in database connection and migration support
- **Real-time Status Management**: Enable/disable services and endpoints on the fly
- **Docker Support**: Generated services include Dockerfile and docker-compose
- **API Documentation**: Auto-generated OpenAPI/Swagger documentation
- **Modern UI**: Tailwind CSS styled interface with smooth animations

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 16 or higher
- MySQL 5.7 or higher
- Git

## 🛠 Installation

1. **Clone the repository:**
```bash
git clone https://github.com/adipras/api-gladiatore.git
cd api-gladiatore
```

2. **Install dependencies:**
```bash
make deps
```

3. **Set up the database:**
```bash
make db-setup
```

4. **Configure environment:**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

5. **Build the project:**
```bash
make setup
```

## 🏃‍♂️ Running the Application

### Development Mode

**Run both frontend and backend in development mode:**
```bash
# Terminal 1 - Backend
make backend-dev

# Terminal 2 - Frontend
make frontend-dev
```

### Production Mode

```bash
make run
```

The application will be available at:
- Frontend: http://localhost:8080
- API: http://localhost:8080/api/v1

## 📝 Usage

### 1. Create a New Service

Click "New Service" and paste your JSON configuration:

```json
{
  "service": {
    "name": "user-service",
    "description": "User management microservice",
    "package": "github.com/example/user-service",
    "port": 8081
  },
  "database": {
    "type": "mysql",
    "host": "localhost",
    "port": 3306,
    "database": "userdb",
    "username": "root",
    "password": "password"
  },
  "tables": [
    {
      "name": "users",
      "fields": [
        {"name": "id", "type": "int64", "primary_key": true},
        {"name": "username", "type": "string", "required": true, "unique": true},
        {"name": "email", "type": "string", "required": true, "unique": true},
        {"name": "password", "type": "string", "required": true},
        {"name": "active", "type": "bool", "default": true}
      ]
    }
  ],
  "endpoints": [
    {"method": "POST", "path": "/users", "handler": "CreateUser", "table": "users", "operation": "create"},
    {"method": "GET", "path": "/users", "handler": "ListUsers", "table": "users", "operation": "list"},
    {"method": "GET", "path": "/users/:id", "handler": "GetUser", "table": "users", "operation": "read"},
    {"method": "PUT", "path": "/users/:id", "handler": "UpdateUser", "table": "users", "operation": "update"},
    {"method": "DELETE", "path": "/users/:id", "handler": "DeleteUser", "table": "users", "operation": "delete"}
  ]
}
```

### 2. Evaluate the Configuration

Click "Evaluate" to:
- Validate your JSON configuration
- Generate the complete microservice code
- Save service metadata to the database

### 3. Manage Services

From the dashboard, you can:
- **Edit**: Modify service configuration and regenerate
- **Details**: View and manage individual endpoints
- **Enable/Disable**: Toggle service and endpoint status
- **Delete**: Remove services from the system

## 🏗 Architecture

```
api-gladiatore/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── generator/       # Code generation engine
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── repository/      # Data access layer
│   ├── templates/       # Go templates for generation
│   └── validator/       # JSON validation
├── frontend/            # React SPA
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── api/         # API client
│   │   └── types/       # TypeScript types
│   └── public/
├── generated/           # Generated microservices output
└── migrations/          # Database migrations
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the root directory:

```env
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=root
DATABASE_PASS=yourpassword
DATABASE_NAME=api_gladiatore

# Paths
GENERATED_PATH=./generated
TEMPLATES_PATH=./internal/templates

# CORS
ALLOW_CORS=true
```

## 📚 API Documentation

### Service Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/services` | List all services |
| GET | `/api/v1/services/:id` | Get service details |
| POST | `/api/v1/services` | Create new service |
| PUT | `/api/v1/services/:id` | Update service |
| DELETE | `/api/v1/services/:id` | Delete service |
| PATCH | `/api/v1/services/:id/status` | Toggle service status |

### Generator Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/validate` | Validate JSON configuration |
| POST | `/api/v1/generate` | Generate service (without saving) |
| GET | `/api/v1/example` | Get example configuration |

## 🐳 Docker Support

Build and run with Docker:

```bash
# Build image
docker build -t api-gladiatore .

# Run container
docker run -p 8080:8080 api-gladiatore
```

## 🧪 Testing

Run tests:
```bash
make test
```

## 📦 Generated Service Structure

Each generated service includes:
- Complete Go project structure
- Database models with GORM
- RESTful API endpoints
- MySQL integration
- Docker configuration
- README documentation
- Migration scripts

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Gin Web Framework
- React and Tailwind CSS
- GORM ORM library
- Heroicons

## 📞 Support

For issues and questions, please open an issue on GitHub.

---

Made with ❤️ by [adipras](https://github.com/adipras)