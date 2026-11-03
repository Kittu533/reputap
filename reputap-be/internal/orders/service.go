package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("order not found")
	ErrInvalidTransition = errors.New("invalid status transition")
)

var allowedTransitions = map[string][]string{
	"pending":       {"confirmed", "cancelled", "on_hold"},
	"confirmed":     {"processing", "cancelled", "on_hold"},
	"processing":    {"production", "cancelled", "on_hold"},
	"production":    {"qc", "on_hold"},
	"qc":            {"ready_to_ship", "production", "on_hold"},
	"ready_to_ship": {"shipped", "on_hold"},
	"shipped":       {"completed"},
	"on_hold":       {"pending", "confirmed", "processing", "production", "qc", "ready_to_ship"},
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest, userID uint64) (*OrderDetail, error) {
	orderNumber, err := s.repo.GenerateOrderNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("generating order number: %w", err)
	}

	var subtotal float64
	for _, item := range req.Items {
		subtotal += float64(item.Quantity) * item.UnitPrice
	}
	total := subtotal - req.Discount

	order := &Order{
		OrderNumber:   orderNumber,
		CustomerID:    req.CustomerID,
		StoreID:       req.StoreID,
		OrderStatus:   "pending",
		PaymentStatus: "unpaid",
		Subtotal:      subtotal,
		Discount:      req.Discount,
		Total:         total,
		Notes:         req.Notes,
		CreatedBy:     userID,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("creating order: %w", err)
	}

	items := make([]OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  float64(item.Quantity) * item.UnitPrice,
		}
	}

	if err := s.repo.CreateItems(ctx, items); err != nil {
		return nil, fmt.Errorf("creating order items: %w", err)
	}

	history := &StatusHistory{
		OrderID:   order.ID,
		ToStatus:  "pending",
		ChangedBy: userID,
	}
	if err := s.repo.CreateStatusHistory(ctx, history); err != nil {
		return nil, fmt.Errorf("creating status history: %w", err)
	}

	created, err := s.repo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching created order: %w", err)
	}

	return &OrderDetail{
		Order: *created,
		Items: items,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*OrderDetail, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching order: %w", err)
	}

	items, err := s.repo.GetItemsByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetching order items: %w", err)
	}

	return &OrderDetail{
		Order: *order,
		Items: items,
	}, nil
}

func (s *Service) List(ctx context.Context, params ListParams) ([]Order, int, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) UpdateStatus(ctx context.Context, id uint64, req UpdateStatusRequest, userID uint64) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching order: %w", err)
	}

	allowed, ok := allowedTransitions[order.OrderStatus]
	if !ok {
		return ErrInvalidTransition
	}

	valid := false
	for _, s := range allowed {
		if s == req.Status {
			valid = true
			break
		}
	}
	if !valid {
		return ErrInvalidTransition
	}

	if err := s.repo.UpdateStatus(ctx, id, req.Status); err != nil {
		return fmt.Errorf("updating order status: %w", err)
	}

	fromStatus := order.OrderStatus
	history := &StatusHistory{
		OrderID:    id,
		FromStatus: &fromStatus,
		ToStatus:   req.Status,
		ChangedBy:  userID,
		Notes:      req.Notes,
	}
	if err := s.repo.CreateStatusHistory(ctx, history); err != nil {
		return fmt.Errorf("creating status history: %w", err)
	}

	return nil
}

func (s *Service) UpdatePaymentStatus(ctx context.Context, id uint64, req UpdatePaymentStatusRequest) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching order: %w", err)
	}

	if err := s.repo.UpdatePaymentStatus(ctx, id, req.PaymentStatus); err != nil {
		return fmt.Errorf("updating payment status: %w", err)
	}

	return nil
}

func (s *Service) GetStatusHistory(ctx context.Context, orderID uint64) ([]StatusHistory, error) {
	_, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching order: %w", err)
	}

	history, err := s.repo.GetStatusHistory(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("fetching status history: %w", err)
	}

	return history, nil
}
