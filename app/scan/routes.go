package scan

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	scans := router.Group("/scans")

	scans.POST("", handler.CreateScan)
}
