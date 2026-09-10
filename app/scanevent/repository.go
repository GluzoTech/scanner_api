package scanevent

import (
	"context"
	"database/sql"
	"errors"
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

// Create stores a scan, or reports that it was already stored.
//
// event_id is the phone's own per-event UUID, so a scan replayed after a
// network timeout — or drained from an offline outbox days later — collides
// with the row already present and is dropped rather than duplicated.
func (r *Repository) Create(
	ctx context.Context,
	eventID string,
	referenceModel string,
	barcodeInput string,
	barcodeFormat string,
	scannedAt time.Time,
) (CreateResult, error) {
	query := `
	INSERT INTO scan_events (
	event_id,
	reference_model,
	barcode_input,
	barcode_format,
	scanned_at
	)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (event_id) DO NOTHING
	RETURNING id, received_at
	`

	// An absent format is stored as NULL rather than an empty string, so
	// "the decoder did not say" and "the decoder said nothing" stay distinct.
	var format any
	if barcodeFormat != "" {
		format = barcodeFormat
	}

	result := CreateResult{EventID: eventID}

	err := r.db.QueryRowContext(
		ctx,
		query,
		eventID,
		referenceModel,
		barcodeInput,
		format,
		scannedAt.UTC(),
	).Scan(&result.ID, &result.ReceivedAt)

	switch {
	case err == nil:
		return result, nil

	case errors.Is(err, sql.ErrNoRows):
		// DO NOTHING suppresses the RETURNING row, so no rows means the
		// event_id was already present. Read the original row back so the
		// caller can answer with the id it was first given.
		return r.findByEventID(ctx, eventID)

	default:
		return CreateResult{}, fmt.Errorf("failed to create scan event: %w", err)
	}
}

func (r *Repository) findByEventID(
	ctx context.Context,
	eventID string,
) (CreateResult, error) {
	query := `
	SELECT id, received_at
	FROM scan_events
	WHERE event_id = $1
	`

	result := CreateResult{EventID: eventID, Duplicate: true}

	if err := r.db.QueryRowContext(ctx, query, eventID).Scan(
		&result.ID,
		&result.ReceivedAt,
	); err != nil {
		return CreateResult{}, fmt.Errorf("failed to read existing scan event: %w", err)
	}

	return result, nil
}
