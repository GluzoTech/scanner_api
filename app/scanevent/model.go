// Package scanevent stores scans captured by the web scanner.
//
// It sits alongside app/scan rather than inside it: that package owns the
// hardware agent's `scans` table, which identifies a capture by the device
// that produced it. A web scan is identified by the reference model the
// operator set before scanning, and shares no columns beyond the barcode
// itself.
package scanevent

import "time"

// CreateScanEventRequest is the body posted to POST /scan/event.
//
// The field names match what the web scanner already sends, so the client
// needs no change. Note `timestamp` rather than app/scan's `timestamp_utc`:
// the deployed frontend spells it this way, and renaming it here would break
// every scan already queued in a phone's outbox.
type CreateScanEventRequest struct {
	EventID string `json:"event_id" binding:"required,uuid"`
	// RFC3339 UTC, read from the phone's clock at the moment of decode.
	Timestamp      string `json:"timestamp" binding:"required"`
	ReferenceModel string `json:"reference_model" binding:"required,max=255"`
	BarcodeInput   string `json:"barcode_input" binding:"required,max=4096"`
	// Optional: the symbology the decoder reported, e.g. ean_13, code_128.
	BarcodeFormat string `json:"barcode_format" binding:"omitempty,max=32"`
}

// ScanEvent is a stored scan.
type ScanEvent struct {
	ID             int64
	EventID        string
	ReferenceModel string
	BarcodeInput   string
	BarcodeFormat  string
	ScannedAt      time.Time
	ReceivedAt     time.Time
}

// CreateResult reports what a write did.
//
// Duplicate carries the outcome the offline queue depends on: replaying a
// scan that already landed is a success, not an error, and the client uses
// this to retire the row from its outbox instead of retrying it forever.
type CreateResult struct {
	ID         int64
	EventID    string
	ReceivedAt time.Time
	Duplicate  bool
}
