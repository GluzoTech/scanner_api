package scan

import "time"

type CreateScanRequest struct {
	EventID      string `json:"event_id" binding:"required,uuid"`
	DeviceID     string `json:"device_id" binding:"required"`
	VendorID     string `json:"vendor_id" binding:"required"`
	ProductID    string `json:"product_id" binding:"required"`
	AgentHost    string `json:"agent_host" binding:"required"`
	BarcodeInput string `json:"barcode" binding:"required"`
	TimestampUTC string `json:"timestamp_utc" binding:"required"`
}

type Scan struct {
	ID           int64
	EventID      string
	DeviceID     string
	VendorID     string
	ProductID    string
	AgentHost    string
	BarcodeInput string
	CreatedAt    time.Time
	CapturedAt   time.Time
}
