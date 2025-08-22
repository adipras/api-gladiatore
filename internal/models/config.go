package models

type ServiceConfig struct {
	Service  ServiceInfo      `json:"service" validate:"required"`
	Database DatabaseConfig   `json:"database" validate:"required"`
	Tables   []TableConfig    `json:"tables" validate:"required,min=1"`
	Endpoints []EndpointConfig `json:"endpoints" validate:"required,min=1"`
}

type ServiceInfo struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Package     string `json:"package" validate:"required"`
	Port        int    `json:"port" validate:"required,min=1024,max=65535"`
}

type DatabaseConfig struct {
	Type     string `json:"type" validate:"required,eq=mysql"`
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	Database string `json:"database" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password"`
}

type TableConfig struct {
	Name   string        `json:"name" validate:"required"`
	Fields []FieldConfig `json:"fields" validate:"required,min=1"`
}

type FieldConfig struct {
	Name       string            `json:"name" validate:"required"`
	Type       string            `json:"type" validate:"required"`
	PrimaryKey bool              `json:"primary_key"`
	Required   bool              `json:"required"`
	Unique     bool              `json:"unique"`
	Default    interface{}       `json:"default"`
	Validation map[string]string `json:"validation"`
}

type EndpointConfig struct {
	Method      string                 `json:"method" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	Path        string                 `json:"path" validate:"required"`
	Handler     string                 `json:"handler" validate:"required"`
	Table       string                 `json:"table"`
	Operation   string                 `json:"operation" validate:"required,oneof=create read update delete list custom"`
	Auth        bool                   `json:"auth"`
	Middleware  []string               `json:"middleware"`
	RequestBody map[string]interface{} `json:"request_body"`
	Response    map[string]interface{} `json:"response"`
}

type GenerationResult struct {
	Success      bool     `json:"success"`
	ServiceName  string   `json:"service_name"`
	OutputPath   string   `json:"output_path"`
	Files        []string `json:"files"`
	Endpoints    []string `json:"endpoints"`
	Message      string   `json:"message"`
	Error        string   `json:"error,omitempty"`
}