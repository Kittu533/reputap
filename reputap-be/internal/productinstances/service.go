package productinstances

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kittu533/reputap-be/pkg/identifier"
)

var (
	ErrNotFound         = errors.New("product instance not found")
	ErrAlreadyGenerated = errors.New("instances already generated for this order")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GenerateForOrder(ctx context.Context, orderID uint64) ([]InstanceWithQR, error) {
	orderItems, err := s.repo.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("fetching order items: %w", err)
	}

	if len(orderItems) == 0 {
		return nil, fmt.Errorf("no order items found for order %d", orderID)
	}

	type pending struct {
		item  OrderItemInfo
		count int
	}

	var toCreate []pending
	for _, item := range orderItems {
		existing, err := s.repo.CountInstancesByOrderItemID(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("counting instances for order item %d: %w", item.ID, err)
		}

		needed := item.Quantity - existing
		if needed > 0 {
			toCreate = append(toCreate, pending{item: item, count: needed})
		}
	}

	if len(toCreate) == 0 {
		return nil, ErrAlreadyGenerated
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	for _, p := range toCreate {
		for i := 0; i < p.count; i++ {
			serialNumber, err := identifier.Generate()
			if err != nil {
				return nil, fmt.Errorf("generating serial number: %w", err)
			}

			itemID := p.item.ID
			storeID := p.item.StoreID
			instance := &ProductInstance{
				OrderItemID:      &itemID,
				StoreID:          &storeID,
				SerialNumber:     serialNumber,
				NFCStatus:        "not_configured",
				ProductionStatus: "pending",
				QCStatus:         "pending",
			}

			if err := s.repo.Create(ctx, tx, instance); err != nil {
				return nil, fmt.Errorf("creating product instance: %w", err)
			}

			if err := s.repo.CreateQRCode(ctx, tx, instance.ID, p.item.StoreID, serialNumber); err != nil {
				return nil, fmt.Errorf("creating qr code: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	instances, err := s.repo.ListByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("fetching generated instances: %w", err)
	}

	return instances, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*InstanceWithQR, error) {
	instance, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching product instance: %w", err)
	}

	return instance, nil
}

func (s *Service) ListByOrderID(ctx context.Context, orderID uint64) ([]InstanceWithQR, error) {
	return s.repo.ListByOrderID(ctx, orderID)
}

func (s *Service) List(ctx context.Context, params ListParams) ([]InstanceWithQR, int, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) UpdateNFCStatus(ctx context.Context, id uint64, req UpdateNFCStatusRequest, userID uint64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching product instance: %w", err)
	}

	var configuredBy *uint64
	if req.Status == "configured" || req.Status == "tested" {
		configuredBy = &userID
	}

	if err := s.repo.UpdateNFCStatus(ctx, id, req.Status, configuredBy); err != nil {
		return fmt.Errorf("updating nfc status: %w", err)
	}

	return nil
}
