package handlers

import (
	"io"
	"net/http"

	"github.com/adipras/api-gladiatore/internal/generator"
	"github.com/adipras/api-gladiatore/internal/models"
	"github.com/adipras/api-gladiatore/internal/validator"
	"github.com/gin-gonic/gin"
)

type GeneratorHandler struct {
	validator *validator.Validator
	generator *generator.Generator
}

func NewGeneratorHandler(templatesPath, outputPath string) *GeneratorHandler {
	return &GeneratorHandler{
		validator: validator.New(),
		generator: generator.New(templatesPath, outputPath),
	}
}

func (h *GeneratorHandler) GenerateService(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	config, err := h.validator.ValidateJSON(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.generator.Generate(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *GeneratorHandler) ValidateConfig(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	config, err := h.validator.ValidateJSON(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":  true,
		"config": config,
	})
}

func (h *GeneratorHandler) GetExample(c *gin.Context) {
	example := models.ServiceConfig{
		Service: models.ServiceInfo{
			Name:        "user-service",
			Description: "User management microservice",
			Package:     "github.com/example/user-service",
			Port:        8081,
		},
		Database: models.DatabaseConfig{
			Type:     "mysql",
			Host:     "localhost",
			Port:     3306,
			Database: "userdb",
			Username: "root",
			Password: "password",
		},
		Tables: []models.TableConfig{
			{
				Name: "users",
				Fields: []models.FieldConfig{
					{Name: "id", Type: "int64", PrimaryKey: true},
					{Name: "username", Type: "string", Required: true, Unique: true},
					{Name: "email", Type: "string", Required: true, Unique: true},
					{Name: "password", Type: "string", Required: true},
					{Name: "full_name", Type: "string"},
					{Name: "active", Type: "bool", Default: true},
				},
			},
		},
		Endpoints: []models.EndpointConfig{
			{Method: "POST", Path: "/users", Handler: "CreateUser", Table: "users", Operation: "create"},
			{Method: "GET", Path: "/users", Handler: "ListUsers", Table: "users", Operation: "list"},
			{Method: "GET", Path: "/users/:id", Handler: "GetUser", Table: "users", Operation: "read"},
			{Method: "PUT", Path: "/users/:id", Handler: "UpdateUser", Table: "users", Operation: "update"},
			{Method: "DELETE", Path: "/users/:id", Handler: "DeleteUser", Table: "users", Operation: "delete"},
		},
	}

	c.JSON(http.StatusOK, example)
}