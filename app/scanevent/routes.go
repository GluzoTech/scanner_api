package scanevent

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the web scanner's endpoints.
//
// Takes gin.IRouter rather than *gin.RouterGroup so the caller decides the
// prefix: mounted on the engine these sit at /scan/event, mounted on the
// existing /api group they become /api/scan/event.
func RegisterRoutes(router gin.IRouter, handler *Handler) {
	events := router.Group("/scan")

	events.POST("/event", handler.CreateScanEvent)
}
