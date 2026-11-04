package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound  = errors.New("product not found")
	ErrSKUExists = errors.New("SKU already exists")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Product, error) {
	exists, err := s.repo.SKUExists(ctx, req.SKU, 0)
	if err != nil {
		return nil, fmt.Errorf("checking SKU uniqueness: %w", err)
	}
	if exists {
		return nil, ErrSKUExists
	}

	status := "active"
	if req.Status != nil {
		status = *req.Status
	}

	product := &Product{
		SKU:            req.SKU,
		Name:           req.Name,
		Description:    req.Description,
		Material:       req.Material,
		WidthMM:        req.WidthMM,
		HeightMM:       req.HeightMM,
		ThicknessMM:    req.ThicknessMM,
		NFCChipType:    req.NFCChipType,
		Finishing:      req.Finishing,
		ProductionCost: req.ProductionCost,
		SellingPrice:   req.SellingPrice,
		Status:         status,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("creating product: %w", err)
	}

	created, err := s.repo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching created product: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching product: %w", err)
	}
	return product, nil
}

func (s *Service) List(ctx context.Context, params ListParams) ([]Product, int, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) Update(ctx context.Context, id uint64, req UpdateRequest) (*Product, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching product: %w", err)
	}

	if req.SKU != nil && *req.SKU != existing.SKU {
		exists, err := s.repo.SKUExists(ctx, *req.SKU, id)
		if err != nil {
			return nil, fmt.Errorf("checking SKU uniqueness: %w", err)
		}
		if exists {
			return nil, ErrSKUExists
		}
		existing.SKU = *req.SKU
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Material = req.Material
	existing.WidthMM = req.WidthMM
	existing.HeightMM = req.HeightMM
	existing.ThicknessMM = req.ThicknessMM
	existing.NFCChipType = req.NFCChipType
	existing.Finishing = req.Finishing
	existing.ProductionCost = req.ProductionCost
	existing.SellingPrice = req.SellingPrice
	if req.Status != nil {
		existing.Status = *req.Status
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("updating product: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetching updated product: %w", err)
	}

	return updated, nil
}
