package scan

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	scans := router.Group("/scans")

	scans.POST("", handler.CreateScan)
	scans.GET("", handler.SearchScans)
}

// RegisterPageRoutes mounts the server rendered views.
func RegisterPageRoutes(router gin.IRouter, handler *Handler) {
	router.GET("/scans", handler.SearchScanPage)
}
