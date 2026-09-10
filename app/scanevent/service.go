package scanevent

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrInvalidRequest marks a body the client got wrong, as opposed to a
// failure on our side.
//
// The distinction is load-bearing rather than cosmetic: the web scanner's
// retry queue treats a 4xx as permanent and discards the scan, and a 5xx as
// transient and retries it. Reporting a database outage as a 4xx would make
// every phone quietly throw away the scans it was holding.
var ErrInvalidRequest = errors.New("invalid request")

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateScanEvent(
	ctx context.Context,
	request CreateScanEventRequest,
) (CreateResult, error) {

	scannedAt, err := time.Parse(time.RFC3339Nano, request.Timestamp)
	if err != nil {
		return CreateResult{}, fmt.Errorf("%w: invalid timestamp: %v", ErrInvalidRequest, err)
	}

	// binding:"required" rejects an absent or empty field, but not one that is
	// entirely whitespace, which would otherwise become a reference model no
	// later lookup could match.
	referenceModel := strings.TrimSpace(request.ReferenceModel)
	if referenceModel == "" {
		return CreateResult{}, fmt.Errorf("%w: reference_model is required", ErrInvalidRequest)
	}

	barcodeInput := strings.TrimSpace(request.BarcodeInput)
	if barcodeInput == "" {
		return CreateResult{}, fmt.Errorf("%w: barcode_input is required", ErrInvalidRequest)
	}

	return s.repository.Create(
		ctx,
		request.EventID,
		referenceModel,
		barcodeInput,
		strings.TrimSpace(request.BarcodeFormat),
		scannedAt,
	)
}
