package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/adipras/api-gladiatore/internal/models"
)

type Generator struct {
	templatesPath string
	outputPath    string
}

func New(templatesPath, outputPath string) *Generator {
	return &Generator{
		templatesPath: templatesPath,
		outputPath:    outputPath,
	}
}

func (g *Generator) Generate(config *models.ServiceConfig) (*models.GenerationResult, error) {
	servicePath := filepath.Join(g.outputPath, config.Service.Name)
	
	if err := g.createDirectoryStructure(servicePath); err != nil {
		return nil, fmt.Errorf("failed to create directory structure: %v", err)
	}
	
	files := []string{}
	
	if err := g.generateMainFile(servicePath, config); err != nil {
		return nil, fmt.Errorf("failed to generate main.go: %v", err)
	}
	files = append(files, "main.go")
	
	if err := g.generateModelsFile(servicePath, config); err != nil {
		return nil, fmt.Errorf("failed to generate models: %v", err)
	}
	files = append(files, "internal/models/models.go")
	
	if err := g.generateHandlersFile(servicePath, config); err != nil {
		return nil, fmt.Errorf("failed to generate handlers: %v", err)
	}
	files = append(files, "internal/handlers/handlers.go")
	
	if err := g.generateDatabaseFile(servicePath, config); err != nil {
		return nil, fmt.Errorf("failed to generate database: %v", err)
	}
	files = append(files, "internal/database/database.go")
	
	if err := g.generateRouterFile(servicePath, config); err != nil {
		return nil, fmt.Errorf("failed to generate router: %v", err)
	}
	files = append(files, "internal/router/router.go")
	
	if err := g.generateGoMod(servicePath, config); err != nil {
		return nil, fmt.Errorf("failed to generate go.mod: %v", err)
	}
	files = append(files, "go.mod")
	
	if err := g.generateDockerfile(servicePath, config); err != nil {
		return nil, fmt.Errorf("failed to generate Dockerfile: %v", err)
	}
	files = append(files, "Dockerfile")
	
	if err := g.generateREADME(servicePath, config); err != nil {
		return nil, fmt.Errorf("failed to generate README: %v", err)
	}
	files = append(files, "README.md")
	
	endpoints := []string{}
	for _, ep := range config.Endpoints {
		endpoints = append(endpoints, fmt.Sprintf("%s %s", ep.Method, ep.Path))
	}
	
	return &models.GenerationResult{
		Success:     true,
		ServiceName: config.Service.Name,
		OutputPath:  servicePath,
		Files:       files,
		Endpoints:   endpoints,
		Message:     fmt.Sprintf("Service '%s' generated successfully", config.Service.Name),
	}, nil
}

func (g *Generator) createDirectoryStructure(servicePath string) error {
	dirs := []string{
		servicePath,
		filepath.Join(servicePath, "internal"),
		filepath.Join(servicePath, "internal", "models"),
		filepath.Join(servicePath, "internal", "handlers"),
		filepath.Join(servicePath, "internal", "database"),
		filepath.Join(servicePath, "internal", "router"),
		filepath.Join(servicePath, "internal", "middleware"),
		filepath.Join(servicePath, "migrations"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	
	return nil
}

func (g *Generator) generateMainFile(servicePath string, config *models.ServiceConfig) error {
	tmpl := `package main

import (
	"fmt"
	"log"

	"{{ .Service.Package }}/internal/database"
	"{{ .Service.Package }}/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	r := gin.Default()
	router.Setup(r, db)

	fmt.Printf("{{ .Service.Name }} is running on port {{ .Service.Port }}\n")
	if err := r.Run(":{{ .Service.Port }}"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
`
	return g.writeTemplate(filepath.Join(servicePath, "main.go"), tmpl, config)
}

func (g *Generator) generateModelsFile(servicePath string, config *models.ServiceConfig) error {
	tmpl := `package models

import (
	"time"
)

{{ range .Tables }}
type {{ .Name | title }} struct {
	{{ range .Fields -}}
	{{ .Name | title }} {{ .Type | goType }} ` + "`" + `json:"{{ .Name | lower }}"{{ if .PrimaryKey }} gorm:"primaryKey"{{ end }}{{ if .Unique }} gorm:"unique"{{ end }}` + "`" + `
	{{ end -}}
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}
{{ end }}
`
	return g.writeTemplate(filepath.Join(servicePath, "internal", "models", "models.go"), tmpl, config)
}

func (g *Generator) generateHandlersFile(servicePath string, config *models.ServiceConfig) error {
	tmpl := `package handlers

import (
	"net/http"
	"strconv"

	"{{ .Service.Package }}/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

{{ range .Endpoints }}
{{ if eq .Operation "create" }}
func (h *Handler) {{ .Handler }}(c *gin.Context) {
	var item models.{{ .Table | title }}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}
{{ else if eq .Operation "list" }}
func (h *Handler) {{ .Handler }}(c *gin.Context) {
	var items []models.{{ .Table | title }}
	if err := h.db.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, items)
}
{{ else if eq .Operation "read" }}
func (h *Handler) {{ .Handler }}(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.{{ .Table | title }}
	
	if err := h.db.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	c.JSON(http.StatusOK, item)
}
{{ else if eq .Operation "update" }}
func (h *Handler) {{ .Handler }}(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.{{ .Table | title }}
	
	if err := h.db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, item)
}
{{ else if eq .Operation "delete" }}
func (h *Handler) {{ .Handler }}(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.db.Delete(&models.{{ .Table | title }}{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
{{ else }}
func (h *Handler) {{ .Handler }}(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Custom endpoint"})
}
{{ end }}
{{ end }}
`
	return g.writeTemplate(filepath.Join(servicePath, "internal", "handlers", "handlers.go"), tmpl, config)
}

func (g *Generator) generateDatabaseFile(servicePath string, config *models.ServiceConfig) error {
	tmpl := `package database

import (
	"fmt"

	"{{ .Service.Package }}/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"{{ .Database.Username }}",
		"{{ .Database.Password }}",
		"{{ .Database.Host }}",
		{{ .Database.Port }},
		"{{ .Database.Database }}")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		{{ range .Tables -}}
		&models.{{ .Name | title }}{},
		{{ end -}}
	); err != nil {
		return nil, err
	}

	return db, nil
}
`
	return g.writeTemplate(filepath.Join(servicePath, "internal", "database", "database.go"), tmpl, config)
}

func (g *Generator) generateRouterFile(servicePath string, config *models.ServiceConfig) error {
	tmpl := `package router

import (
	"{{ .Service.Package }}/internal/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB) {
	h := handlers.New(db)

	api := r.Group("/api/v1")
	{
		{{ range .Endpoints -}}
		api.{{ .Method }}("{{ .Path }}", h.{{ .Handler }})
		{{ end -}}
	}
}
`
	return g.writeTemplate(filepath.Join(servicePath, "internal", "router", "router.go"), tmpl, config)
}

func (g *Generator) generateGoMod(servicePath string, config *models.ServiceConfig) error {
	tmpl := `module {{ .Service.Package }}

go 1.21

require (
	github.com/gin-gonic/gin v1.10.0
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)
`
	return g.writeTemplate(filepath.Join(servicePath, "go.mod"), tmpl, config)
}

func (g *Generator) generateDockerfile(servicePath string, config *models.ServiceConfig) error {
	tmpl := `FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o {{ .Service.Name }} .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/{{ .Service.Name }} .
EXPOSE {{ .Service.Port }}
CMD ["./{{ .Service.Name }}"]
`
	return g.writeTemplate(filepath.Join(servicePath, "Dockerfile"), tmpl, config)
}

func (g *Generator) generateREADME(servicePath string, config *models.ServiceConfig) error {
	tmpl := `# {{ .Service.Name }}

{{ .Service.Description }}

## API Endpoints

{{ range .Endpoints -}}
- ` + "`" + `{{ .Method }} {{ .Path }}` + "`" + ` - {{ .Handler }}
{{ end }}

## Database Tables

{{ range .Tables -}}
### {{ .Name }}
{{ range .Fields -}}
- {{ .Name }} ({{ .Type }}){{ if .PrimaryKey }} - Primary Key{{ end }}{{ if .Unique }} - Unique{{ end }}
{{ end }}
{{ end }}

## Setup

1. Install dependencies:
` + "```bash" + `
go mod download
` + "```" + `

2. Run the service:
` + "```bash" + `
go run main.go
` + "```" + `

3. Build Docker image:
` + "```bash" + `
docker build -t {{ .Service.Name }} .
` + "```" + `

## Configuration

The service runs on port {{ .Service.Port }} by default.

Database configuration:
- Host: {{ .Database.Host }}
- Port: {{ .Database.Port }}
- Database: {{ .Database.Database }}
`
	return g.writeTemplate(filepath.Join(servicePath, "README.md"), tmpl, config)
}

func (g *Generator) writeTemplate(filepath string, tmplStr string, data interface{}) error {
	funcMap := template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"goType": func(t string) string {
			typeMap := map[string]string{
				"string":   "string",
				"int":      "int",
				"int64":    "int64",
				"float":    "float64",
				"bool":     "bool",
				"datetime": "time.Time",
				"text":     "string",
				"json":     "json.RawMessage",
			}
			if goType, ok := typeMap[t]; ok {
				return goType
			}
			return "string"
		},
	}

	tmpl, err := template.New("template").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}