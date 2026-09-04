package scan

import "time"

type CreateScanRequest struct {
	DeviceID     string `json:"device_id" binding:"required"`
	BarcodeInput string `json:"barcode_input" binding:"required"`
	TimestampUTC string `json:"timestampt_utc" binding:"required"`
}

type Scan struct {
	ID           int64
	DeviceID     string
	BarcodeInput string
	CreatedAt    time.Time
	CapturedAt   time.Time
}
