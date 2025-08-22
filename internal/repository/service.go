package repository

import (
	"encoding/json"

	"github.com/adipras/api-gladiatore/internal/models"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(config *models.ServiceConfig, outputPath string) (*models.Service, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	service := &models.Service{
		Name:          config.Service.Name,
		Description:   config.Service.Description,
		Package:       config.Service.Package,
		Port:          config.Service.Port,
		Status:        true,
		Configuration: string(configJSON),
		OutputPath:    outputPath,
	}

	if err := r.db.Create(service).Error; err != nil {
		return nil, err
	}

	// Create endpoints
	for _, ep := range config.Endpoints {
		endpoint := models.Endpoint{
			ServiceID: service.ID,
			Method:    ep.Method,
			Path:      ep.Path,
			Handler:   ep.Handler,
			Table:     ep.Table,
			Operation: ep.Operation,
			Status:    true,
			Auth:      ep.Auth,
		}
		if err := r.db.Create(&endpoint).Error; err != nil {
			return nil, err
		}
	}

	// Load endpoints
	r.db.Preload("Endpoints").First(service, service.ID)
	return service, nil
}

func (r *ServiceRepository) GetAll() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Preload("Endpoints").Find(&services).Error
	return services, err
}

func (r *ServiceRepository) GetByID(id uint) (*models.Service, error) {
	var service models.Service
	err := r.db.Preload("Endpoints").First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) GetByName(name string) (*models.Service, error) {
	var service models.Service
	err := r.db.Preload("Endpoints").Where("name = ?", name).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) Update(id uint, config *models.ServiceConfig) (*models.Service, error) {
	var service models.Service
	if err := r.db.First(&service, id).Error; err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	service.Name = config.Service.Name
	service.Description = config.Service.Description
	service.Package = config.Service.Package
	service.Port = config.Service.Port
	service.Configuration = string(configJSON)

	if err := r.db.Save(&service).Error; err != nil {
		return nil, err
	}

	// Delete old endpoints and create new ones
	r.db.Where("service_id = ?", id).Delete(&models.Endpoint{})

	for _, ep := range config.Endpoints {
		endpoint := models.Endpoint{
			ServiceID: service.ID,
			Method:    ep.Method,
			Path:      ep.Path,
			Handler:   ep.Handler,
			Table:     ep.Table,
			Operation: ep.Operation,
			Status:    true,
			Auth:      ep.Auth,
		}
		r.db.Create(&endpoint)
	}

	r.db.Preload("Endpoints").First(&service, id)
	return &service, nil
}

func (r *ServiceRepository) UpdateStatus(id uint, status bool) error {
	return r.db.Model(&models.Service{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ServiceRepository) Delete(id uint) error {
	// Delete endpoints first
	if err := r.db.Where("service_id = ?", id).Delete(&models.Endpoint{}).Error; err != nil {
		return err
	}
	// Delete service
	return r.db.Delete(&models.Service{}, id).Error
}

func (r *ServiceRepository) GetEndpointsByServiceID(serviceID uint) ([]models.Endpoint, error) {
	var endpoints []models.Endpoint
	err := r.db.Where("service_id = ?", serviceID).Find(&endpoints).Error
	return endpoints, err
}

func (r *ServiceRepository) UpdateEndpointStatus(endpointID uint, status bool) error {
	return r.db.Model(&models.Endpoint{}).Where("id = ?", endpointID).Update("status", status).Error
}

func (r *ServiceRepository) GetEndpointByID(id uint) (*models.Endpoint, error) {
	var endpoint models.Endpoint
	err := r.db.First(&endpoint, id).Error
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (r *ServiceRepository) CheckServiceExists(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Service{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
