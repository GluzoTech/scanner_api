package scanevent

import (
	"errors"
	"log"
	"net/http"
	"time"

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

func (h *Handler) CreateScanEvent(c *gin.Context) {
	var request CreateScanEventRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	result, err := h.service.CreateScanEvent(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, ErrInvalidRequest) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		// Anything else is ours, and must answer 5xx so the scanner keeps the
		// scan queued and retries it rather than discarding it as rejected.
		log.Printf("scan event %s failed: %v", request.EventID, err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "could not store the scan, please retry",
		})
		return
	}

	// 201 for a new row, 200 for one already stored. A replay is a success:
	// the scan is safely recorded, which is all the client needs to know to
	// drop it from its outbox.
	status := http.StatusCreated
	message := "Scan event captured successfully"

	if result.Duplicate {
		status = http.StatusOK
		message = "Scan event already recorded"
	}

	c.JSON(status, gin.H{
		"success":     true,
		"message":     message,
		"id":          result.ID,
		"event_id":    result.EventID,
		"duplicate":   result.Duplicate,
		"received_at": result.ReceivedAt.UTC().Format(time.RFC3339Nano),
	})
}
