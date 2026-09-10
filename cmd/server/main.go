package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"scanner-api/app/scan"
	"scanner-api/app/scanevent"
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
	scanEventRepository := scanevent.NewRepository(db)

	// Service
	scanService := scan.NewService(scanRepository)
	scanEventService := scanevent.NewService(scanEventRepository)

	// Handler
	scanHandler := scan.NewHandler(scanService)
	scanEventHandler := scanevent.NewHandler(scanEventService)

	router := gin.Default()

	// Server rendered views live in TEMPLATES_DIR (default "templates"),
	// resolved the same way migrations/ is so an installed binary does not
	// depend on the directory it was started from.
	templatesGlob := filepath.Join(templatesDirectory(), "*.html")

	// LoadHTMLGlob panics when the pattern matches nothing, which on a fresh
	// deployment surfaces as a stack trace rather than the missing directory.
	// Check first so the failure names the path and how to change it.
	matches, err := filepath.Glob(templatesGlob)
	if err != nil {
		log.Fatalf("invalid templates pattern %s: %v", templatesGlob, err)
	}

	if len(matches) == 0 {
		log.Fatalf(
			"no templates found at %s (set TEMPLATES_DIR to override)",
			templatesGlob,
		)
	}

	router.LoadHTMLGlob(templatesGlob)

	router.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")

	scan.RegisterRoutes(api, scanHandler)

	// Web scanner routes, mounted on the api group so the endpoint is
	// /api/scan/event, which is what the deployed scanner posts to.
	scanevent.RegisterRoutes(api, scanEventHandler)

	// Page routes
	scan.RegisterPageRoutes(router, scanHandler)

	address := listenAddress()

	log.Printf("Scanner API running on %s", address)

	if err := router.Run(address); err != nil {
		log.Fatal(err)
	}
}

// templatesDirectory resolves where the .html views live, mirroring
// MIGRATIONS_DIR so both runtime asset paths are configured the same way.
func templatesDirectory() string {
	if dir := strings.TrimSpace(os.Getenv("TEMPLATES_DIR")); dir != "" {
		return dir
	}

	return "templates"
}

// listenAddress builds the address to bind from PORT, defaulting to 8000.
//
// PORT is accepted either bare ("8000") or already prefixed (":8000"), since
// hosts differ on which they hand over, and a full host:port ("0.0.0.0:8000")
// is passed through so a deployment can bind one interface.
func listenAddress() string {
	port := strings.TrimSpace(os.Getenv("PORT"))

	if port == "" {
		return ":8000"
	}

	if strings.Contains(port, ":") {
		return port
	}

	return ":" + port
}
