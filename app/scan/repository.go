package scan

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

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
	deviceID string,
	barcodeInput string,
	capturedAt time.Time,
) error {
	query := `
	INSERT INTO scans (
	device_id
	barcode_input,
	created_at,
	captured_at
	)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		deviceID,
		barcodeInput,
		capturedAt.Format("2006-01-02"),
		capturedAt.Format("15:04:05.999999"),
	)

	if err != nil {
		return fmt.Errorf("failed to create scan: %w", err)
	}

	return nil
}
