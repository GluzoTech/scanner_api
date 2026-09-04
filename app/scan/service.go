package scan

import (
	"context"
	"fmt"
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
		request.DeviceID,
		request.BarcodeInput,
		timestamp,
	)
}
