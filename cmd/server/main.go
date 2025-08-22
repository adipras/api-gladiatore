package main

import (
	"log"
	"net/http"

	"github.com/adipras/api-gladiatore/internal/config"
	"github.com/adipras/api-gladiatore/internal/database"
	"github.com/adipras/api-gladiatore/internal/handlers"
	"github.com/adipras/api-gladiatore/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	serviceRepo := repository.NewServiceRepository(db.DB)

	// Initialize handlers
	serviceHandler := handlers.NewServiceHandler(cfg.TemplatesPath, cfg.GeneratedPath, serviceRepo)
	generatorHandler := handlers.NewGeneratorHandler(cfg.TemplatesPath, cfg.GeneratedPath)

	// Setup Gin router
	router := gin.Default()

	// Configure CORS
	if cfg.AllowCORS {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
		router.Use(cors.New(corsConfig))
	}

	// API routes
	api := router.Group("/api/v1")
	{
		// Service management endpoints
		api.GET("/services", serviceHandler.ListServices)
		api.GET("/services/:id", serviceHandler.GetService)
		api.GET("/services/:id/config", serviceHandler.GetServiceConfig)
		api.POST("/services", serviceHandler.GenerateService)
		api.PUT("/services/:id", serviceHandler.UpdateService)
		api.PATCH("/services/:id/status", serviceHandler.ToggleServiceStatus)
		api.DELETE("/services/:id", serviceHandler.DeleteService)

		// Endpoint management
		api.GET("/services/:id/endpoints", serviceHandler.GetServiceEndpoints)
		api.PATCH("/endpoints/:id/status", serviceHandler.ToggleEndpointStatus)

		// Generator utilities
		api.POST("/generate", generatorHandler.GenerateService)
		api.POST("/validate", generatorHandler.ValidateConfig)
		api.GET("/example", generatorHandler.GetExample)
	}

	// Serve static files (React app)
	router.Static("/static", "./frontend/build/static")
	router.StaticFile("/", "./frontend/build/index.html")
	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/build/index.html")
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "api-gladiatore",
		})
	})

	// Start server
	log.Printf("API Gladiatore is running on port %s", cfg.ServerPort)
	log.Printf("Frontend: http://localhost:%s", cfg.ServerPort)
	log.Printf("API Base: http://localhost:%s/api/v1", cfg.ServerPort)
	
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}