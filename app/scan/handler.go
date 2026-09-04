package scan

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateScan(c *gin.Context) {
	var request CreateScanRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.CreateScan(c.Request.Context(), request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Scan captured successfully",
	})
}

// orderIDQuery reads the searched value off the query string. The page view and
// the JSON endpoint both call it "order id"; "barcode" stays accepted so the
// value posted to CreateScan can be handed straight back for a lookup.
func orderIDQuery(c *gin.Context) string {
	orderID := strings.TrimSpace(c.Query("order_id"))

	if orderID == "" {
		orderID = strings.TrimSpace(c.Query("barcode"))
	}

	return orderID
}

func (h *Handler) SearchScans(c *gin.Context) {
	orderID := orderIDQuery(c)

	scans, err := h.service.FindByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(scans),
		"data":    scans,
	})
}

// SearchScanPage renders the search form and, once an order id has been
// entered, the scans captured for it.
func (h *Handler) SearchScanPage(c *gin.Context) {
	orderID := orderIDQuery(c)

	page := gin.H{
		"OrderID":  orderID,
		"Searched": orderID != "",
	}

	if orderID != "" {
		scans, err := h.service.FindByOrderID(c.Request.Context(), orderID)
		if err != nil {
			log.Printf("scan lookup for %q failed: %v", orderID, err)
			page["Error"] = "Could not fetch the scan details. Please try again."
		} else {
			page["Scans"] = scans
		}
	}

	c.HTML(http.StatusOK, "scan_search.html", page)
}
