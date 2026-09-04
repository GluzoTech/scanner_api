package scan

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// IST is the Asia/Kolkata offset (+05:30). India does not observe DST, so a
// fixed zone is safe and avoids depending on the host tzdata being present.
var IST = time.FixedZone("IST", 5*60*60+30*60)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	eventID string,
	deviceID string,
	vendorID string,
	productID string,
	agentHost string,
	barcodeInput string,
	capturedAt time.Time,
) error {
	query := `
	INSERT INTO scans (
	event_id,
	device_id,
	vendor_id,
	product_id,
	agent_host,
	barcode_input,
	created_at,
	captured_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (event_id) DO NOTHING
	`

	// event_id is the agent's own per-event UUID, so a scan replayed after a
	// network timeout collides with the row already stored and is dropped
	// rather than duplicated.

	// created_at and captured_at are stored as bare DATE / TIME values, so the
	// instant has to be shifted into IST before it is formatted.
	capturedAtIST := capturedAt.In(IST)

	_, err := r.db.ExecContext(
		ctx,
		query,
		eventID,
		deviceID,
		vendorID,
		productID,
		agentHost,
		barcodeInput,
		capturedAtIST.Format("2006-01-02"),
		capturedAtIST.Format("15:04:05.999999"),
	)

	if err != nil {
		return fmt.Errorf("failed to create scan: %w", err)
	}

	return nil
}
