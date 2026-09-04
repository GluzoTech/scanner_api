package scan

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateScan(
	ctx context.Context,
	request CreateScanRequest,
) error {

	timestamp, err := time.Parse(
		time.RFC3339Nano,
		request.TimestampUTC,
	)

	if err != nil {
		return fmt.Errorf("Invalid timestamp_utc: %w", err)
	}

	return s.repository.Create(
		ctx,
		request.EventID,
		request.DeviceID,
		request.VendorID,
		request.ProductID,
		request.AgentHost,
		request.BarcodeInput,
		timestamp,
	)
}

// FindByOrderID looks a scan up by the value the scanner captured. "Order id"
// is the label the page view puts on it; in storage it is barcode_input.
func (s *Service) FindByOrderID(
	ctx context.Context,
	orderID string,
) ([]ScanDetail, error) {

	orderID = strings.TrimSpace(orderID)

	if orderID == "" {
		return nil, fmt.Errorf("order id is required")
	}

	return s.repository.FindByBarcode(ctx, orderID)
}
