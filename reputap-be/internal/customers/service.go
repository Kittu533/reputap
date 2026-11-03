package customers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("customer not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Customer, error) {
	code, err := s.repo.GetNextCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("generating customer code: %w", err)
	}

	customer := &Customer{
		CustomerCode: code,
		Name:         req.Name,
		BusinessName: req.BusinessName,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		City:         req.City,
		Province:     req.Province,
		Notes:        req.Notes,
	}

	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("creating customer: %w", err)
	}

	created, err := s.repo.GetByID(ctx, customer.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching created customer: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*Customer, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching customer: %w", err)
	}
	return customer, nil
}

func (s *Service) List(ctx context.Context, params ListParams) ([]Customer, int, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) Update(ctx context.Context, id uint64, req UpdateRequest) (*Customer, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching customer: %w", err)
	}

	existing.Name = req.Name
	existing.BusinessName = req.BusinessName
	existing.Email = req.Email
	existing.Phone = req.Phone
	existing.Address = req.Address
	existing.City = req.City
	existing.Province = req.Province
	existing.Notes = req.Notes

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("updating customer: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetching updated customer: %w", err)
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching customer: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting customer: %w", err)
	}

	return nil
}
