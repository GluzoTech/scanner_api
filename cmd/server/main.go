package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"scanner-api/app/scan"
	"scanner-api/database"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Database migrations completed")

	// Repository
	scanRepository := scan.NewRepository(db)

	// Service
	scanService := scan.NewService(scanRepository)

	// Handler
	scanHandler := scan.NewHandler(scanService)

	router := gin.Default()

	// Server rendered views live in templates/, read relative to the working
	// directory the same way migrations/ is.
	router.LoadHTMLGlob("templates/*.html")

	router.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")

	scan.RegisterRoutes(api, scanHandler)

	// Page routes
	scan.RegisterPageRoutes(router, scanHandler)

	log.Println("Scanner API running on :8000")

	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
