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

// ScanDetail is the read model returned by an order id lookup. created_at and
// captured_at are bare DATE / TIME columns, so they are read back as already
// formatted IST strings rather than being re-interpreted as instants.
type ScanDetail struct {
	ID           int64  `json:"id"`
	EventID      string `json:"event_id"`
	DeviceID     string `json:"device_id"`
	VendorID     string `json:"vendor_id"`
	ProductID    string `json:"product_id"`
	AgentHost    string `json:"agent_host"`
	BarcodeInput string `json:"barcode"`
	CreatedAt    string `json:"created_at"`
	CapturedAt   string `json:"captured_at"`
}
