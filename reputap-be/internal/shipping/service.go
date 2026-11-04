package shipping

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("shipment not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Shipment, error) {
	shipment := &Shipment{
		OrderID:          req.OrderID,
		RecipientName:    req.RecipientName,
		RecipientAddress: req.RecipientAddress,
		RecipientPhone:   req.RecipientPhone,
		Courier:          req.Courier,
		TrackingNumber:   req.TrackingNumber,
	}

	if err := s.repo.Create(ctx, shipment); err != nil {
		return nil, fmt.Errorf("creating shipment: %w", err)
	}

	created, err := s.repo.GetByID(ctx, shipment.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching created shipment: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*Shipment, error) {
	shipment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching shipment: %w", err)
	}

	return shipment, nil
}

func (s *Service) GetByOrderID(ctx context.Context, orderID uint64) (*Shipment, error) {
	shipment, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching shipment by order: %w", err)
	}

	return shipment, nil
}

func (s *Service) List(ctx context.Context, params ListParams) ([]Shipment, int, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) UpdateStatus(ctx context.Context, id uint64, req UpdateStatusRequest) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching shipment: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, id, req.Status, req.TrackingNumber, req.Courier); err != nil {
		return fmt.Errorf("updating shipment status: %w", err)
	}

	return nil
}
