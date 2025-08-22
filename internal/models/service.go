package models

import (
	"time"
)

type Service struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"unique;not null"`
	Description   string    `json:"description"`
	Package       string    `json:"package"`
	Port          int       `json:"port"`
	Status        bool      `json:"status" gorm:"default:true"`      // true=enabled, false=disabled
	Configuration string    `json:"configuration" gorm:"type:text"`  // JSON config stored as text
	OutputPath    string    `json:"output_path"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Endpoints     []Endpoint `json:"endpoints" gorm:"foreignKey:ServiceID"`
}

type Endpoint struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ServiceID   uint      `json:"service_id"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Handler     string    `json:"handler"`
	Table       string    `json:"table"`
	Operation   string    `json:"operation"`
	Status      bool      `json:"status" gorm:"default:true"` // true=enabled, false=disabled
	Auth        bool      `json:"auth"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ServiceListResponse struct {
	Services []Service `json:"services"`
	Total    int64     `json:"total"`
}

type EndpointListResponse struct {
	Endpoints []Endpoint `json:"endpoints"`
	Total     int64      `json:"total"`
}