package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adipras/api-gladiatore/internal/generator"
	"github.com/adipras/api-gladiatore/internal/models"
	"github.com/adipras/api-gladiatore/internal/repository"
	"github.com/adipras/api-gladiatore/internal/validator"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	validator  *validator.Validator
	generator  *generator.Generator
	repository *repository.ServiceRepository
}

func NewServiceHandler(templatesPath, outputPath string, repo *repository.ServiceRepository) *ServiceHandler {
	return &ServiceHandler{
		validator:  validator.New(),
		generator:  generator.New(templatesPath, outputPath),
		repository: repo,
	}
}

func (h *ServiceHandler) ListServices(c *gin.Context) {
	services, err := h.repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch services",
		})
		return
	}

	c.JSON(http.StatusOK, models.ServiceListResponse{
		Services: services,
		Total:    int64(len(services)),
	})
}

func (h *ServiceHandler) GetService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	service, err := h.repository.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) GenerateService(c *gin.Context) {
	var config models.ServiceConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate configuration
	configJSON, _ := json.Marshal(config)
	validatedConfig, err := h.validator.ValidateJSON(configJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if service already exists
	exists, _ := h.repository.CheckServiceExists(validatedConfig.Service.Name)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Service with this name already exists"})
		return
	}

	// Generate service
	result, err := h.generator.Generate(validatedConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save to database
	service, err := h.repository.Create(validatedConfig, result.OutputPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Service generated but failed to save to database",
			"result": result,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"service": service,
		"result":  result,
	})
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	var config models.ServiceConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate configuration
	configJSON, _ := json.Marshal(config)
	validatedConfig, err := h.validator.ValidateJSON(configJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Regenerate service files
	result, err := h.generator.Generate(validatedConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update database
	service, err := h.repository.Update(uint(id), validatedConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service": service,
		"result":  result,
	})
}

func (h *ServiceHandler) ToggleServiceStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	service, err := h.repository.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	newStatus := !service.Status

	if err := h.repository.UpdateStatus(uint(id), newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service status updated",
		"status":  newStatus,
	})
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	if err := h.repository.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

func (h *ServiceHandler) GetServiceEndpoints(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	endpoints, err := h.repository.GetEndpointsByServiceID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch endpoints"})
		return
	}

	c.JSON(http.StatusOK, models.EndpointListResponse{
		Endpoints: endpoints,
		Total:     int64(len(endpoints)),
	})
}

func (h *ServiceHandler) ToggleEndpointStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
		return
	}

	endpoint, err := h.repository.GetEndpointByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
		return
	}

	newStatus := !endpoint.Status

	if err := h.repository.UpdateEndpointStatus(uint(id), newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update endpoint status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Endpoint status updated",
		"status":  newStatus,
	})
}

func (h *ServiceHandler) GetServiceConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	service, err := h.repository.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	var config models.ServiceConfig
	if err := json.Unmarshal([]byte(service.Configuration), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}